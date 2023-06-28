package execution

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestUpdateRhapsody(t *testing.T) {
	s1 := NewInstallStatus(types.InstallerContext{}, nil, nil)
	s1.Statuses = append(s1.Statuses, &RecipeStatus{Name: "infrastructure-agent-installer", Status: RecipeStatusTypes.INSTALLED}, &RecipeStatus{Name: "logs-integration", Status: RecipeStatusTypes.INSTALLED})
	tests := []struct {
		name    string
		s       *InstallStatus
		wantErr bool
	}{
		{"do the thing", s1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateRhapsody(tt.s); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRhapsody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
