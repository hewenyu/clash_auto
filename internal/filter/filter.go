package filter

import (
	"strings"

	"github.com/hewenyu/clash_auto/internal/types"
)

// FilterProxies filters a list of proxy definitions based on include keywords.
// It checks if the proxy's name contains any of the keywords.
func FilterProxies(proxies []types.Proxy, keywords []string) []types.Proxy {
	var filteredProxies []types.Proxy

	if len(keywords) == 0 {
		return proxies
	}

	for _, proxy := range proxies {
		name, ok := proxy["name"].(string)
		if !ok {
			continue // Skip proxies without a valid name
		}

		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(name), strings.ToLower(keyword)) {
				filteredProxies = append(filteredProxies, proxy)
				break // Move to the next proxy once a match is found
			}
		}
	}

	return filteredProxies
}
