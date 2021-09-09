package service

import v1 "github.com/klovercloud-ci/core/v1"

type AgentEvent interface {
	Listen(v1.Subject)
}