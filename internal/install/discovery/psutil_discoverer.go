package discovery

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type PSUtilDiscoverer struct {
	processFilterer recipes.ProcessFilterer
}

func NewPSUtilDiscoverer(f recipes.ProcessFilterer) *PSUtilDiscoverer {
	d := PSUtilDiscoverer{
		processFilterer: f,
	}

	return &d
}

func (p *PSUtilDiscoverer) Discover(ctx context.Context) (*types.DiscoveryManifest, error) {
	i, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	m := types.DiscoveryManifest{
		Hostname:        i.Hostname,
		KernelArch:      i.KernelArch,
		KernelVersion:   i.KernelVersion,
		OS:              i.OS,
		Platform:        i.Platform,
		PlatformFamily:  i.PlatformFamily,
		PlatformVersion: i.PlatformVersion,
	}

	m = filterValues(m)

	pids, err := process.PidsWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve processes: %s", err)
	}

	processes := []types.GenericProcess{}
	for _, pid := range pids {
		var pp *process.Process
		pp, err = process.NewProcess(pid)
		if err != nil {
			if err != process.ErrorProcessNotRunning {
				log.Debugf("cannot read pid %d: %s", pid, err)
			}
			continue
		}

		processes = append(processes, PSUtilProcess(*pp))
	}

	m.DiscoveredProcesses = processes

	return &m, nil
}

func filterValues(m types.DiscoveryManifest) types.DiscoveryManifest {
	if !isValidOpenInstallationPlatform(m.Platform) {
		m.Platform = ""
	}

	if !isValidOpenInstallationPlatformFamily(m.PlatformFamily) {
		m.PlatformFamily = ""
	}

	return m
}

func isValidOpenInstallationPlatform(platform string) bool {
	s := reflect.ValueOf(&types.OpenInstallationPlatformTypes).Elem()

	for i := 0; i < s.NumField(); i++ {
		v := s.Field(i).Interface().(types.OpenInstallationPlatform)
		if strings.EqualFold(string(v), platform) {
			return true
		}
	}

	return false
}

func isValidOpenInstallationPlatformFamily(platformFamily string) bool {
	s := reflect.ValueOf(&types.OpenInstallationPlatformFamilyTypes).Elem()

	for i := 0; i < s.NumField(); i++ {
		v := s.Field(i).Interface().(types.OpenInstallationPlatformFamily)
		if strings.EqualFold(string(v), platformFamily) {
			return true
		}
	}

	return false
}
