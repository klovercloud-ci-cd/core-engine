package in_memory

import v1 "github.com/klovercloud-ci/core/v1"

var IndexedLogEvents map[string]map[int32]v1.LogEvent


var LogEventsStore []v1.LogEvent
