package packs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"

	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

type ServicePacksInstaller struct {
	client            *newrelic.NewRelic
	progressIndicator ux.ProgressIndicator
	installStatus     *execution.InstallStatus
}

func NewServicePacksInstaller(client *newrelic.NewRelic, s *execution.InstallStatus) *ServicePacksInstaller {
	return &ServicePacksInstaller{
		client:            client,
		progressIndicator: ux.NewPlainProgress(),
		installStatus:     s,
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
		p.installStatus.ObservabilityPackInstallPending(execution.ObservabilityPackStatusEvent{ObservabilityPack: pack})
		fmt.Printf("\n  Installing pack: %s\n", pack.Name)

		// Only installing dashboards currently
		// failure to create any dashboards results in the whole pack installation being marked as failure
		if pack.Dashboards != nil {
			for _, dashboard := range pack.Dashboards {
				if _, err := p.createObservabilityPackDashboard(ctx, dashboard); err != nil {
					errMsg := fmt.Sprintf("failed to create observability pack dashboard [%s]: %s", dashboard.Name, err)

					p.installStatus.ObservabilityPackInstallFailed(execution.ObservabilityPackStatusEvent{ObservabilityPack: pack, Msg: errMsg})
					p.progressIndicator.Fail(msg)

					return fmt.Errorf(errMsg)
				}
			}
		}

		p.installStatus.ObservabilityPackInstallSuccess(execution.ObservabilityPackStatusEvent{ObservabilityPack: pack})
	}

	p.progressIndicator.Success(msg)
	return nil
}

func (p *ServicePacksInstaller) createObservabilityPackDashboard(ctx context.Context, d types.OpenInstallationObservabilityPackDashboard) (*dashboards.DashboardCreateResult, error) {
	accountID := configAPI.GetActiveProfileAccountID()

	body, err := getJSONfromURL(d.URL)
	if err != nil {
		return nil, err
	}

	dashboard, err := transformDashboardJSON(body, accountID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("  ==> Creating dashboard: %s\n", dashboard.Name)

	entitySearchResult, err := p.client.Entities.GetEntitySearchWithContext(
		ctx,
		entities.EntitySearchOptions{},
		"",
		entities.EntitySearchQueryBuilder{
			Name: dashboard.Name,
			Type: entities.EntitySearchQueryBuilderTypeTypes.DASHBOARD,
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error checking if dashboard already exists: %s", err)
	}

	if len(entitySearchResult.Results.Entities) > 0 {
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
