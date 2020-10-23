package diagnose

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

	"github.com/newrelic/newrelic-cli/internal/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var options struct {
	suites        string
	listSuites    bool
	verbose       bool
	attachmentKey string
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Troubleshoot your New Relic-instrumented application",
	Long: `Troubleshoot your New Relic-instrumented application

The diagnose command runs New Relic Diagnostics, our troubleshooting suite. The first time you run this command the nrdiag binary appropriate for your system will be downloaded to .newrelic/bin in your home directory.\n
`,
	Example: "\tnewrelic diagnose run --suites java,infra",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := ensureBinaryExists()
		if err != nil {
			log.Fatal(err)
		}
		nrdiag := exec.Command(path)
		nrdiag.Stdout = os.Stdout
		nrdiag.Stderr = os.Stderr
		nrdiag.Env = append(nrdiag.Env, "NEWRELIC_CLI_SUBPROCESS=true")
		if options.listSuites {
			nrdiag.Args = append(nrdiag.Args, "-help", "suites")
		} else if options.suites != "" {
			nrdiag.Args = append(nrdiag.Args, "-suites", options.suites)
		}
		err = nrdiag.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// TODO: flesh this out - do we want the fetch/update process to be entirely transparent, or accessible to the user?

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update the New Relic Diagnostics binary if necessary",
	Long: `Update the New Relic Diagnostics binary for your system, if it is out of date.

Checks the currently-installed version against the latest version, and if they are different, fetches the latest New Relic Diagnostics build from https://download.newrelic.com/nrdiag.`,
	Example: "newrelic diagnose update",
	Run: func(cmd *cobra.Command, args []string) {
		// FIXME: do something
	},
}

const downloadURL = "http://download.newrelic.com/nrdiag/nrdiag_latest.zip"

func ensureBinaryExists() (string, error) {
	configDirectory, err := utils.GetDefaultConfigDirectory()
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
		log.Infof("nrdiag binary not found in %s", binaryPath)
		return binaryPath, downloadBinary(binaryPath)
	}

	return binaryPath, nil
}

// TODO: break this up?
func downloadBinary(destination string) error {
	log.Info("Determining OS...")
	var executable string
	if bits.UintSize == 64 {
		executable = "nrdiag_x64"
	} else {
		executable = "nrdiag"
	}

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

	log.Infof("Downloading %s", downloadURL)
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

	log.Info("Extracting... ")
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer out.Close()

	r, err := zipped.Open()
	if err != nil {
		return err
	}
	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	Command.AddCommand(cmdRun)
	cmdRun.Flags().StringVar(&options.attachmentKey, "attachment-key", "", "Attachment key for automatic upload to a support ticket (get key from an existing ticket).")
	cmdRun.Flags().BoolVar(&options.verbose, "verbose", false, "Display verbose logging during task execution.")
	cmdRun.Flags().StringVar(&options.suites, "suites", "", "The task suite or comma-separated list of suites to run. Use --list-suites for a list of available suites.")
	cmdRun.Flags().BoolVar(&options.listSuites, "list-suites", false, "List the task suites available for the --suites argument.")
	Command.AddCommand(cmdUpdate)
}
