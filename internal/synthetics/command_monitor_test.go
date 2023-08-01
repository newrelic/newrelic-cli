package synthetics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestSyntheticsMonitor(t *testing.T) {
	assert.Equal(t, "monitor", cmdMon.Name())

	testcobra.CheckCobraMetadata(t, cmdMon)
	testcobra.CheckCobraRequiredFlags(t, cmdMon, []string{})
}

func TestSyntheticsMonitorGet(t *testing.T) {
	assert.Equal(t, "get", cmdMonGet.Name())

	testcobra.CheckCobraMetadata(t, cmdMonGet)
	testcobra.CheckCobraRequiredFlags(t, cmdMonGet, []string{})
}

func TestSyntheticsMonitorSearch(t *testing.T) {
	assert.Equal(t, "search", cmdMonSearch.Name())

	testcobra.CheckCobraMetadata(t, cmdMonSearch)
	testcobra.CheckCobraRequiredFlags(t, cmdMonSearch, []string{})
}

func TestSyntheticsMonitorList(t *testing.T) {
	assert.Equal(t, "list", cmdMonList.Name())

	testcobra.CheckCobraMetadata(t, cmdMonList)
	testcobra.CheckCobraRequiredFlags(t, cmdMonList, []string{})
}

// CHANGE THIS VARIABLE to use the three mock JSON scenarios
var scenario = 1

func TestNewFunctionality(t *testing.T) {
	//original : batchID would come from the result of the syntheticsStartAutomationTest mutation
	batchID := "36488dff-9a8a-4358-8ef2-e73da3c118e0"
	apiTimeout := time.Minute * 5

	start := time.Now()
	i := 0
	for {
		if time.Since(start) > apiTimeout {
			fmt.Println("Halting execution : reached timeout.")
			break
		}

		i++

		//result, err := functionInClientGo(automatedTestResultQueryInput{batchID: batchID})
		root := fakeAutomatedTestResultQuery(batchID, i)
		//if err != nil {
		//	return fmt.Errorf("Some error")
		//}

		if root.Status == "TIMED_OUT" || root.Status == "PASSED" || root.Status == "FAILED" {
			fmt.Printf("Execution stopped - Status: %s\n", root.Status)
			fakePrintMonitorStatus(root)
			break
		} else if root.Status == "IN_PROGRESS" {
			fmt.Println("Status still IN_PROGRESS, calling API again in 15 seconds")
			fakePrintMonitorStatus(root)
			fmt.Println()
			time.Sleep(time.Second * 15)
		} else {
			fmt.Printf("Unexpected status: %s\n", root.Status)
			break
		}
	}

}

func fakeAutomatedTestResultQuery(batchID string, index int) (r Root) {
	directory := fmt.Sprintf("mock_json/Scenario %d", scenario)
	filePath := fmt.Sprintf("%s/response_%d.json", directory, index)
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var root Root
	if err := json.Unmarshal(byteValue, &root); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	return root
}

func fakePrintMonitorStatus(root Root) {
	countSuccess := 0
	countFailure := 0
	countProgress := 0
	tests := root.Tests

	for _, test := range tests {
		if test.Result == "SUCCESS" {
			countSuccess++
		} else if test.Result == "FAILED" {
			countFailure++
		} else if test.Result == "IN_PROGRESS" || test.Result == "" {
			countProgress++
		}
	}

	fmt.Println("Successful Tests: ", countSuccess)
	fmt.Println("Failed Tests: ", countFailure)
}

// Temporary Data Structures
type Test struct {
	Id          string `json:"id"`
	BatchID     string `json:"batchId"`
	MonitorID   string `json:"monitorId"`
	MonitorGUID string `json:"monitorGuid"`
	MonitorName string `json:"monitorName"`
	Result      string `json:"result"`
}

type Root struct {
	Tests  []Test                 `json:"tests"`
	Config map[string]interface{} `json:"config"`
	Status string                 `json:"status"`
}
