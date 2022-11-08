package discovery

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/newrelic/newrelic-cli/internal/install/discovery/mocks"
	"testing"
)

var mockContext context.Context
var mockCtrl *gomock.Controller
var mockDiscoverer *mocks.MockDiscoverer

func setup(t *testing.T) {
	mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDiscoverer = mocks.NewMockDiscoverer(mockCtrl)
	mockContext = context.Background()
}
