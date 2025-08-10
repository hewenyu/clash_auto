package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type TestProxy map[string]interface{}

func replaceUnicodeEscapes(s string) string {
	// Find all Unicode escape sequences like \U0001F1E8
	re := regexp.MustCompile(`\\U([0-9A-Fa-f]{8})`)
	
	return re.ReplaceAllStringFunc(s, func(match string) string {
		// Extract the hex digits
		hex := match[2:] // Remove \U prefix
		
		// Convert hex to int64
		if codePoint, err := strconv.ParseInt(hex, 16, 64); err == nil {
			// Convert to rune and then to string
			return string(rune(codePoint))
		}
		
		// If conversion fails, return original
		return match
	})
}

func main() {
	// Create a test proxy with Unicode characters (Chinese flag emoji)
	proxy := TestProxy{
		"name": "🇨🇳 Hong Kong",
		"type": "ss", 
		"server": "example.com",
		"port": 443,
	}

	// Test current approach (what's in the code)
	fmt.Println("=== Current approach (with manual post-processing) ===")
	var buf1 bytes.Buffer
	encoder1 := yaml.NewEncoder(&buf1)
	encoder1.SetIndent(2)
	encoder1.Encode(proxy)
	
	// Manual post-processing (current code)
	quoted := `"` + string(buf1.Bytes()) + `"`
	unquoted, err := strconv.Unquote(quoted)
	if err != nil {
		fmt.Printf("Error with unquote: %v\n", err)
		unquoted = string(buf1.Bytes())
	}
	fmt.Printf("Result: %s\n", unquoted)
	
	// Check if Chinese flag emoji is properly rendered
	if strings.Contains(unquoted, "🇨🇳") {
		fmt.Println("✅ Chinese flag emoji rendered correctly")
	} else {
		fmt.Println("❌ Chinese flag emoji NOT rendered correctly")
		if strings.Contains(unquoted, "\\U") {
			fmt.Println("   Found Unicode escape sequences")
		}
	}
	
	// Test alternative approaches
	fmt.Println("\n=== Alternative approach (no post-processing) ===")
	var buf2 bytes.Buffer
	encoder2 := yaml.NewEncoder(&buf2)
	encoder2.SetIndent(2)
	encoder2.Encode(proxy)
	result2 := string(buf2.Bytes())
	fmt.Printf("Result: %s\n", result2)
	
	if strings.Contains(result2, "🇨🇳") {
		fmt.Println("✅ Chinese flag emoji rendered correctly")
	} else {
		fmt.Println("❌ Chinese flag emoji NOT rendered correctly")
		if strings.Contains(result2, "\\U") {
			fmt.Println("   Found Unicode escape sequences")
		}
	}

	// Test manual replacement approach
	fmt.Println("\n=== Manual regex replacement approach ===")
	var buf3 bytes.Buffer
	encoder3 := yaml.NewEncoder(&buf3)
	encoder3.SetIndent(2)
	encoder3.Encode(proxy)
	result3 := string(buf3.Bytes())
	
	// Replace Unicode escape sequences with actual characters
	result3 = replaceUnicodeEscapes(result3)
	fmt.Printf("Result: %s\n", result3)
	
	if strings.Contains(result3, "🇨🇳") {
		fmt.Println("✅ Chinese flag emoji rendered correctly")
	} else {
		fmt.Println("❌ Chinese flag emoji NOT rendered correctly")
	}
}