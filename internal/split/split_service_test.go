package split

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/splitio/go-client/v6/splitio/client"
	"github.com/stretchr/testify/require"
)

// Used for unit-testing
var mockSplits = "feature1 on\nfeature2 off\n"

func TestNewSplitService(t *testing.T) {
	service, _ := NewService("localhost")

	actualType := reflect.TypeOf(service)
	expectedType := reflect.TypeOf(&Srvc{client: &client.SplitClient{}})

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

func setup() *Srvc {
	// Set a custom split config
	splitConfig.SplitFile = createMockSplits()
	service, _ := NewService("localhost")
	return service
}

// Creates a temporary file with splits used for unit-testing
func createMockSplits() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("could not get user home directory: %s", err)
	}
	// Create a temporary file that holds test splits for testing purposes
	blob := []byte(mockSplits)
	filename := dir + "/mock.split"
	err = ioutil.WriteFile(filename, blob, 0777)
	if err != nil {
		log.Errorf("could not create temp file: %s", err)
	}

	return filename
}
