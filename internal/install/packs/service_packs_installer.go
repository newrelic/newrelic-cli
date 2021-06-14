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

// nolint: golint
type PacksInstaller interface {
	Install(ctx context.Context, packs []types.OpenInstallationObservabilityPack) error
}

type ServicePacksInstaller struct {
	client            *newrelic.NewRelic
	progressIndicator ux.ProgressIndicator
}

func NewServicePacksInstaller(client *newrelic.NewRelic) *ServicePacksInstaller {
	return &ServicePacksInstaller{
		client:            client,
		progressIndicator: ux.NewPlainProgress(),
	}
}

func defaultHTTPGetFunc(dashboardURL string) (*http.Response, error) {
	return http.Get(dashboardURL)
}

func (p *ServicePacksInstaller) Install(ctx context.Context, packs []types.OpenInstallationObservabilityPack) error {
	msg := "Installing observability packs"
	p.progressIndicator.Start(msg)
	defer p.progressIndicator.Stop()

	for _, pack := range packs {
		fmt.Printf("\n  Installing pack: %s\n", pack.Name)

		// Only installing dashboards currently
		if pack.Dashboards != nil {
			for _, dashboard := range pack.Dashboards {
				if _, err := p.createObservabilityPackDashboard(ctx, dashboard); err != nil {
					p.progressIndicator.Fail(msg)
					return fmt.Errorf("Failed to create observability pack dashboard [%s]: %s", dashboard.Name, err) // nolint: golint
				}
			}
		}
	}

	p.progressIndicator.Success(msg)
	return nil
}

func (p *ServicePacksInstaller) createObservabilityPackDashboard(ctx context.Context, d types.OpenInstallationObservabilityPackDashboard) (*dashboards.DashboardCreateResult, error) {
	defaultProfile := credentials.DefaultProfile()
	accountID := defaultProfile.AccountID

	body, err := getJSONfromURL(d.URL)
	if err != nil {
		return nil, err
	}

	dashboard, err := transformDashboardJSON(body, defaultProfile.AccountID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("  ==> Creating dashboard: %s\n", dashboard.Name)

	// Check for existence of dashboard before creating a new one
	dashboards, err := p.client.Dashboards.ListDashboards(&dashboards.ListDashboardsParams{
		Title: dashboard.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("error checking if dashboard already exists: %s", err)
	}

	if len(dashboards) > 0 {
		fmt.Printf("  ==> Dashboard [%s] already exists, skipping\n", dashboard.Name)
		return nil, nil
	}

	// Dashboard doesn't exist yet, proceed with dashboard create
	created, err := p.client.Dashboards.DashboardCreateWithContext(ctx, accountID, dashboard)
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

func transformDashboardJSON(body []byte, accountID int) (dashboards.DashboardInput, error) {
	dashboard := dashboards.DashboardInput{}
	re := regexp.MustCompile("\"accountId\": 0")
	dash := re.ReplaceAllString(string(body), fmt.Sprintf("\"accountId\": %d", accountID))

	if err := json.Unmarshal([]byte(dash), &dashboard); err != nil {
		return dashboard, err
	}

	dashboard.Permissions = entities.DashboardPermissionsTypes.PUBLIC_READ_WRITE
	log.Tracef("Dashboard definition: %+v", dashboard)

	return dashboard, nil
}
