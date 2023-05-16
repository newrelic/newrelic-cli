package execution

import (
	"github.com/newrelic/newrelic-cli/internal/install/segment"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type SegmentReporter struct {
	sg *segment.Segment
}

func NewSegmentReporter(sg *segment.Segment) *SegmentReporter {
	r := SegmentReporter{sg}
	return &r
}

func (r SegmentReporter) RecipeDetected(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r SegmentReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r SegmentReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r SegmentReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r SegmentReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r SegmentReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r SegmentReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	return nil
}

func (r SegmentReporter) RecipeAvailable(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r SegmentReporter) RecipeCanceled(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r SegmentReporter) InstallStarted(status *InstallStatus) error {
	return nil
}

func (r SegmentReporter) InstallComplete(status *InstallStatus) error {
	if r.sg == nil {
		return nil
	}

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

	m := make(map[string]int)
	for _, s := range status.Statuses {
		k := string("count" + s.Status)
		if v, ok := m[k]; ok {
			m[k] = v + 1
		} else {
			m[k] = 1
		}
	}

	ei := segment.NewEventInfo(types.EventTypes.InstallCompleted, "")
	for key, v := range m {
		ei.WithAdditionalInfo(key, v)
	}

	r.sg.TrackInfo(ei)
	return nil
}

func (r SegmentReporter) InstallCanceled(status *InstallStatus) error {

	return nil
}

func (r SegmentReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	if r.sg == nil {
		return nil
	}
	r.sg.Track(types.EventTypes.LicenseKeyFetchedOk)
	return nil
}

func (r SegmentReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r SegmentReporter) UpdateRequired(status *InstallStatus) error {
	return nil
}
