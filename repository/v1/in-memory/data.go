package in_memory

import v1 "github.com/klovercloud-ci/core/v1"

var indexedLogEvents map[string]map[string]v1.LogEvent

var LogEventsStore []v1.LogEvent
