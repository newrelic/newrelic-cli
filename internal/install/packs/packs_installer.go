package packs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	log "github.com/sirupsen/logrus"
)

// Maybe both fetchPacks and installPacks should be one interface
// vs separate interfaces?
type PacksInstaller interface {
	Install(ctx context.Context, packs []types.OpenInstallationObservabilityPack) error
}

type ConcretePacksInstaller struct {
	client            *newrelic.NewRelic
	progressIndicator ux.ProgressIndicator
}

func NewConcretePacksInstaller(client *newrelic.NewRelic) *ConcretePacksInstaller {
	return &ConcretePacksInstaller{
		client:            client,
		progressIndicator: ux.NewPlainProgress(),
	}
}

func defaultHTTPGetFunc(dashboardURL string) (*http.Response, error) {
	return http.Get(dashboardURL)
}

func (p *ConcretePacksInstaller) Install(ctx context.Context, packs []types.OpenInstallationObservabilityPack) error {
	msg := fmt.Sprintf("Installing observability packs")
	p.progressIndicator.Start(msg)
	defer func() { p.progressIndicator.Stop() }()

	for _, pack := range packs {

		// Only installing dashboards currently
		if pack.Dashboards != nil {
			for _, dashboard := range pack.Dashboards {
				if _, err := p.createObservabilityPackDashboard(ctx, dashboard); err != nil {
					return fmt.Errorf("Failed to create observability pack dashboard [%s]: %s", dashboard.Name, err)
				}
			}
		}
	}

	p.progressIndicator.Success(msg)
	return nil
}

func (p *ConcretePacksInstaller) createObservabilityPackDashboard(ctx context.Context, d types.OpenInstallationObservabilityPackDashboard) (*dashboards.DashboardCreateResult, error) {
	// TODO: check for existance of dashboard before creating a new one

	defaultProfile := credentials.DefaultProfile()
	accountId := defaultProfile.AccountID

	body, err := getJSONfromURL(d.URL)
	if err != nil {
		return nil, err
	}

	dashboard, err := transformDashboardJSON(body, defaultProfile.AccountID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("  Creating dashboard: %s\n", dashboard.Name)
	created, err := p.client.Dashboards.DashboardCreate(accountId, dashboard)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func getJSONfromURL(url string) ([]byte, error) {
	response, err := defaultHTTPGetFunc(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("received non-2xx status code %d when retrieving recipe", response.StatusCode)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Tracef("Resp from Url[%s]: %s", url, string(body))

	return body, nil
}

func transformDashboardJSON(body []byte, accountId int) (dashboards.DashboardInput, error) {
	dashboard := dashboards.DashboardInput{}
	re := regexp.MustCompile("\"accountId\": 0")
	dash := re.ReplaceAllString(string(body), fmt.Sprintf("\"accountId\": %d", accountId))

	if err := json.Unmarshal([]byte(dash), &dashboard); err != nil {
		return dashboard, err
	}

	dashboard.Permissions = entities.DashboardPermissionsTypes.PUBLIC_READ_WRITE
	log.Tracef("Dashboard definition: %+v", dashboard)

	return dashboard, nil
}

const (
	createDashboardMutation = `
	mutation ($accountId: Int!, $dashboard: DashboardInput!) {
		dashboardCreate(accountId: $accountId, dashboard: $dashboard) {
			errors {
				description
				type
			}
			entityResult {
				guid
			}
		}
	}
	`
)
