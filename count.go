package metrics

type MetricCount struct {
	metricCommonInt64
}

func newMetricCount(key string, tags AnyTags) *MetricCount {
	metric := metricCountPool.Get().(*MetricCount)
	metric.init(key, tags)
	return metric
}

func (m *MetricCount) init(key string, tags AnyTags) {
	m.metricCommonInt64.init(m, key, tags)
}

func Count(key string, tags AnyTags) *MetricCount {
	if IsDisabled() {
		return (*MetricCount)(nil)
	}

	m := metricsRegistry.Get(TypeCount, key, tags)
	if m != nil {
		return m.(*MetricCount)
	}

	return newMetricCount(key, tags)
}

func (m *MetricCount) GetType() Type {
	return TypeCount
}