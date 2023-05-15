package segment

import (
	"gopkg.in/segmentio/analytics-go.v3"
)

type MockSegmentClient struct {
	EnqueueCallCount int
	MessageQueued    []analytics.Track
	IsCloseCalled    bool
}

func (m *MockSegmentClient) Enqueue(msg analytics.Message) error {
	m.EnqueueCallCount = m.EnqueueCallCount + 1
	m.MessageQueued = append(m.MessageQueued, msg.(analytics.Track))
	return nil
}

func (m MockSegmentClient) Close() error {
	m.IsCloseCalled = true
	return nil
}
