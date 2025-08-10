package generator

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// GenerateConfig reads a template, adds proxies and rules, and writes the final config.
func GenerateConfig(templatePath, outputPath string, proxies []map[string]interface{}, additionalRules []string) error {
	// 1. Read template file
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// 2. Unmarshal template into a map
	var config map[string]interface{}
	if err := yaml.Unmarshal(templateData, &config); err != nil {
		return fmt.Errorf("failed to unmarshal template yaml: %w", err)
	}

	// 3. Add proxies to the config map
	config["proxies"] = proxies

	// 4. Extract proxy names
	var proxyNames []string
	for _, p := range proxies {
		if name, ok := p["name"].(string); ok {
			proxyNames = append(proxyNames, name)
		}
	}

	// 5. Add proxy names to proxy groups
	updateProxyGroups(config, proxyNames)

	// 6. Merge and deduplicate rules
	config["rules"] = mergeRules(config["rules"], additionalRules)

	// 7. Marshal the updated map back to YAML
	finalYAML, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal final yaml: %w", err)
	}

	// 8. Write to output file
	if err := os.WriteFile(outputPath, finalYAML, 0644); err != nil {
		return fmt.Errorf("failed to write final config file: %w", err)
	}

	return nil
}

// updateProxyGroups finds a specific proxy group and populates it with the given proxy names.
func updateProxyGroups(config map[string]interface{}, proxyNames []string) {
	proxyGroups, ok := config["proxy-groups"].([]interface{})
	if !ok {
		return // No proxy groups found
	}

	// The name of the group to populate with all filtered proxies.
	// This could be made configurable in config.yaml in the future.
	const targetGroupName = "线路选择"

	for i, group := range proxyGroups {
		groupMap, ok := group.(map[string]interface{})
		if !ok {
			continue
		}

		groupName, ok := groupMap["name"].(string)
		if !ok {
			continue
		}

		if groupName == targetGroupName {
			// Convert []string to []interface{} for YAML marshaling
			namesAsInterface := make([]interface{}, len(proxyNames))
			for j, name := range proxyNames {
				namesAsInterface[j] = name
			}

			// Set the proxies for the target group
			groupMap["proxies"] = namesAsInterface
			proxyGroups[i] = groupMap

			// We found and updated the target group, no need to check others.
			break
		}
	}
}

// mergeRules combines rules from the template and subscriptions, removing duplicates.
func mergeRules(templateRules interface{}, additionalRules []string) []string {
	seen := make(map[string]bool)
	var finalRules []string

	// Add rules from template first
	if tRules, ok := templateRules.([]interface{}); ok {
		for _, r := range tRules {
			if ruleStr, ok := r.(string); ok {
				if _, exists := seen[ruleStr]; !exists {
					seen[ruleStr] = true
					finalRules = append(finalRules, ruleStr)
				}
			}
		}
	}

	// Add rules from subscriptions
	for _, ruleStr := range additionalRules {
		if _, exists := seen[ruleStr]; !exists {
			seen[ruleStr] = true
			finalRules = append(finalRules, ruleStr)
		}
	}

	return finalRules
}
