package logic

import (
	"bufio"
	"context"
	"github.com/klovercloud-ci-cd/core-engine/config"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"strings"
	"time"
)

type k8sService struct {
	Kcs          *kubernetes.Clientset
	tekton       service.Tekton
	observerList []service.Observer
}

func (k8s k8sService) DeletePersistentVolumeClaimByProcessId(processId string) error {
	list, err := k8s.Kcs.CoreV1().PersistentVolumeClaims(config.CiNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: "processId=" + processId,
	})
	if err != nil {
		log.Println("[WARNING]:", err.Error())
		return err
	}
	for _, each := range list.Items {
		err = k8s.Kcs.CoreV1().PersistentVolumeClaims(config.CiNamespace).Delete(context.Background(), each.Name, metav1.DeleteOptions{})
		if err != nil {
			log.Println("[ERROR]:", err.Error())
		}
	}
	return nil
}

func (k8s k8sService) CreatePersistentVolumeClaim(source corev1.PersistentVolumeClaim) error {
	_, err := k8s.Kcs.CoreV1().PersistentVolumeClaims(source.Namespace).Create(context.Background(), &source, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (k8s k8sService) InitPersistentVolumeClaim(step v1.Step, label map[string]string, processId string) corev1.PersistentVolumeClaim {
	label = make(map[string]string)
	label["step"] = step.Name
	label["processId"] = processId
	pvc := corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      step.Name + "-" + processId,
			Labels:    label,
			Namespace: config.CiNamespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.PersistentVolumeAccessMode(step.Params[enums.RESOURCE_ACCESS_MODE])},
			Resources: corev1.ResourceRequirements{
				Limits:   nil,
				Requests: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: resource.MustParse(step.Params[enums.RESOURCE_SIZE])},
			},
		},
	}
	return pvc
}

func (k8s k8sService) GetContainerLog(namespace, podName, containerName string, taskRunLabel map[string]string) (io.ReadCloser, error) {
	req := k8s.RequestContainerLog(namespace, podName, containerName)
	readCloser, err := req.Stream(context.Background())
	if err != nil {
		log.Println(err.Error())

		if strings.Contains(err.Error(), "not found") {
			return nil, err
		}
		time.Sleep(time.Second)
		taskRunName := taskRunLabel["step"] + "-" + taskRunLabel["processId"]
		tRun, tRunError := k8s.tekton.GetTaskRun(taskRunName, true)

		if tRunError == nil && !tRun.IsCancelled() {
			return k8s.GetContainerLog(namespace, podName, containerName, taskRunLabel)
		}
		return nil, err
	}

	return readCloser, nil
}

func (k8s k8sService) FollowContainerLifeCycle(namespace, podName, containerName, step string, processId string, stepType enums.STEP_TYPE) {
	processEventData := make(map[string]interface{})
	processEventData["step"] = step
	req := k8s.RequestContainerLog(namespace, podName, containerName)
	readCloser, err := req.Stream(context.Background())
	for err != nil {
		listener := v1.Subject{Pipeline: v1.Pipeline{ProcessId: processId}, Log: err.Error(), Step: step}
		if strings.Contains(err.Error(), "image can't be pulled") || strings.Contains(err.Error(), "pods \""+podName+"\" not found") || strings.Contains(err.Error(), "pod \""+podName+"\" is terminated") {
			if stepType == enums.BUILD {
				processEventData["status"] = enums.STEP_FAILED
				processEventData["reason"] = err.Error()
				listener.EventData = processEventData
				go k8s.notifyAll(listener)
			}
			return
		}
		readCloser, err = req.Stream(context.Background())
	}
	reader := bufio.NewReaderSize(readCloser, 64)
	lastLine := ""
	for {
		data, isPrefix, err := reader.ReadLine()
		if err != nil {
			listener := v1.Subject{Pipeline: v1.Pipeline{ProcessId: processId}, Log: err.Error(), Step: step}
			processEventData["reason"] = err.Error()
			listener.EventData = make(map[string]interface{})
			listener.EventData = processEventData
			go k8s.notifyAll(listener)
			return
		}
		lines := strings.Split(string(data), "\r")
		length := len(lines)
		if len(lastLine) > 0 {
			lines[0] = lastLine + lines[0]
			lastLine = ""
		}
		if isPrefix {
			lastLine = lines[length-1]
			lines = lines[:(length - 1)]
		}
		for _, line := range lines {
			log.Println(containerName, podName, "+++++++", line)
			temp := strings.ToLower(line)
			processEventData["log"] = temp
			processEventData["reason"] = "n/a"
			listener := v1.Subject{Pipeline: v1.Pipeline{ProcessId: processId}, Log: temp, Step: step}
			listener.EventData = make(map[string]interface{})
			listener.EventData = processEventData
			go k8s.notifyAll(listener)
			if (!strings.HasPrefix(temp, "progress") && (!strings.HasSuffix(temp, " mb") || !strings.HasSuffix(temp, " kb"))) && !strings.HasPrefix(temp, "downloading from") {

			}
			if strings.Contains(line, "image can't be pulled") || strings.Contains(line, "pods \""+podName+"\" not found") {
				if strings.Contains(line, "image can't be pulled") {

				}
				break
			}

		}
	}
	if readCloser != nil {
		readCloser.Close()
	}

}

func (k8s k8sService) RequestContainerLog(namespace string, podName string, containerName string) *rest.Request {
	return k8s.Kcs.CoreV1().Pods(namespace).GetLogs(
		podName,
		&corev1.PodLogOptions{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Task",
				APIVersion: "tekton.dev/v1",
			},
			Container: containerName,
			Follow:    true,
		},
	)
}

func (k8s *k8sService) GetSecret(name, namespace string) (corev1.Secret, error) {
	sec, err := k8s.Kcs.CoreV1().
		Secrets(namespace).
		Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return corev1.Secret{}, err
	}
	return *sec, nil
}

func (k8s *k8sService) GetConfigMap(name, namespace string) (corev1.ConfigMap, error) {
	sec, err := k8s.Kcs.CoreV1().
		ConfigMaps(namespace).
		Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return corev1.ConfigMap{}, err
	}
	return *sec, nil
}

func (k8s *k8sService) GetPodListByProcessId(namespace, processId string, option v1.PodListGetOption) *corev1.PodList {
	listener := v1.Subject{}
	data := make(map[string]interface{})
	listener.Pipeline.ProcessId = processId
	labelSelector := "processId=" + processId
	podList, err := k8s.Kcs.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Println(err.Error())
	}
	count := 0
	for len(podList.Items) == 0 && option.Wait {
		data["status"] = enums.POD_INITIALIZING
		listener.EventData = data
		go k8s.notifyAll(listener)
		count++
		podList, err = k8s.Kcs.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err == nil && len(podList.Items) == 0 && count > option.Duration {
			return podList
		}
		time.Sleep(1 * time.Second)
	}
	return podList
}

func (k8s *k8sService) WaitAndGetInitializedPods(namespace, processId, step string) *corev1.PodList {
	var podList *corev1.PodList
	listener := v1.Subject{}
	data := make(map[string]interface{})
	data["step"] = step
	listener.Pipeline.ProcessId = processId
	podList = k8s.GetPodListByProcessId(namespace, processId, v1.PodListGetOption{
		Wait:     true,
		Duration: enums.DEFAULT_POD_INITIALIZATION_WAIT_DURATION,
	})
	if len(podList.Items) < 1 {
		data["status"] = enums.WAITING
		listener.EventData = data
		go k8s.notifyAll(listener)
		return podList
	}
	podStatus := podList.Items[0].Status.Phase
	if enums.POD_STATUS(podStatus) == enums.POD_TERMINATING {
		data["status"] = enums.TERMINATING
		listener.EventData = data
		go k8s.notifyAll(listener)
		return k8s.WaitAndGetInitializedPods(namespace, processId, step)
	}
	if enums.POD_STATUS(podStatus) == enums.POD_INITIALIZING {
		data["status"] = enums.INITIALIZING
		listener.EventData = data
		go k8s.notifyAll(listener)
	}
	return podList
}

func (k8s k8sService) notifyAll(listener v1.Subject) {
	for _, observer := range k8s.observerList {
		go observer.Listen(listener)
	}
}

// NewK8sService returns K8s type service.
func NewK8sService(Kcs *kubernetes.Clientset, tekton service.Tekton, observerList []service.Observer) service.K8s {
	return &k8sService{
		Kcs:          Kcs,
		tekton:       tekton,
		observerList: observerList,
	}
}
