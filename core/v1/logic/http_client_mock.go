package logic

import (
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/core/v1/service"
)

type httpClientMockService struct {
}

func (h httpClientMockService) Post(url string, header map[string]string, body []byte) error {
	panic("implement me")
}

func (h httpClientMockService) Get(url string, header map[string]string) ([]byte, error) {
	panic("implement me")
}

// NewHttpClientMockService returns HttpClient type service
func NewHttpClientMockService() service.HttpClient {
	return httpClientMockService{}
}
