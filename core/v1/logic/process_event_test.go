package logic

//func TestProcessEventService_GetByProcessId(t *testing.T) {
//	type TestCase struct {
//		processId string
//		data      []v1.ProcessEvent
//		expected  map[string]interface{}
//		actions   []string
//		actual    map[string]interface{}
//	}
//	var testCases []TestCase
//	testCases = append(testCases, TestCase{
//		processId: "1",
//		data: []v1.ProcessEvent{{
//			ProcessId: "1",
//			Data:      map[string]interface{}{"step": "build", "status": "Processing"},
//		},
//			{
//				ProcessId: "1",
//				Data:      map[string]interface{}{"step": "build", "status": "Pod Initializing"},
//			},
//		},
//		actions:  []string{"getFront"},
//		expected: map[string]interface{}{"step": "build", "status": "Processing"},
//		actual:   map[string]interface{}{"step": "build", "status": "Processing"},
//	})
//
//	testCases = append(testCases, TestCase{
//		processId: "1",
//		actions:   []string{"getFront", "removeFront", "removeFront"},
//		expected:  map[string]interface{}{"step": "build", "status": "Pod Initializing"},
//		actual:    map[string]interface{}{"step": "build", "status": "Pod Initializing"},
//	})
//	repo := NewProcessEventService(in_memory.NewProcessEventRepository())
//	for _, each := range testCases {
//		if len(each.data) > 0 {
//			for _, d := range each.data {
//				repo.Store(d)
//			}
//		}
//		for _, action := range each.actions {
//			if action == "getFront" {
//				each.actual = repo.GetByCompanyId(each.processId)
//			} else if action == "removeFront" {
//				each.actual = repo.DequeueByCompanyIdAndUserId(each.processId)
//			}
//		}
//		if !reflect.DeepEqual(each.expected, each.actual) {
//			assert.ElementsMatch(t, each.expected, each.actual)
//		}
//	}
//}
//
//func TestProcessEventService_Store(t *testing.T) {
//	type TestCase struct {
//		data     v1.ProcessEvent
//		expected map[string]interface{}
//		actions  []string
//		actual   map[string]interface{}
//	}
//
//	var testCases []TestCase
//
//	testCases = append(testCases, TestCase{
//		data: v1.ProcessEvent{
//			ProcessId: "1",
//			Data:      map[string]interface{}{"step": "build", "status": "Processing"},
//		},
//		actions:  []string{"getFront"},
//		expected: map[string]interface{}{"step": "build", "status": "Processing"},
//		actual:   map[string]interface{}{"step": "build", "status": "Processing"},
//	})
//
//	testCases = append(testCases, TestCase{
//		data: v1.ProcessEvent{
//			ProcessId: "2",
//			Data:      map[string]interface{}{"step": "build", "status": "Processing"},
//		},
//		actions:  []string{"getFront"},
//		expected: map[string]interface{}{"step": "build", "status": "Processing"},
//		actual:   map[string]interface{}{"step": "build", "status": "Processing"},
//	})
//
//	testCases = append(testCases, TestCase{
//		data: v1.ProcessEvent{
//			ProcessId: "1",
//			Data:      map[string]interface{}{"step": "build", "status": "Pod Initializing"},
//		},
//		actions:  []string{"getFront", "removeFront", "getFront"},
//		expected: map[string]interface{}{"step": "build", "status": "Pod Initializing"},
//		actual:   map[string]interface{}{"step": "build", "status": "Pod Initializing"},
//	})
//
//	testCases = append(testCases, TestCase{
//		data: v1.ProcessEvent{
//			ProcessId: "1",
//			Data:      map[string]interface{}{"step": "build", "status": "Image pulling"},
//		},
//		actions:  []string{"removeFront", "removeFront", "removeFront"},
//		expected: nil,
//		actual:   nil,
//	})
//
//	repo := NewProcessEventService(in_memory.NewProcessEventRepository())
//	for _, each := range testCases {
//		repo.Store(each.data)
//		for _, action := range each.actions {
//			if action == "getFront" {
//				each.actual = repo.GetByCompanyId(each.data.ProcessId)
//			} else if action == "removeFront" {
//				each.actual = repo.DequeueByCompanyIdAndUserId(each.data.ProcessId)
//			}
//		}
//		if !reflect.DeepEqual(each.expected, each.actual) {
//			assert.ElementsMatch(t, each.expected, each.actual)
//		}
//	}
//}
