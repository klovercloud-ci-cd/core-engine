package service

import "github.com/klovercloud-ci/core/v1"
type Observer interface {
	Listen(v1.Subject)
}