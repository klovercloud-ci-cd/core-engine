package service

import "github.com/klovercloud-ci-cd/core-engine/core/v1"

// Observer Observer operations.
type Observer interface {
	Listen(subject v1.Subject)
}
