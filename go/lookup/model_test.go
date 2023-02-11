package lookup

import (
	"net/http"
	"testing"
)

func TestLookupModel(t *testing.T) {
	var routeMap = RouteMap{
		"route1": {
			"method1": {"metrics1", http.StatusOK},
			"method2": {"metrics2", http.StatusOK},
		},
		"route3": {
			"method3": {"metrics3", http.StatusOK},
		},
	}

	tests := []struct {
		route  string
		method string
	}{
		{
			"route1",
			"method1",
		},
		{
			"route1",
			"method2",
		},
		{
			"route1",
			"method3",
		},
		{
			"route2",
			"method1",
		},
	}

	for _, test := range tests {
		if route, ok := routeMap[test.route]; ok {
			if operation, ok := route[test.method]; ok {
				t.Logf("%v", operation)
			} else {
				t.Logf("%s/%s not found", test.route, test.method)
			}
		} else {
			t.Logf("%s not found", test.route)
		}
	}
}
