package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/hewenyu/clash_auto/internal/types"
	"gopkg.in/yaml.v3"
)

// GenerateConfig reads a template, adds proxies and rules, and writes the final config.
func GenerateConfig(templatePath, outputPath string, proxies []types.Proxy, additionalRules []string) error {
	// 1. Read template file
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// 2. Unmarshal template into our strong-typed struct
	var finalConfig types.Config
	if err := yaml.Unmarshal(templateData, &finalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal template yaml: %w", err)
	}

	// 3. Add the filtered proxies from subscriptions
	finalConfig.Proxies = proxies

	// 4. Extract proxy names
	var proxyNames []string
	for _, p := range proxies {
		if name, ok := p["name"].(string); ok {
			proxyNames = append(proxyNames, name)
		}
	}

	// 5. Add proxy names to the target proxy group
	updateProxyGroups(&finalConfig, proxyNames)

	// 6. Merge and deduplicate rules
	finalConfig.Rules = mergeRules(finalConfig.Rules, additionalRules)

	// 7. Marshal the updated struct back to YAML
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2) // for pretty printing
	if err := encoder.Encode(&finalConfig); err != nil {
		return fmt.Errorf("failed to marshal final config: %w", err)
	}

	// 8. Ensure the output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 9. Post-process the output to un-escape unicode characters
	// The yaml encoder escapes emojis and other Unicode characters, so we replace them manually.
	result := replaceUnicodeEscapes(buf.String())

	// 10. Write to output file
	if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
		return fmt.Errorf("failed to write final config file: %w", err)
	}

	return nil
}

// updateProxyGroups finds a specific proxy group and populates it with the given proxy names.
func updateProxyGroups(config *types.Config, proxyNames []string) {
	const targetGroupName = "线路选择"

	for i, group := range config.ProxyGroups {
		if group.Name == targetGroupName {
			// Set the proxies for the target group
			config.ProxyGroups[i].Proxies = proxyNames
			// We found and updated the target group, no need to check others.
			return
		}
	}
}

// mergeRules combines rules from the template and subscriptions, removing duplicates.
func mergeRules(templateRules []string, additionalRules []string) []string {
	seen := make(map[string]bool)
	var finalRules []string

	// Add rules from template first
	for _, ruleStr := range templateRules {
		if _, exists := seen[ruleStr]; !exists {
			seen[ruleStr] = true
			finalRules = append(finalRules, ruleStr)
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

// replaceUnicodeEscapes replaces Unicode escape sequences like \U0001F1E8 and \u0065 with actual Unicode characters.
func replaceUnicodeEscapes(s string) string {
	// Find all Unicode escape sequences like \U0001F1E8 (8 hex digits)
	re8 := regexp.MustCompile(`\\U([0-9A-Fa-f]{8})`)
	s = re8.ReplaceAllStringFunc(s, func(match string) string {
		hex := match[2:] // Remove \U prefix
		if codePoint, err := strconv.ParseInt(hex, 16, 64); err == nil {
			return string(rune(codePoint))
		}
		return match
	})

	// Find all Unicode escape sequences like \u0065 (4 hex digits)
	re4 := regexp.MustCompile(`\\u([0-9A-Fa-f]{4})`)
	s = re4.ReplaceAllStringFunc(s, func(match string) string {
		hex := match[2:] // Remove \u prefix
		if codePoint, err := strconv.ParseInt(hex, 16, 64); err == nil {
			return string(rune(codePoint))
		}
		return match
	})

	return s
}
