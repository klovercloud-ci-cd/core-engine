package service

import "github.com/klovercloud-ci/core/v1"
type Observer interface {
	Listen(subject v1.Subject)
}