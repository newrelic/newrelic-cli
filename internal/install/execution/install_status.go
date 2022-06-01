package execution

import (
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http/httpproxy"

	"github.com/newrelic/newrelic-cli/internal/cli"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// nolint: maligned
type InstallStatus struct {
	InstallID             string                  `json:"installId"`
	Complete              bool                    `json:"complete"`
	DiscoveryManifest     types.DiscoveryManifest `json:"discoveryManifest"`
	EntityGUIDs           []string                `json:"entityGuids"`
	Error                 StatusError             `json:"error"`
	LogFilePath           string                  `json:"logFilePath"`
	Statuses              []*RecipeStatus         `json:"recipes"`
	Timestamp             int64                   `json:"timestamp"`
	CLIVersion            string                  `json:"cliVersion"`
	InstallLibraryVersion string                  `json:"installLibraryVersion"`
	HasInstalledRecipes   bool                    `json:"hasInstalledRecipes"`
	HasCanceledRecipes    bool                    `json:"hasCanceledRecipes"`
	HasSkippedRecipes     bool                    `json:"hasSkippedRecipes"`
	HasFailedRecipes      bool                    `json:"hasFailedRecipes"`
	HasUnsupportedRecipes bool                    `json:"hasUnsupportedRecipes"`
	Skipped               []*RecipeStatus         `json:"recipesSkipped"`
	Canceled              []*RecipeStatus         `json:"recipesCanceled"`
	Failed                []*RecipeStatus         `json:"recipesFailed"`
	Installed             []*RecipeStatus         `json:"recipesInstalled"`
	RedirectURL           string                  `json:"redirectUrl"`
	HTTPSProxy            string                  `json:"httpsProxy"`
	UpdateRequired        bool                    `json:"updateRequired"`
	DocumentID            string
	targetedInstall       bool
	statusSubscriber      []StatusSubscriber
	successLinkConfig     types.OpenInstallationSuccessLinkConfig
	PlatformLinkGenerator LinkGenerator
}

type RecipeStatus struct {
	DisplayName string           `json:"displayName"`
	Error       StatusError      `json:"error"`
	Name        string           `json:"name"`
	Status      RecipeStatusType `json:"status"`
	EntityGUID  string           `json:"entityGuid,omitempty"`
	// validationDurationMs is duration in Milliseconds that a recipe took to validate data was flowing.
	ValidationDurationMs int64 `json:"validationDurationMs,omitempty"`
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
	DETECTED    RecipeStatusType
	NULL        RecipeStatusType
}{
	AVAILABLE:   "AVAILABLE",
	CANCELED:    "CANCELED",
	INSTALLING:  "INSTALLING",
	FAILED:      "FAILED",
	INSTALLED:   "INSTALLED",
	SKIPPED:     "SKIPPED",
	RECOMMENDED: "RECOMMENDED",
	UNSUPPORTED: "UNSUPPORTED",
	DETECTED:    "DETECTED",
	NULL:        "",
}

type StatusError struct {
	Message  string   `json:"message"`
	Details  string   `json:"details"`
	TaskPath []string `json:"taskPath"`
}

var StatusIconMap = map[RecipeStatusType]string{
	RecipeStatusTypes.INSTALLED:   ux.IconSuccess,
	RecipeStatusTypes.FAILED:      ux.IconError,
	RecipeStatusTypes.UNSUPPORTED: ux.IconUnsupported,
	RecipeStatusTypes.SKIPPED:     ux.IconMinus,
	RecipeStatusTypes.CANCELED:    ux.IconMinus,
}

func NewInstallStatus(reporters []StatusSubscriber, PlatformLinkGenerator LinkGenerator) *InstallStatus {
	s := InstallStatus{
		InstallID:             uuid.New().String(),
		DocumentID:            uuid.New().String(),
		Timestamp:             utils.GetTimestamp(),
		LogFilePath:           config.GetDefaultLogFilePath(),
		statusSubscriber:      reporters,
		PlatformLinkGenerator: PlatformLinkGenerator,
		HTTPSProxy:            httpproxy.FromEnvironment().HTTPSProxy,
		CLIVersion:            cli.Version(),
	}

	return &s
}

func (s *InstallStatus) DiscoveryComplete(dm types.DiscoveryManifest) {
	s.withDiscoveryInfo(dm)

	for _, r := range s.statusSubscriber {
		if err := r.DiscoveryComplete(s, dm); err != nil {
			log.Debugf("Could not report discovery info: %s", err)
		}
	}
}

// RecipeDetected handles setting the DETECTED status for the provided
// recipe as well sending out the DETECTED status event to the install events service.
// RecipeDetected is called when a recipe is available and passes the checks in both
// the process match and the pre-install steps of recipe execution.
func (s *InstallStatus) RecipeDetected(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.DETECTED)
	for _, r := range s.statusSubscriber {
		if err := r.RecipeDetected(s, event); err != nil {
			log.Debugf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *InstallStatus) RecipeCanceled(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.CANCELED)
	for _, r := range s.statusSubscriber {
		if err := r.RecipeCanceled(s, event); err != nil {
			log.Debugf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *InstallStatus) RecipeAvailable(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.AVAILABLE)
	for _, ss := range s.statusSubscriber {
		if err := ss.RecipeAvailable(s, event); err != nil {
			log.Debugf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *InstallStatus) RecipeInstalled(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.INSTALLED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeInstalled(s, event); err != nil {
			log.Debugf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
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
			log.Debugf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeInstalling(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.INSTALLING)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeInstalling(s, event); err != nil {
			log.Debugf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeFailed(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.FAILED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeFailed(s, event); err != nil {
			log.Debugf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeSkipped(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.SKIPPED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeSkipped(s, event); err != nil {
			log.Debugf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeUnsupported(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.UNSUPPORTED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeUnsupported(s, event); err != nil {
			log.Debugf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) InstallStarted() {
	s.started()

	for _, r := range s.statusSubscriber {
		if err := r.InstallStarted(s); err != nil {
			log.Debugf("Error writing execution status: %s", err)
		}
	}
}

func (s *InstallStatus) InstallComplete(err error) {
	s.completed(err)

	for _, r := range s.statusSubscriber {
		if err := r.InstallComplete(s); err != nil {
			log.Debugf("Error writing execution status: %s", err)
		}
	}
}

func (s *InstallStatus) RecipeHasStatus(recipeName string, status RecipeStatusType) bool {
	for _, s := range s.Statuses {
		if s.Name == recipeName && s.Status == status {
			return true
		}
	}
	return false
}

func (s *InstallStatus) ReportStatus(status RecipeStatusType, event RecipeStatusEvent) {

	switch status {
	case RecipeStatusTypes.AVAILABLE:
		s.RecipeAvailable(event)
	case RecipeStatusTypes.CANCELED:
		s.RecipeCanceled(event)
	case RecipeStatusTypes.DETECTED:
		s.RecipeDetected(event)
	case RecipeStatusTypes.FAILED:
		s.RecipeFailed(event)
	case RecipeStatusTypes.INSTALLED:
		s.RecipeInstalled(event)
	case RecipeStatusTypes.INSTALLING:
		s.RecipeInstalling(event)
	case RecipeStatusTypes.SKIPPED:
		s.RecipeSkipped(event)
	case RecipeStatusTypes.UNSUPPORTED:
		s.RecipeUnsupported(event)
	case RecipeStatusTypes.RECOMMENDED:
		s.RecipeRecommended(event)
	case RecipeStatusTypes.NULL:
		// Not used
	default:
		log.Warnf("Unknown status to report: %s, ignoring", status)
	}
}

func (s *InstallStatus) InstallCanceled() {
	s.canceled()

	for _, r := range s.statusSubscriber {
		if err := r.InstallCanceled(s); err != nil {
			log.Debugf("Error writing execution status: %s", err)
		}
	}
}

func (s *InstallStatus) WasSuccessful() bool {
	return s.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED)
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

func (s *InstallStatus) SetVersions(installLibraryVersion string) {
	s.InstallLibraryVersion = installLibraryVersion
	s.CLIVersion = cli.Version()
}

func (s *InstallStatus) withDiscoveryInfo(dm types.DiscoveryManifest) {
	s.DiscoveryManifest = dm
	s.Timestamp = utils.GetTimestamp()
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

		if e.ValidationDurationMs > 0 {
			found.ValidationDurationMs = e.ValidationDurationMs
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

		if e.ValidationDurationMs > 0 {
			recipeStatus.ValidationDurationMs = e.ValidationDurationMs
		}

		s.Statuses = append(s.Statuses, recipeStatus)
	}

	s.Timestamp = utils.GetTimestamp()

	log.WithFields(log.Fields{
		"recipe_name":          e.Recipe.Name,
		"status":               rs,
		"error":                statusError.Message,
		"tasks":                statusError.TaskPath,
		"guid":                 e.EntityGUID,
		"validationDurationMs": e.ValidationDurationMs,
		"statusCount":          len(s.Statuses),
	}).Debug("recipe event")
}

func (s *InstallStatus) started() {
	s.Timestamp = utils.GetTimestamp()

	log.WithFields(log.Fields{
		"timestamp": s.Timestamp,
	}).Debug("started")
}

func (s *InstallStatus) completed(err error) {
	isUnsupported := false
	s.Complete = true
	s.Timestamp = utils.GetTimestamp()

	if err != nil {
		statusError := StatusError{
			Message: err.Error(),
		}

		if e, ok := err.(*types.UpdateRequiredError); ok {
			statusError.Details = e.Details
			s.CLIVersion = cli.Version()
		}

		if e, ok := err.(types.GoTaskError); ok {
			statusError.TaskPath = e.TaskPath()
		}

		if _, ok := err.(*types.UnsupportedOperatingSystemError); ok {
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
	s.setRedirectURL()
}

func (s *InstallStatus) getStatus(r types.OpenInstallationRecipe) *RecipeStatus {
	for _, rs := range s.Statuses {
		if rs.Name == r.Name {
			return rs
		}
	}

	return nil
}

// This function handles updating the final recipe statuses and top-level installation status.
// Canceling (e.g. ctl+c) will cause unresolved recipes to be marked as canceled.
// Exiting early (i.e. an error occurred) will cause unresolved recipes to be marked as failed.
func (s *InstallStatus) updateFinalInstallationStatuses(installCanceled bool, isUnsupported bool) {
	s.updateRecipeStatuses(installCanceled, isUnsupported)

	log.WithFields(log.Fields{
		"hasInstalledRecipes": s.HasInstalledRecipes,
		"hasSkippedRecipes":   s.HasSkippedRecipes,
		"hasCanceledRecipes":  s.HasCanceledRecipes,
		"hasFailedRecipes":    s.HasFailedRecipes,
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
