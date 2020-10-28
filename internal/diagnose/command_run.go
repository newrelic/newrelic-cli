package diagnose

import (
	"archive/zip"
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
		nrdiagArgs := make([]string, 0)
		if options.listSuites {
			nrdiagArgs = append(nrdiagArgs, "-help", "suites")
		} else if options.suites != "" {
			nrdiagArgs = append(nrdiagArgs, "-suites", options.suites)
		}
		err := runDiagnostics(nrdiagArgs...)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update the New Relic Diagnostics binary if necessary",
	Long: `Update the New Relic Diagnostics binary for your system, if it is out of date.

Checks the currently-installed version against the latest version, and if they are different, fetches and installs the latest New Relic Diagnostics build from https://download.newrelic.com/nrdiag.`,
	Example: "newrelic diagnose update",
	Run: func(cmd *cobra.Command, args []string) {
		err := runDiagnostics("-q", "-version")
		if err == nil {
			return
		}
		exitError, ok := err.(*exec.ExitError)
		if !ok || ok && exitError.ProcessState.ExitCode() != 1 {
			// Unexpected error
			log.Fatal(err)
		}
		err = downloadBinary()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func runDiagnostics(args ...string) error {
	err := ensureBinaryExists()
	if err != nil {
		return err
	}
	diagnostics := exec.Command(getBinaryPath())
	diagnostics.Stdout = os.Stdout
	diagnostics.Stderr = os.Stderr
	diagnostics.Env = append(diagnostics.Env, "NEWRELIC_CLI_SUBPROCESS=true")
	diagnostics.Args = append(diagnostics.Args, args...)
	return diagnostics.Run()
}

func getBinaryPath() string {
	configDirectory, err := utils.GetDefaultConfigDirectory()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(configDirectory, "bin", "nrdiag")
}

func ensureBinaryExists() error {
	destination := getBinaryPath()
	err := os.MkdirAll(path.Dir(destination), 0777)
	if err != nil {
		return err
	}

	if _, err = os.Stat(destination); os.IsNotExist(err) {
		log.Infof("nrdiag binary not found in %s", destination)
		return downloadBinary()
	}
	return nil
}

const downloadURL = "https://download.newrelic.com/nrdiag/nrdiag_latest.zip"

func downloadBinary() error {
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
		return fmt.Errorf("unknown operating system: %s", runtime.GOOS)
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

	targetPath := path.Join("nrdiag", subdir, executable)
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
	out, err := os.OpenFile(getBinaryPath(), os.O_CREATE|os.O_WRONLY, 0777)
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
