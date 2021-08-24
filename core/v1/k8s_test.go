package v1

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"reflect"
	"strconv"
	"testing"
)

func Test_getSecret(t *testing.T) {
	type TestData struct {
		data                  []v1.Secret
		expected v1.Secret
		actual v1.Secret
	}

	var testCases [] TestData
	data:=InitSecrets()
	for i:=0;i<10;i++{
		resource := &MockK8sResource{}
		secret,_:=resource.getSecret("secret"+strconv.Itoa(i),"klovercloud")
		testCases= append(testCases, TestData{
			data: data,
			expected: data[i],
			actual: secret,
		})
	}

	for _,each:=range testCases{
		if !reflect.DeepEqual(each.expected, each.actual){
			assert.ElementsMatch(t, each.expected, each.actual)
		}
	}

}