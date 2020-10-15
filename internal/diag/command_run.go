package diag

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/bits"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdDiag = &cobra.Command{
	Use:   "run",
	Short: "Troubleshoot your New Relic-instrumented application",
	Long: `Troubleshoot your New Relic-instrumented application

The diag command runs New Relic Diagnostic, our troubleshooting suite. The first time you run this command the nrdiag binary appropriate for your system will be downloaded to .newrelic/bin in your home directory.
`,
	Example: "newrelic diag",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := ensureBinaryExists()
		if err != nil {
			log.Fatal(err)
		}
		nrdiag := exec.Command(path)
		nrdiag.Stdout = os.Stdout
		nrdiag.Stderr = os.Stderr
		nrdiag.Args = args
		err = nrdiag.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// TODO: flesh this out - do we want the fetch/update process to be entirely transparent, or accessible to the user?

/*var cmdFetch = &cobra.Command{
	Use: "fetch",
	Short: "Download the New Relic Diagnostics binary",
	Long: `Download the New Relic Diagnostics binary for your system.

Fetches the latest New Relic Diagnostics build from https://download.newrelic.com/nrdiag. If a binary is already installed, it will be updated if necessary.`,
	Example: "newrelic diag fetch",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
*/

const downloadURL = "http://download.newrelic.com/nrdiag/nrdiag_latest.zip"

// FIXME: this should be somewhere globally-available (copied from internal/config/config.go)
func getDefaultConfigDirectory() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.newrelic", home), nil
}

func ensureBinaryExists() (string, error) {
	configDirectory, err := getDefaultConfigDirectory()
	if err != nil {
		return "", err
	}

	binPath := path.Join(configDirectory, "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		err = os.Mkdir(binPath, 0777)
		if err != nil {
			return "", err
		}
	}

	binaryPath := path.Join(binPath, "nrdiag")
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return binaryPath, downloadBinary(binaryPath)
	}

	return binaryPath, nil
}

// TODO: break this up?
func downloadBinary(destination string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "nrdiag-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return err
	}
	tmpFile.Close()

	zipReader, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		return err
	}
	defer zipReader.Close()
	var executable string
	// FIXME: is this the best way to detect this?
	if bits.UintSize == 64 {
		executable = "nrdiag_x64"
	} else {
		executable = "nrdiag"
	}
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, 0777)
	var subdir string
	if runtime.GOOS == "windows" {
		subdir = "win"
		executable = executable + ".exe"
	} else if runtime.GOOS == "darwin" {
		subdir = "mac"
	} else if runtime.GOOS == "linux" {
		subdir = "linux"
	} else {
		return errors.New("Unknown operating system: " + runtime.GOOS)
	}
	if err != nil {
		return err
	}

	targetPath := fmt.Sprintf("nrdiag/%s/%s", subdir, executable)
	var zipped *zip.File
	for _, f := range zipReader.File {
		if f.Name == targetPath {
			zipped = f
			break
		}
	}
	if zipped == nil {
		return fmt.Errorf("executable %s not found in zip file", targetPath)
	}
	r, err := zipped.Open()
	if err != nil {
		return err
	}
	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	return os.Chmod(destination, 0755)
}

func init() {
	Command.AddCommand(cmdDiag)
}
