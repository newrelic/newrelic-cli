package install

import "golang.org/x/net/http/httpproxy"

func IsProxyConfigured() bool {
	proxyConfig := httpproxy.FromEnvironment()
	return proxyConfig.HTTPProxy != "" || proxyConfig.HTTPSProxy != "" || proxyConfig.NoProxy != ""
}

// ShouldWarnAboutProxy reports whether the proxy environment looks misconfigured:
// HTTP_PROXY is set (suggesting the user intends to use a proxy) but HTTPS_PROXY
// is absent, meaning the CLI's HTTPS calls to New Relic will bypass the proxy entirely.
func ShouldWarnAboutProxy(cfg httpproxy.Config) bool {
	return cfg.HTTPProxy != "" && cfg.HTTPSProxy == ""
}
