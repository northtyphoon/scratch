package lookup

type Metrics struct {
	MetricsName   string
	SucceededCode int
}

type MethodName string

type MethodMap map[MethodName]Metrics

type RouteName string

type RouteMap map[RouteName]MethodMap
