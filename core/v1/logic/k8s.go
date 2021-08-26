package logic

import (
	"bufio"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"log"
	"strings"
	"time"
)

type k8sService struct {
	Kcs *kubernetes.Clientset
	repo repository.LogEventRepository
}

func (k8s k8sService) LogContainer(namespace, podName, containerName, step, buildId string,stepType enums.STEP_TYPE ) {
	req := k8s.Kcs.CoreV1().Pods(namespace).GetLogs(
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

	readCloser, err := req.Stream()
	for err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "image can't be pulled") || strings.Contains(err.Error(), "pods \""+podName+"\" not found") || strings.Contains(err.Error(), "pod \""+podName+"\" is terminated") {

			if stepType==enums.BUILD{
					//build step failed
			}else if stepType==enums.DEPLOY{
				//deploy step failed
			}
			return
		}else{
			readCloser, err = req.Stream()
		}

	}
	reader := bufio.NewReaderSize(readCloser, 64)
	lastLine := ""
	for {
		data, isPrefix, err := reader.ReadLine()
		if err != nil {
			log.Println(err)
			//log.Println("appId=" + taskrun.AppId + ", appType=" + taskrun.AppType + ", buildId=" + buildId + ", taskType=" + taskType + ", revision=" + taskrun.Input.Revision + ", error=" + err.Error())
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
			temp := strings.ToLower(line)
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



func (k8s * k8sService) GetSecret(name,namespace string)(corev1.Secret,error){
	sec,err:=k8s.Kcs.CoreV1().
		Secrets(namespace).
		Get(name,metav1.GetOptions{})
	if err!=nil{
		return corev1.Secret{}, err
	}
	return *sec,nil
}

func (k8s * k8sService) GetConfigMap(name,namespace string)(corev1.ConfigMap,error){
	sec,err:=k8s.Kcs.CoreV1().
		ConfigMaps(namespace).
		Get(name,metav1.GetOptions{})
	if err!=nil{
		return corev1.ConfigMap{}, err
	}
	return *sec,nil
}

func (k8s * k8sService) GetPodListByBuildId(namespace,buildId string,option v1.PodListGetOption) *corev1.PodList{
	labelSelector := "buildId=" + buildId
	podList, err := k8s.Kcs.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	count := 0
	for len(podList.Items) == 0 && option.Wait {
		count++
		podList, err = k8s.Kcs.CoreV1().Pods(namespace).List(metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err == nil && len(podList.Items) == 0 && count > option.Duration {
			return podList
		}
		time.Sleep(1 * time.Second)
	}
	return podList
}

func (k8s * k8sService) WaitAndGetInitializedPods(namespace,buildId string) *corev1.PodList{
	var podList *corev1.PodList

	podList = k8s.GetPodListByBuildId(namespace,buildId,v1.PodListGetOption{
		Wait:     true,
		Duration: enums.DEFAULT_POD_INITIALIZATION_WAIT_DURATION,
	})
	if len(podList.Items) < 1 {
		return podList
	}
	podStatus := podList.Items[0].Status.Phase
	if podStatus == enums.TERMINATING {
		log.Println("Pod is:", podStatus)
		return k8s.WaitAndGetInitializedPods(namespace, buildId)
	}

	if podStatus == enums.POD_INITIALIZING {
		log.Println("Pod is:", podStatus)
		podStatus = podList.Items[0].Status.Phase
	}
	return podList
}

func NewK8sService(Kcs *kubernetes.Clientset,repo repository.LogEventRepository) service.K8s {
	return &k8sService{
		Kcs:  Kcs,
		repo: repo,
	}
}