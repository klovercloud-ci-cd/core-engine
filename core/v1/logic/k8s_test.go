package logic

import (
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"strconv"
	"testing"
)

func Test_GetSecret(t *testing.T) {
	type TestData struct {
		data     []corev1.Secret
		expected corev1.Secret
		actual   corev1.Secret
	}

	var testCases []TestData
	data := InitSecrets()
	for i := 0; i < 10; i++ {
		resource := &mockK8sService{}
		secret, _ := resource.GetSecret("secret"+strconv.Itoa(i), "klovercloud")
		testCases = append(testCases, TestData{
			data:     data,
			expected: data[i],
			actual:   secret,
		})
	}

	for _, each := range testCases {
		if !reflect.DeepEqual(each.expected, each.actual) {
			assert.ElementsMatch(t, each.expected, each.actual)
		}
	}

}
