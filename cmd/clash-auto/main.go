package main

import (
	"flag"
	"log"

	"github.com/hewenyu/clash_auto/internal/config"
	"github.com/hewenyu/clash_auto/internal/downloader"
	"github.com/hewenyu/clash_auto/internal/filter"
	"github.com/hewenyu/clash_auto/internal/generator"
	"github.com/hewenyu/clash_auto/internal/parser"
	"github.com/hewenyu/clash_auto/internal/types"
)

func main() {
	configPath := flag.String("c", "./config/config.yaml", "Path to the configuration file")
	flag.Parse()

	// 1. Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Println("Configuration loaded successfully.")

	// 2. Download and parse all subscriptions
	var allProxies []types.Proxy
	var allRules []string
	for _, subURL := range cfg.Subscriptions {
		log.Printf("Downloading subscription from: %s", subURL)
		subData, err := downloader.Download(subURL)
		if err != nil {
			log.Printf("Failed to download from %s: %v. Skipping.", subURL, err)
			continue
		}

		proxies, rules, err := parser.Parse(subData)
		if err != nil {
			log.Printf("Failed to parse subscription from %s: %v. Skipping.", subURL, err)
			continue
		}
		if proxies != nil {
			allProxies = append(allProxies, proxies...)
			log.Printf("Successfully parsed %d proxies from %s.", len(proxies), subURL)
		}
		if rules != nil {
			allRules = append(allRules, rules...)
			log.Printf("Successfully parsed %d rules from %s.", len(rules), subURL)
		}
	}

	if len(allProxies) == 0 {
		log.Fatalf("No proxies were successfully parsed from any subscription. Aborting.")
	}
	log.Printf("Total proxies collected: %d", len(allProxies))

	// 3. Filter proxies
	filteredProxies := filter.FilterProxies(allProxies, cfg.FilterRules.IncludeKeywords)
	log.Printf("Filtered proxies: %d remaining.", len(filteredProxies))

	if len(filteredProxies) == 0 {
		log.Fatalf("No proxies left after filtering. Aborting.")
	}

	// 4. Generate final config
	err = generator.GenerateConfig(cfg.TemplatePath, cfg.OutputPath, filteredProxies, allRules)
	if err != nil {
		log.Fatalf("Error generating final config: %v", err)
	}

	log.Printf("Successfully generated config file at: %s", cfg.OutputPath)
}
