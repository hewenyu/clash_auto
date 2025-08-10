package parser

import (
	"fmt"

	"github.com/hewenyu/clash_auto/internal/types"
	"gopkg.in/yaml.v3"
)

// Parse decodes subscription content and extracts a list of proxy and rule definitions.
func Parse(data []byte) ([]types.Proxy, []string, error) {
	var subscription struct {
		Proxies []types.Proxy `yaml:"proxies"`
		Rules   []string      `yaml:"rules"`
	}

	err := yaml.Unmarshal(data, &subscription)
	if err != nil {
		// Attempt to parse only proxies if the full structure fails
		var proxiesOnly struct {
			Proxies []types.Proxy `yaml:"proxies"`
		}
		if yaml.Unmarshal(data, &proxiesOnly) == nil {
			return proxiesOnly.Proxies, nil, nil
		}
		return nil, nil, fmt.Errorf("failed to unmarshal subscription yaml: %w", err)
	}

	return subscription.Proxies, subscription.Rules, nil
}
