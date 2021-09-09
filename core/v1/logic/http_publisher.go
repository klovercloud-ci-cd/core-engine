package logic

import "github.com/klovercloud-ci/core/v1/service"

type httpPublisherService struct {

}

func (h httpPublisherService) Post(url string, header map[string]string, body []byte) error{
	return nil
}

func NewHttpPublisherService() service.HttpPublisher {
	return &httpPublisherService{
	}
}