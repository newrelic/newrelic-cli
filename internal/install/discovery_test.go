// +build unit

package install

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverOSInfo_Linux(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "TestCollectLinuxOSInfo")
	if err != nil {
		t.Fatalf("could not create mock OS release file")
	}

	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write([]byte(mockOSReleaseFile)); err != nil {
		t.Fatalf("could not write to mock OS release file")
	}

	linuxOSReleaseFile = tmpFile.Name()

	osInfo, err := discoverOSInfo("linux")

	require.NoError(t, err)
	assert.NotNil(t, osInfo)
	assert.Equal(t, "Ubuntu", osInfo["NAME"])

	if err = tmpFile.Close(); err != nil {
		t.Fatalf("could not close mock OS release file")
	}
}

func TestDiscoverOSInfo_Default(t *testing.T) {
	_, err := discoverOSInfo("not recognized")

	assert.Errorf(t, err, "unsupported system type")
}

const (
	mockOSReleaseFile = `
NAME="Ubuntu"
VERSION="20.04.1 LTS (Focal Fossa)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 20.04.1 LTS"
VERSION_ID="20.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=focal
UBUNTU_CODENAME=focal`
)
