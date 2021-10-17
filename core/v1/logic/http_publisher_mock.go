package logic

import (
	"github.com/klovercloud-ci/core/v1/service"
)

type httpPublisherMockService struct {
}

func (h httpPublisherMockService) Get(url string, header map[string]string) (error, []byte) {
	panic("implement me")
}

func (h httpPublisherMockService) Post(url string, header map[string]string, body []byte) error {
	return nil
}

func NewHttpPublisherMockService() service.HttpPublisher {
	return &httpPublisherMockService{}
}
