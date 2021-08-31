package in_memory

import (
	"container/list"
	v1 "github.com/klovercloud-ci/core/v1"
)

var IndexedLogEvents map[string][]v1.LogEvent
var ProcessEventStore map[string]*list.List


