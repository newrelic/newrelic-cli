package execution

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/segment"
)

type SegmentReporter struct {
	sg *segment.Segment
}

func NewSegmentReporter(sg *segment.Segment) *SegmentReporter {
	if sg == nil {
		log.Debugf("Segment reporter init with no-op")
		sg = segment.NewNoOp()
	}
	r := SegmentReporter{sg}
	return &r
}

func (r *SegmentReporter) RecipeDetected(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r *SegmentReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r *SegmentReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r *SegmentReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r *SegmentReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r *SegmentReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r *SegmentReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	return nil
}

func (r *SegmentReporter) RecipeAvailable(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r *SegmentReporter) RecipeCanceled(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r *SegmentReporter) InstallStarted(status *InstallStatus) error {
	return nil
}

func (r *SegmentReporter) InstallComplete(status *InstallStatus) error {
	if status.Error.Message != "" {
		et, ok := types.TryParseEventType(status.Error.Message)
		if ok {
			ei := segment.NewEventInfo(et, status.Error.Details)
			r.sg.TrackInfo(ei)
		} else {
			// If unclassified error, the detail would just be the error
			ei := segment.NewEventInfo(types.EventTypes.OtherError, status.Error.Message+" "+status.Error.Details)
			r.sg.TrackInfo(ei)
		}
	}

	ei := buildInstallCompleteEvent(status, types.EventTypes.InstallCompleted)
	r.sg.TrackInfo(ei)
	return nil
}

func (r *SegmentReporter) InstallCanceled(status *InstallStatus) error {
	ei := buildInstallCompleteEvent(status, types.EventTypes.InstallCancelled)
	r.sg.TrackInfo(ei)
	return nil
}

func buildInstallCompleteEvent(status *InstallStatus, et types.EventType) *segment.EventInfo {

	ei := segment.NewEventInfo(et, "")
	ei.WithAdditionalInfo("countDetected", len(status.Detected))
	ei.WithAdditionalInfo("countSkipped", len(status.Skipped))
	ei.WithAdditionalInfo("countCanceled", len(status.Canceled))
	ei.WithAdditionalInfo("countFailed", len(status.Failed))
	ei.WithAdditionalInfo("countInstalled", len(status.Installed))
	ei.WithAdditionalInfo("countUnsupported", len(status.Unsupported))
	return ei
}

func (r *SegmentReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	r.sg.Track(types.EventTypes.LicenseKeyFetchedOk)
	return nil
}

func (r *SegmentReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r *SegmentReporter) UpdateRequired(status *InstallStatus) error {
	return nil
}
