// +build unit

package apm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestApmDeployment(t *testing.T) {
	assert.Equal(t, "deployment", cmdDeployment.Name())

	testcobra.CheckCobraMetadata(t, cmdDeployment)
}

func TestApmDeploymentList(t *testing.T) {
	assert.Equal(t, "list", cmdDeploymentList.Name())

	testcobra.CheckCobraMetadata(t, cmdDeploymentList)
	testcobra.CheckCobraRequiredFlags(t, cmdDeploymentList, []string{})
}

func TestApmDeploymentCreate(t *testing.T) {
	assert.Equal(t, "create", cmdDeploymentCreate.Name())

	testcobra.CheckCobraMetadata(t, cmdDeploymentCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdDeploymentCreate,
		[]string{"revision"})

}

func TestApmDeleteDeployment(t *testing.T) {
	assert.Equal(t, "delete", cmdDeploymentDelete.Name())

	testcobra.CheckCobraMetadata(t, cmdDeploymentDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdDeploymentDelete,
		[]string{"deploymentID"})
}
