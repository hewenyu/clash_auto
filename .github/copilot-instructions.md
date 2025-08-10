# Clash Auto - VPN Subscription Merger

Clash Auto is a Go application that automatically fetches multiple Clash VPN subscriptions, filters proxy nodes by keywords, and generates a merged `config.yaml` file using a template system.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Bootstrap and Build
- Ensure Go 1.24.2+ is installed: `go version`
- Download dependencies: `go mod tidy` -- takes 3-5 seconds
- Build the application: `go build -o clash-auto cmd/clash-auto/main.go` -- takes 15 seconds. NEVER CANCEL. Set timeout to 30+ seconds.
- Or run directly: `go run cmd/clash-auto/main.go` -- takes 5 seconds

### Configuration Setup
- Copy example config: `cp config/config.yaml.example config/config.yaml`
- Edit `config/config.yaml` to add your subscription URLs
- Ensure `dist/` directory exists: `mkdir -p dist`
- Template file is at `config/template.yaml` - modify as needed

### Running the Application
- With default config: `./clash-auto`
- With custom config: `./clash-auto -c path/to/config.yaml`
- Application execution is very fast (under 1 second) unless network downloads are slow
- **CRITICAL**: Application requires network connectivity to download subscriptions

### Development Tools
- Format code: `go fmt ./...` -- instant
- Check for errors: `go vet ./...` -- takes 2 seconds
- Organize imports: `~/go/bin/goimports -w .` (install first with `go install golang.org/x/tools/cmd/goimports@latest`)
- Run tests: `go test ./...` -- takes 2 seconds. NEVER CANCEL. Set timeout to 10+ seconds.

## Validation

### Manual Testing Scenarios
**ALWAYS test these scenarios after making changes:**

1. **Basic functionality test:**
   ```bash
   # Create a test server (in separate terminal/session)
   cat > /tmp/test_server.go << 'EOF'
   package main
   import ("fmt"; "net/http"; "net/http/httptest"; "time")
   func main() {
       testYAML := `proxies:
     - name: "HK-Test-Node"
       type: ss
       server: 1.2.3.4
       port: 8080
       cipher: aes-128-gcm
       password: test123
     - name: "SG-Singapore-Test"
       type: ss
       server: 5.6.7.8
       port: 8080
       cipher: aes-128-gcm
       password: test456
   rules:
     - "DOMAIN-SUFFIX,example.com,PROXY"`
       server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           w.WriteHeader(http.StatusOK); w.Write([]byte(testYAML))
       }))
       fmt.Printf("Test server: %s\n", server.URL)
       time.Sleep(60 * time.Second)
   }
   EOF
   go run /tmp/test_server.go &
   SERVER_PID=$!
   
   # Update SERVER_URL with the printed URL from above
   SERVER_URL="http://127.0.0.1:XXXXX"  # Replace XXXXX with actual port
   
   # Create test config
   cat > /tmp/test_config.yaml << EOF
   subscriptions:
     - "$SERVER_URL"
   filter_rules:
     include_keywords: ["HK", "SG"]
   template_path: "./config/template.yaml"
   output_path: "./dist/test_config.yaml"
   EOF
   
   # Test the application
   ./clash-auto -c /tmp/test_config.yaml
   
   # Verify output
   cat ./dist/test_config.yaml
   
   # Cleanup
   kill $SERVER_PID
   ```

2. **Error handling test:**
   ```bash
   # Test missing config
   ./clash-auto -c nonexistent.yaml  # Should fail gracefully
   
   # Test invalid URLs in config
   echo "subscriptions: ['http://invalid-url']" > /tmp/bad_config.yaml
   ./clash-auto -c /tmp/bad_config.yaml  # Should skip invalid URLs
   ```

3. **Build validation:**
   ```bash
   # Clean build test
   go clean
   go build -o clash-auto cmd/clash-auto/main.go
   ./clash-auto --help  # Should show usage
   ```

### Expected Behavior
- **Success case**: Application downloads, filters, and generates config in under 5 seconds
- **Network failure**: Application continues and skips failed subscriptions
- **No matches**: Application fails if no proxies match filter keywords
- **Generated config**: Should contain filtered proxies in the "线路选择" proxy group

## Common Tasks

### Repository Structure
```
/home/runner/work/clash_auto/clash_auto/
├── cmd/clash-auto/main.go          # Main application entry point
├── internal/                       # Internal modules
│   ├── config/config.go           # Configuration loading
│   ├── downloader/downloader.go   # HTTP downloading
│   ├── filter/filter.go           # Proxy filtering logic
│   ├── generator/generator.go     # Config file generation
│   ├── parser/parser.go           # YAML parsing
│   └── types/types.go             # Data structures
├── config/
│   ├── config.yaml.example        # Example configuration
│   └── template.yaml              # Clash config template
├── dist/                          # Generated output files
├── go.mod                         # Go module definition
└── README.md                      # Chinese documentation
```

### Key Files Content
**config/config.yaml format:**
```yaml
subscriptions:
  - "https://your-subscription-url.com/sub.yaml"
filter_rules:
  include_keywords: ["HK", "SG", "US"]
template_path: "./config/template.yaml"
output_path: "./dist/config.yaml"
```

**Expected subscription format:**
```yaml
proxies:
  - name: "HK-Node-01"
    type: ss
    server: ip.address
    port: 8080
    cipher: aes-128-gcm
    password: password
rules:
  - "DOMAIN-SUFFIX,example.com,PROXY"
```

## Troubleshooting

### Common Issues
- **"No such file or directory"**: Copy `config.yaml.example` to `config.yaml`
- **"No proxies were successfully parsed"**: Check subscription URLs and network connectivity
- **"No proxies left after filtering"**: Verify filter keywords match proxy names
- **Build errors**: Run `go mod tidy` to fix dependency issues

### Network Dependencies
- Application requires internet access to download subscriptions
- Uses standard HTTP GET requests
- No authentication required for subscription URLs
- Supports HTTP and HTTPS protocols

### Performance Notes
- Build time: ~15 seconds (small codebase, minimal dependencies)
- Runtime: <1 second for local processing, variable for network downloads
- Memory usage: Minimal (processes YAML in memory)
- No persistent storage required

## Development Guidelines

### Code Organization
- **cmd/**: Application entry points
- **internal/**: Private application code (not importable by other projects)
- **config/**: User configuration and templates
- **dist/**: Generated output (excluded from git)

### Dependencies
- **Only dependency**: `gopkg.in/yaml.v3` for YAML processing
- **Standard library**: HTTP client, file I/O, string processing
- **No external tools required** beyond Go toolchain

### Making Changes
1. Always run `go fmt ./...` before committing
2. Verify with `go vet ./...` to catch common errors
3. Test with the validation scenarios above
4. Build and run the application to ensure it works
5. Check that generated config files are valid YAML

**NEVER** modify working proxy filtering or YAML generation logic without thorough testing - these are core functions.