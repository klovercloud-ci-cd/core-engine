package inmemory

import (
	"container/list"
	v1 "github.com/klovercloud-ci/core/v1"
)

// IndexedLogEvents log events store
var IndexedLogEvents map[string][]v1.LogEvent

// ProcessEventStore process events store
var ProcessEventStore map[string]*list.List
