package execution

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	statusNotStarted = "NOT_STARTED"
	statusStarted    = "STARTED"
	statusDone       = "DONE"
	statusFailed     = "FAILED"
	statusSkipped    = "SKIPPED"
)

type RhapsodyClient struct {
	client *http.Client
	url    string
}

func NewRhapsodyClient() *RhapsodyClient {
	client := http.Client{}
	url := os.Getenv("NEW_RELIC_RHAPSODY_URL")
	if url == "" {
		url = "http://127.0.0.1:4000" // apollo default
	}
	return &RhapsodyClient{client: &client, url: url}
}

// func setInstallItemStatus(status RecipeStatus) string {
// 	installItemStatus := ""
// 	switch status.Status {
// 	case RecipeStatusTypes.INSTALLED:
// 		installItemStatus = statusDone

// 	case RecipeStatusTypes.FAILED:
// 		installItemStatus = statusFailed

// 	case RecipeStatusTypes.SKIPPED:
// 		installItemStatus = statusSkipped
// 	default:
// 	}
// 	return installItemStatus
// }

type InstallItem struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	LastUpdated  time.Time `json:"lastUpdated"`
	DataSourceID string    `json:"dataSourceId"`
	InstallID    string    `json:"installId"`
	IsVerified   bool      `json:"verified"`
	RecipeName   string    `json:"name"`
}

type rhapsodyRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

var mutationCreateInstallAttemptItemCLI = `mutation CreateInstallAttemptItemCLI($attemptId: ID!, $recipeName: String!, $installId: String!) { createInstallAttemptItemCLI(attemptId: $attemptId, recipeName: $recipeName, installId: $installId) { id } }`
var mutationUpdateInstallItemStatus = `
mutation UpdateInstallItemStatus($itemId: ID!, $status: InstallStatus!) {
	updateInstallItemStatus(itemId: $itemId, status: $status) {
		id
	}
}
`

func UpdateRhapsody(s *InstallStatus) error {
	c := NewRhapsodyClient()

	attemptID := os.Getenv("RHAPSODY_ATTEMPT_ID")
	if attemptID == "" {
		return fmt.Errorf("must set RHAPSODY_ATTEMPT_ID")
	}

	items := make([]InstallItem, len(s.Statuses))
	for i, status := range s.Statuses {
		items[i] = InstallItem{
			RecipeName: status.Name,
		}
	}
	// create items
	for i, item := range items {
		request := rhapsodyRequest{
			Query: mutationCreateInstallAttemptItemCLI,
			Variables: map[string]string{
				"attemptId":  attemptID,
				"recipeName": item.RecipeName,
				"installId":  s.InstallID,
			},
		}

		reqBody, err := json.Marshal(request)
		if err != nil {
			return err
		}
		res, err := c.client.Post(c.url, "application/json", bytes.NewReader(reqBody))
		if err != nil {
			return err
		}
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		var createItemRes createInstallAttemptItemCLIResponse
		if err := json.Unmarshal(resBody, &createItemRes); err != nil {
			return err
		}
		items[i].ID = createItemRes.Data.CreateInstallAttemptItemCLI.ID
	}

	// update items
	for _, item := range items {
		request := rhapsodyRequest{
			Query: mutationUpdateInstallItemStatus,
			Variables: map[string]string{
				"itemId": item.ID,
				"status": statusDone,
			},
		}

		reqBody, err := json.Marshal(request)
		if err != nil {
			return err
		}
		if _, err := c.client.Post(c.url, "application/json", bytes.NewReader(reqBody)); err != nil {
			return err
		}
	}
	return nil
}

type createInstallAttemptItemCLIResponse struct {
	Data struct {
		CreateInstallAttemptItemCLI struct {
			ID string `json:"id"`
		} `json:"createInstallAttemptItemCLI"`
	} `json:"data"`
}

// StatusSubscriber solution

// type InstallAttempt struct {
// 	ID                 string        `json:"id"`
// 	TargetDataSourceID string        `json:"targetDataSourceId"`
// 	InstallItems       []InstallItem `json:"items"`
// 	CliInstallID       string
// }

// type InstallItem struct {
// 	ID           string    `json:"id"`
// 	Status       string    `json:"status"`
// 	LastUpdated  time.Time `json:"lastUpdated"`
// 	DataSourceID string    `json:"dataSourceId"`
// 	InstallID    string    `json:"installId"`
// 	IsVerified   bool      `json:"verified"`
// 	RecipeName   string    `json:"name"`
// }

// var dataSourceIDLookup = map[string]string{
// 	"node-agent-installer":           "node-js",
// 	"infrastructure-agent-installer": "infra", //todo -- IDs are os/distro specific
// 	"logs-integration":               "logs",  //todo -- same as above
// }

// func (c *RhapsodyClient) createRhapsodyInstallEvent(status *InstallStatus) (string, error) {
// 	items := make([]InstallItem, len(status.Detected))
// 	for i, detected := range status.Detected {
// 		items[i] = InstallItem{
// 			Status:       statusNotStarted,
// 			LastUpdated:  time.Now(),
// 			DataSourceID: dataSourceIDLookup[detected.Name],
// 		}
// 	}
// 	i := InstallAttempt{
// 		InstallItems: items,
// 	}
// 	reqBodyJson, err := json.Marshal(i)
// 	if err != nil {
// 		return "", err
// 	}
// 	reqBody := bytes.NewReader(reqBodyJson)
// 	res, err := c.client.Post(c.url, "application/json", reqBody)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer res.Body.Close()
// 	// resBody, err := io.ReadAll(res.Body)
// 	// if err != nil {
// 	// 	return "", err
// 	// }
// 	// fmt.Println(resBody)
// 	// json.Unmarshal(resBody, &foo)
// 	return "id", nil
// }

// func (c *RhapsodyClient) updateRhapsodyInstallEvent(id string, status *InstallStatus) error {
// 	//todo
// 	if id == "" {
// 		return fmt.Errorf("no rhapsody id provided")
// 	}
// 	return nil
// }

// type RhapsodyReporter struct {
// 	client     *RhapsodyClient
// 	accountID  int
// 	rhapsodyID string
// }

// func NewRhapsodyReporter() *RhapsodyReporter {
// 	rhapsodyID := os.Getenv("NEW_RELIC_INSTALL_ID")
// 	if rhapsodyID != "" {
// 		log.Debugf("found NEW_RELIC_INSTALL_ID %s", rhapsodyID)
// 	}

// 	r := RhapsodyReporter{
// 		client:     NewRhapsodyClient(),
// 		accountID:  configAPI.GetActiveProfileAccountID(),
// 		rhapsodyID: rhapsodyID,
// 	}
// 	return &r
// }

// func (r *RhapsodyReporter) createRhapsodyInstallEvent(status *InstallStatus) error {
// 	id, err := r.client.createRhapsodyInstallEvent(status)
// 	if err != nil {
// 		return err
// 	}
// 	r.rhapsodyID = id
// 	return nil
// }

// func (r *RhapsodyReporter) updateRhapsodyInstallEvent(status *InstallStatus) error {
// 	return r.client.updateRhapsodyInstallEvent(r.rhapsodyID, status)
// }

// // StatusSubscriber methods -- implemented

// func (r *RhapsodyReporter) InstallStarted(status *InstallStatus) error {
// 	if r.rhapsodyID != "" {
// 		log.Debugf("rhapsody: skipping creating event, using ID %s", r.rhapsodyID)
// 	}
// 	return r.createRhapsodyInstallEvent(status)
// }
// func (r *RhapsodyReporter) InstallComplete(status *InstallStatus) error {
// 	return r.updateRhapsodyInstallEvent(status)
// }
// func (r *RhapsodyReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
// 	return r.updateRhapsodyInstallEvent(status)
// }
// func (r *RhapsodyReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
// 	return r.updateRhapsodyInstallEvent(status)
// }

// // StatusSubscriber methods -- not implemented

// func (r *RhapsodyReporter) UpdateRequired(status *InstallStatus) error {
// 	return nil
// }
// func (r *RhapsodyReporter) InstallCanceled(status *InstallStatus) error {
// 	return nil
// }
// func (r *RhapsodyReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
// 	return nil
// }
// func (r *RhapsodyReporter) RecipeDetected(status *InstallStatus, event RecipeStatusEvent) error {
// 	return nil
// }
// func (r *RhapsodyReporter) RecipeCanceled(status *InstallStatus, event RecipeStatusEvent) error {
// 	return nil
// }
// func (r *RhapsodyReporter) RecipeAvailable(status *InstallStatus, event RecipeStatusEvent) error {
// 	return nil
// }
// func (r *RhapsodyReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
// 	return nil
// }
// func (r *RhapsodyReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
// 	return nil
// }
// func (r *RhapsodyReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
// 	return nil
// }
// func (r *RhapsodyReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
// 	return nil
// }
// func (r *RhapsodyReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
// 	return nil
// }
