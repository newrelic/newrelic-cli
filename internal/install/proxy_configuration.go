package install

import "golang.org/x/net/http/httpproxy"

func IsProxyConfigured() bool {
	proxyConfig := httpproxy.FromEnvironment()
	return proxyConfig.HTTPProxy != "" || proxyConfig.HTTPSProxy != "" || proxyConfig.NoProxy != ""
}
