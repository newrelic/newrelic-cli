package install

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

var (
	// nolint:unused
	linuxOSReleaseFile = "/etc/os-release"
)

type discoveryManifest struct {
	processes  []genericProcess
	systemType string
	arch       string

	// nolint:unused,structcheck
	distro string
}

func newDiscoveryManifest() *discoveryManifest {
	d := discoveryManifest{
		systemType: runtime.GOOS,
		arch:       runtime.GOARCH,
	}

	return &d
}

func (d *discoveryManifest) AddProcess(p *process.Process) {
	d.processes = append(d.processes, p)
}

type genericProcess interface {
	Name() (string, error)
}

type discoverer interface {
	discover() (*discoveryManifest, error)
}

type psUtilDiscoverer struct{}

func (p *psUtilDiscoverer) discover() (*discoveryManifest, error) {
	d := newDiscoveryManifest()

	// osInfo, err := discoverOSInfo(runtime.GOOS)
	// if err != nil {
	// 	return nil, err
	// }

	// for key, val := range osInfo {
	// 	if key == "NAME" {
	// 		d.distro = value
	// 	}
	// }

	pids, err := process.PidsWithContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve processes: %s", err)
	}

	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			log.Debugf("cannot read pid %d: %s", pid, err)
			continue
		}

		d.AddProcess(p)
	}

	return d, nil
}

// nolint:deadcode,unused
func discoverOSInfo(systemType string) (map[string]string, error) {
	switch systemType {
	case "linux":
		osInfo, err := discoverLinuxOSInfo()
		if err != nil {
			return nil, err
		}

		return osInfo, nil
	default:
		return nil, fmt.Errorf("unsupported system type %s", systemType)
	}
}

// nolint:unused
func discoverLinuxOSInfo() (map[string]string, error) {
	osInfo := map[string]string{}

	file, err := os.Open(linuxOSReleaseFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err != nil && err == io.EOF {
			break
		}

		info := strings.Split(line, "=")
		if len(info) != 2 {
			continue
		}

		osInfo[info[0]] = strings.Trim(info[1], "\"\n")
	}

	return osInfo, nil
}
