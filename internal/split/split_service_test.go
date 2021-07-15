package split

import (
	"reflect"
	"testing"

	"github.com/splitio/go-client/v6/splitio/client"
	"github.com/stretchr/testify/require"
)

func TestNewSplitService(t *testing.T) {
	service, _ := NewSplitService("localhost", "staging")

	actualType := reflect.TypeOf(service)
	expectedType := reflect.TypeOf(&Service{client: &client.SplitClient{}})

	require.NotNil(t, service)
	require.Equal(t, expectedType, actualType)
}

func TestGet(t *testing.T) {
	service := setup()

	actualTreatment := service.Get("feature1")
	expectedTreatment := "on"

	require.Equal(t, expectedTreatment, actualTreatment)
}

func TestGetAll(t *testing.T) {
	service := setup()

	splits := []string{"feature1", "feature2"}
	actualTreatments := service.GetAll(splits)
	expectedTreatments := map[string]string{
		"feature1": "on",
		"feature2": "off",
	}

	require.Equal(t, expectedTreatments, actualTreatments)
}

func setup() *Service {
	service, _ := NewSplitService("localhost", "staging")

	return service
}
