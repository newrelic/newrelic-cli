package discovery

import (
	"context"
	"reflect"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Discoverer

type Discoverer interface {
	Discover(context.Context) (*types.DiscoveryManifest, error)
}

type PSUtilDiscoverer struct{}

func NewPSUtilDiscoverer() *PSUtilDiscoverer {
	return &PSUtilDiscoverer{}
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

	log.Debugf("discovered manifest %+v", m)

	m = filterValues(m)

	log.Debugf("filtered manifest %+v", m)

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
