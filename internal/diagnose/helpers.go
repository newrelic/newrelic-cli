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

	"github.com/newrelic/newrelic-cli/internal/config"

	log "github.com/sirupsen/logrus"
)

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
		log.Warnf("failed to download the latest nrdiag: %s", err)
		log.Infof("If this problem persists, you can download the zip file from https://download.newrelic.com/nrdiag/nrdiag_latest.zip and place the appropriate binary for your system in %s/bin, then try again.", config.ConfigDir)
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
	return path.Join(config.ConfigDir, "bin", "nrdiag")
}

const downloadURL = "https://download.newrelic.com/nrdiag/nrdiag_latest.zip"
