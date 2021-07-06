package execution

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http/httpproxy"

	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// nolint: maligned
type InstallStatus struct {
	Complete                  bool                       `json:"complete"`
	DiscoveryManifest         types.DiscoveryManifest    `json:"discoveryManifest"`
	EntityGUIDs               []string                   `json:"entityGuids"`
	Error                     StatusError                `json:"error"`
	LogFilePath               string                     `json:"logFilePath"`
	Statuses                  []*RecipeStatus            `json:"recipes"`
	ObservabilityPackStatuses []*ObservabilityPackStatus `json:"packs"`
	Timestamp                 int64                      `json:"timestamp"`
	CLIVersion                string                     `json:"cliVersion"`
	HasInstalledRecipes       bool                       `json:"hasInstalledRecipes"`
	HasCanceledRecipes        bool                       `json:"hasCanceledRecipes"`
	HasSkippedRecipes         bool                       `json:"hasSkippedRecipes"`
	HasFailedRecipes          bool                       `json:"hasFailedRecipes"`
	HasUnsupportedRecipes     bool                       `json:"hasUnsupportedRecipes"`
	HasInstalledPacks         bool                       `json:"hasInstalledPacks"`
	HasCanceledPacks          bool                       `json:"hasCanceledPacks"`
	HasFailedPacks            bool                       `json:"hasFailedPacks"`
	Skipped                   []*RecipeStatus            `json:"recipesSkipped"`
	Canceled                  []*RecipeStatus            `json:"recipesCanceled"`
	Failed                    []*RecipeStatus            `json:"recipesFailed"`
	Installed                 []*RecipeStatus            `json:"recipesInstalled"`
	CanceledPacks             []*ObservabilityPackStatus `json:"packsCanceled"`
	FailedPacks               []*ObservabilityPackStatus `json:"packsFailed"`
	InstalledPacks            []*ObservabilityPackStatus `json:"packslInstalled"`
	RedirectURL               string                     `json:"redirectUrl"`
	HTTPSProxy                string                     `json:"httpsProxy"`
	DocumentID                string
	targetedInstall           bool
	statusSubscriber          []StatusSubscriber
	successLinkConfig         types.OpenInstallationSuccessLinkConfig
	PlatformLinkGenerator     LinkGenerator
}

type RecipeStatus struct {
	DisplayName string           `json:"displayName"`
	Error       StatusError      `json:"error"`
	Name        string           `json:"name"`
	Status      RecipeStatusType `json:"status"`
	EntityGUID  string           `json:"entityGuid,omitempty"`
	// ValidationDurationMilliseconds is duration in Milliseconds that a recipe took to validate data was flowing.
	ValidationDurationMilliseconds int64 `json:"validationDurationMilliseconds,omitempty"`
}

type RecipeStatusType string

var RecipeStatusTypes = struct {
	AVAILABLE   RecipeStatusType
	CANCELED    RecipeStatusType
	INSTALLING  RecipeStatusType
	FAILED      RecipeStatusType
	INSTALLED   RecipeStatusType
	SKIPPED     RecipeStatusType
	RECOMMENDED RecipeStatusType
	UNSUPPORTED RecipeStatusType
}{
	AVAILABLE:   "AVAILABLE",
	CANCELED:    "CANCELED",
	INSTALLING:  "INSTALLING",
	FAILED:      "FAILED",
	INSTALLED:   "INSTALLED",
	SKIPPED:     "SKIPPED",
	RECOMMENDED: "RECOMMENDED",
	UNSUPPORTED: "UNSUPPORTED",
}

type ObservabilityPackStatus struct {
	Error  StatusError                 `json:"error"`
	Name   string                      `json:"name"`
	Status ObservabilityPackStatusType `json:"status"`
}

type ObservabilityPackStatusType string

var ObservabilityPackStatusTypes = struct {
	FetchPending   ObservabilityPackStatusType
	FetchSuccess   ObservabilityPackStatusType
	FetchFailed    ObservabilityPackStatusType
	InstallPending ObservabilityPackStatusType
	InstallSuccess ObservabilityPackStatusType
	InstallFailed  ObservabilityPackStatusType
	Canceled       ObservabilityPackStatusType
}{
	FetchPending:   "FETCH_PENDING",
	FetchSuccess:   "FETCH_SUCCESS",
	FetchFailed:    "FETCH_FAILED",
	InstallPending: "INSTALL_PENDING",
	InstallSuccess: "INSTALL_SUCCESS",
	InstallFailed:  "INSTALL_FAILED",
	Canceled:       "CANCELED",
}

type StatusError struct {
	Message  string   `json:"message"`
	Details  string   `json:"details"`
	TaskPath []string `json:"taskPath"`
}

func NewInstallStatus(reporters []StatusSubscriber, PlatformLinkGenerator LinkGenerator) *InstallStatus {
	s := InstallStatus{
		DocumentID:            uuid.New().String(),
		Timestamp:             utils.GetTimestamp(),
		LogFilePath:           configuration.BasePath + "/" + configuration.DefaultLogFile,
		statusSubscriber:      reporters,
		PlatformLinkGenerator: PlatformLinkGenerator,
		HTTPSProxy:            httpproxy.FromEnvironment().HTTPSProxy,
	}

	return &s
}

func (s *InstallStatus) DiscoveryComplete(dm types.DiscoveryManifest) {
	s.withDiscoveryInfo(dm)

	for _, r := range s.statusSubscriber {
		if err := r.DiscoveryComplete(s, dm); err != nil {
			log.Errorf("Could not report discovery info: %s", err)
		}
	}
}

func (s *InstallStatus) RecipeAvailable(recipe types.OpenInstallationRecipe) {
	s.withAvailableRecipe(recipe)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeAvailable(s, recipe); err != nil {
			log.Errorf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *InstallStatus) RecipesSelected(recipes []types.OpenInstallationRecipe) {
	for _, r := range s.statusSubscriber {
		if err := r.RecipesSelected(s, recipes); err != nil {
			log.Errorf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *InstallStatus) ObservabilityPackFetchPending(event ObservabilityPackStatusEvent) {
	s.withObservabilityPackEvent(event, ObservabilityPackStatusTypes.FetchPending)

	for _, r := range s.statusSubscriber {
		if err := r.ObservabilityPackFetchPending(s); err != nil {
			log.Errorf("Error writing observabilityPack status for pack %s: %s", event.ObservabilityPack.Name, err)
		}
	}
}

func (s *InstallStatus) ObservabilityPackFetchSuccess(event ObservabilityPackStatusEvent) {
	s.withObservabilityPackEvent(event, ObservabilityPackStatusTypes.FetchSuccess)

	for _, r := range s.statusSubscriber {
		if err := r.ObservabilityPackFetchSuccess(s); err != nil {
			log.Errorf("Error writing observabilityPack status for pack %s: %s", event.ObservabilityPack.Name, err)
		}
	}
}

func (s *InstallStatus) ObservabilityPackFetchFailed(event ObservabilityPackStatusEvent) {
	s.withObservabilityPackEvent(event, ObservabilityPackStatusTypes.FetchFailed)

	for _, r := range s.statusSubscriber {
		if err := r.ObservabilityPackFetchFailed(s); err != nil {
			log.Errorf("Error writing observabilityPack status for pack %s: %s", event.ObservabilityPack.Name, err)
		}
	}
}

func (s *InstallStatus) ObservabilityPackInstallPending(event ObservabilityPackStatusEvent) {
	s.withObservabilityPackEvent(event, ObservabilityPackStatusTypes.InstallPending)

	for _, r := range s.statusSubscriber {
		if err := r.ObservabilityPackInstallPending(s); err != nil {
			log.Errorf("Error writing observabilityPack status for pack %s: %s", event.ObservabilityPack.Name, err)
		}
	}
}

func (s *InstallStatus) ObservabilityPackInstallSuccess(event ObservabilityPackStatusEvent) {
	s.withObservabilityPackEvent(event, ObservabilityPackStatusTypes.InstallSuccess)

	for _, r := range s.statusSubscriber {
		if err := r.ObservabilityPackInstallSuccess(s); err != nil {
			log.Errorf("Error writing observabilityPack status for pack %s: %s", event.ObservabilityPack.Name, err)
		}
	}
}

func (s *InstallStatus) ObservabilityPackInstallFailed(event ObservabilityPackStatusEvent) {
	s.withObservabilityPackEvent(event, ObservabilityPackStatusTypes.InstallFailed)

	for _, r := range s.statusSubscriber {
		if err := r.ObservabilityPackInstallFailed(s); err != nil {
			log.Errorf("Error writing observabilityPack status for pack %s: %s", event.ObservabilityPack.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeInstalled(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.INSTALLED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeInstalled(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

// RecipeRecommended is responsible for setting the nerstorage scopes
// when a recipe is recommended.  This is used when a recipe is found, but not
// a "HOST" type, and is used to indicate to the user that it is something they
// should consider integrating, but not something that the recipe framework
// will currently assist with.
func (s *InstallStatus) RecipeRecommended(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.RECOMMENDED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeRecommended(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeInstalling(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.INSTALLING)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeInstalling(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeFailed(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.FAILED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeFailed(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeSkipped(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.SKIPPED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeSkipped(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeUnsupported(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.UNSUPPORTED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeUnsupported(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) InstallComplete(err error) {
	s.completed(err)

	for _, r := range s.statusSubscriber {
		if err := r.InstallComplete(s); err != nil {
			log.Errorf("Error writing execution status: %s", err)
		}
	}
}

func (s *InstallStatus) InstallCanceled() {
	s.canceled()

	for _, r := range s.statusSubscriber {
		if err := r.InstallCanceled(s); err != nil {
			log.Errorf("Error writing execution status: %s", err)
		}
	}
}

func (s *InstallStatus) WasSuccessful() bool {
	return s.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED)
}

func (s *InstallStatus) recommendations() []*RecipeStatus {
	var statuses []*RecipeStatus

	for _, st := range s.Statuses {
		if st.Status == RecipeStatusTypes.RECOMMENDED {
			statuses = append(statuses, st)
		}
	}

	return statuses
}

func (s *InstallStatus) hasAnyRecipeStatus(status RecipeStatusType) bool {
	for _, ss := range s.Statuses {
		if ss.Status == status {
			return true
		}
	}

	return false
}

func (s *InstallStatus) SetTargetedInstall() {
	s.targetedInstall = true
}

func (s *InstallStatus) IsTargetedInstall() bool {
	return s.targetedInstall
}

func (s *InstallStatus) HostEntityGUID() string {
	var guid string

	// When we have performed a targeted installation, we want to roll up to the last GUID in the list.
	if len(s.EntityGUIDs) > 0 {
		if s.IsTargetedInstall() {
			guid = s.EntityGUIDs[len(s.EntityGUIDs)-1]
		} else {
			guid = s.EntityGUIDs[0]
		}
	}

	return guid
}

func (s *InstallStatus) setRedirectURL() {
	s.RedirectURL = s.PlatformLinkGenerator.GenerateRedirectURL(*s)
}

func (s *InstallStatus) withAvailableRecipes(recipes []types.OpenInstallationRecipe) {
	for _, r := range recipes {
		s.withAvailableRecipe(r)
	}
}

func (s *InstallStatus) withAvailableRecipe(r types.OpenInstallationRecipe) {
	e := RecipeStatusEvent{Recipe: r}
	s.withRecipeEvent(e, RecipeStatusTypes.AVAILABLE)
}

func (s *InstallStatus) withSuccessLinkConfig(l types.OpenInstallationSuccessLinkConfig) {
	s.successLinkConfig = l
}

func (s *InstallStatus) withEntityGUID(entityGUID string) {
	for _, e := range s.EntityGUIDs {
		if e == entityGUID {
			return
		}
	}

	log.WithFields(log.Fields{
		"guid": entityGUID,
	}).Debug("new GUID")

	s.EntityGUIDs = append(s.EntityGUIDs, entityGUID)
}

func (s *InstallStatus) withDiscoveryInfo(dm types.DiscoveryManifest) {
	s.DiscoveryManifest = dm
	s.Timestamp = utils.GetTimestamp()

	version := os.Getenv("NEW_RELIC_CLI_VERSION")
	if version != "" {
		s.CLIVersion = version
	}
}

func (s *InstallStatus) withObservabilityPackEvent(e ObservabilityPackStatusEvent, opst ObservabilityPackStatusType) {
	statusError := StatusError{
		Message: e.Msg,
	}

	var name string
	if e.Name != "" {
		name = e.Name
	} else {
		name = e.ObservabilityPack.Name
	}

	// Not using this logic for now: these events are sent too
	// quick for the UI to keep up when we modify an existing status.
	// Instead, we'll now send a list of events
	//
	// We can switch back to using this once the install-events-service
	// is in place
	//
	// found := s.getObservabilityPackStatusByPackName(name)
	//
	// if found != nil {
	// 	found.Status = opst
	//
	// 	if e.Msg != "" {
	// 		found.Error = statusError
	// 	}
	// } else {
	// observabilityPackStatus := &ObservabilityPackStatus{
	// 	Name:   name,
	// 	Error:  statusError,
	// 	Status: opst,
	// }
	// s.ObservabilityPackStatuses = append(s.ObservabilityPackStatuses, observabilityPackStatus)
	// }

	observabilityPackStatus := &ObservabilityPackStatus{
		Name:   name,
		Error:  statusError,
		Status: opst,
	}
	s.ObservabilityPackStatuses = append(s.ObservabilityPackStatuses, observabilityPackStatus)

	log.WithFields(log.Fields{
		"observabilityPack_name": e.ObservabilityPack.Name,
		"status":                 opst,
		"error":                  statusError.Message,
		"statusCount":            len(s.Statuses),
	}).Debug("observabilityPack event")
}

func (s *InstallStatus) withRecipeEvent(e RecipeStatusEvent, rs RecipeStatusType) {
	if e.EntityGUID != "" {
		s.withEntityGUID(e.EntityGUID)
	}

	s.withSuccessLinkConfig(e.Recipe.SuccessLinkConfig)

	statusError := StatusError{
		Message:  e.Msg,
		TaskPath: e.TaskPath,
	}

	s.Error = statusError

	found := s.getStatus(e.Recipe)

	if found != nil {
		found.Status = rs

		if e.EntityGUID != "" {
			found.EntityGUID = e.EntityGUID
		}

		if e.ValidationDurationMilliseconds > 0 {
			found.ValidationDurationMilliseconds = e.ValidationDurationMilliseconds
		}

		if e.Msg != "" {
			found.Error = statusError
		}
	} else {
		recipeStatus := &RecipeStatus{
			Name:        e.Recipe.Name,
			DisplayName: e.Recipe.DisplayName,
			Status:      rs,
			Error:       statusError,
		}

		if e.EntityGUID != "" {
			recipeStatus.EntityGUID = e.EntityGUID
		}

		if e.ValidationDurationMilliseconds > 0 {
			recipeStatus.ValidationDurationMilliseconds = e.ValidationDurationMilliseconds
		}

		s.Statuses = append(s.Statuses, recipeStatus)
	}

	s.Timestamp = utils.GetTimestamp()

	log.WithFields(log.Fields{
		"recipe_name":                    e.Recipe.Name,
		"status":                         rs,
		"error":                          statusError.Message,
		"tasks":                          statusError.TaskPath,
		"guid":                           e.EntityGUID,
		"validationDurationMilliseconds": e.ValidationDurationMilliseconds,
		"statusCount":                    len(s.Statuses),
	}).Debug("recipe event")
}

func (s *InstallStatus) completed(err error) {
	isUnsupported := false
	s.Complete = true
	s.Timestamp = utils.GetTimestamp()

	if err != nil {
		statusError := StatusError{
			Message: err.Error(),
		}

		if e, ok := err.(types.GoTaskError); ok {
			statusError.TaskPath = e.TaskPath()
		}

		if _, ok := err.(*types.UnsupportedOperatingSytemError); ok {
			isUnsupported = true
		}

		s.Error = statusError
	}

	log.WithFields(log.Fields{
		"timestamp": s.Timestamp,
	}).Debug("completed")

	s.updateFinalInstallationStatuses(false, isUnsupported)
	s.setRedirectURL()
}

func (s *InstallStatus) canceled() {
	s.Timestamp = utils.GetTimestamp()

	log.WithFields(log.Fields{
		"timestamp": s.Timestamp,
	}).Debug("canceled")

	s.updateFinalInstallationStatuses(true, false)
}

func (s *InstallStatus) getStatus(r types.OpenInstallationRecipe) *RecipeStatus {
	for _, recipe := range s.Statuses {
		if recipe.Name == r.Name {
			return recipe
		}
	}

	return nil
}

// This is unused for now: these events are sent too
// quick for the UI to keep up when we modify an existing status.
// Instead, we'll now send a list of events
//
// We can switch back to using this once the install-events-service
// is in place
// func (s *InstallStatus) getObservabilityPackStatusByPackName(name string) *ObservabilityPackStatus {
// 	for _, pack := range s.ObservabilityPackStatuses {
// 		if pack.Name == name {
// 			return pack
// 		}
// 	}

// 	return nil
// }

func (s *InstallStatus) getObservabilityPackStatusByPackStatusType(st ObservabilityPackStatusType) *ObservabilityPackStatus {
	for _, pack := range s.ObservabilityPackStatuses {
		if pack.Status == st {
			return pack
		}
	}

	return nil
}

// This function handles updating the final recipe statuses and top-level installation status.
// Canceling (e.g. ctl+c) will cause unresolved recipes to be marked as canceled.
// Exiting early (i.e. an error occurred) will cause unresolved recipes to be marked as failed.
func (s *InstallStatus) updateFinalInstallationStatuses(installCanceled bool, isUnsupported bool) {
	s.updateRecipeStatuses(installCanceled, isUnsupported)
	packs := s.collectStatuses()
	s.updateObservabililtyPackStatuses(packs, installCanceled)

	log.WithFields(log.Fields{
		"hasInstalledRecipes": s.HasInstalledRecipes,
		"hasSkippedRecipes":   s.HasSkippedRecipes,
		"hasCanceledRecipes":  s.HasCanceledRecipes,
		"hasFailedRecipes":    s.HasFailedRecipes,
		"hasInstalledPacks":   s.HasInstalledPacks,
		"hasCanceledPacks":    s.HasCanceledPacks,
		"hasFailedPacks":      s.HasFailedPacks,
	}).Debug("final installation statuses updated")
}

func (s *InstallStatus) updateRecipeStatuses(installCanceled bool, isUnsupported bool) {
	for i, ss := range s.Statuses {
		if ss.Status == RecipeStatusTypes.AVAILABLE || ss.Status == RecipeStatusTypes.INSTALLING {
			debugMsg := "failed"

			if installCanceled {
				debugMsg = "canceled"
			}

			if isUnsupported {
				debugMsg = "unsupported"
			}

			log.WithFields(log.Fields{
				"recipe": s.Statuses[i].Name,
			}).Debug(fmt.Sprintf("marking recipe %s", debugMsg))

			if installCanceled {
				s.Statuses[i].Status = RecipeStatusTypes.CANCELED
			} else {
				if ss.Status == RecipeStatusTypes.INSTALLING {
					s.Statuses[i].Status = RecipeStatusTypes.FAILED
				}
			}
		}

		// Installed
		if ss.Status == RecipeStatusTypes.INSTALLED {
			s.Installed = append(s.Installed, ss)
			s.HasInstalledRecipes = true
		}

		// Skipped
		if ss.Status == RecipeStatusTypes.SKIPPED {
			s.Skipped = append(s.Skipped, ss)
			s.HasSkippedRecipes = true
		}

		// Canceled
		if ss.Status == RecipeStatusTypes.CANCELED {
			s.Canceled = append(s.Canceled, ss)
			s.HasCanceledRecipes = true
		}

		// Errored
		if ss.Status == RecipeStatusTypes.FAILED {
			s.Failed = append(s.Failed, ss)
			s.HasFailedRecipes = true
		}

		// Unsupported
		if ss.Status == RecipeStatusTypes.UNSUPPORTED {
			s.HasUnsupportedRecipes = true
		}
	}
}

func (s *InstallStatus) updateObservabililtyPackStatuses(packs map[string][]ObservabilityPackStatusType, installCanceled bool) {
	for i, ops := range s.ObservabilityPackStatuses {
		// Compare ops.Status w/ the last known status
		// If they're the same && not FETCH_FAILED/INSTALL_SUCCESS/INSTALL_FAILED (these are final statuses), update to CANCELED/INSTALL_FAILED
		if v, ok := packs[ops.Name]; ok {
			lastStatus := v[len(v)-1]

			if lastStatus == ops.Status && (lastStatus != ObservabilityPackStatusTypes.FetchFailed &&
				lastStatus != ObservabilityPackStatusTypes.InstallSuccess &&
				lastStatus != ObservabilityPackStatusTypes.InstallFailed) {
				debugMsg := "failed"

				if installCanceled {
					debugMsg = "canceled"
				}

				log.WithFields(log.Fields{
					"lastStatus":        lastStatus,
					"observabilityPack": s.ObservabilityPackStatuses[i].Name,
				}).Debug(fmt.Sprintf("marking observabilityPack %s", debugMsg))

				if installCanceled {
					s.ObservabilityPackStatuses[i].Status = ObservabilityPackStatusTypes.Canceled
				} else {
					s.ObservabilityPackStatuses[i].Status = ObservabilityPackStatusTypes.InstallFailed
				}
			}
		}

		// Report out the final statuses
		// Installed
		if ops.Status == ObservabilityPackStatusTypes.InstallSuccess {
			s.InstalledPacks = append(s.InstalledPacks, ops)
			s.HasInstalledPacks = true
		}

		// Canceled
		if ops.Status == ObservabilityPackStatusTypes.Canceled {
			s.CanceledPacks = append(s.CanceledPacks, ops)
			s.HasCanceledPacks = true
		}

		// Errored
		if ops.Status == ObservabilityPackStatusTypes.InstallFailed {
			s.FailedPacks = append(s.FailedPacks, ops)
			s.HasFailedPacks = true
		}
	}
}

/**
 * Collect every pack's status in a map for ease of deciding the final state
 * of a given pack
 */
func (s *InstallStatus) collectStatuses() map[string][]ObservabilityPackStatusType {
	res := map[string][]ObservabilityPackStatusType{}

	for _, s := range s.ObservabilityPackStatuses {
		if v, ok := res[s.Name]; ok {
			res[s.Name] = append(v, s.Status)
		} else {
			res[s.Name] = []ObservabilityPackStatusType{
				s.Status,
			}
		}
	}
	log.Tracef("[InstallStatus.collectStatuses]: %+v", res)
	return res
}
