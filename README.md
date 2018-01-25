# easy-load-test-plan

## Build:
```bash
sh build.sh
```

## Run Binary (Unix)
```bash
./eztest.linux32 plan.yml
./eztest.linux64 plan.yml
./eztest.mac64 plan.yml
```

## Run Binary (Win)
```bash
eztest.win32.exe plan.yml
eztest.win64.exe plan.yml
```

## Go Run:
```bash
go run ez-test.go plan.yml
```

## Plan.yml Template:
```YAML
rate: 1          # Requests per second
duration: 5      # Duration in seconds of the test [0 = forever]
result:
  stdout: true   # Stdout
  csv: true      # Generate csv report
  plot: true     # Generate plot report

# Default Values
defaults:
  test-domain: http://test.example.com
  token: b82d30f3f1fc4e43b3f427ba3d7b9a50
  qid: 12345
  body-name: somebody

# Test Requests
requests:
  API-POST-JSON:                                        # Unique API Name
    method: POST
    url: ${test-domain}/api1/new?qid=${qid}             # URL
    headers:                                            # Headers
      Content-Type: application/json                    # JSON Body Header
      token: ${token}                                   # Replace ${key} to defaults.key's value
    body:                                               # Round-robins Bodies
      - '{"id":1, "name": "nobody"}'
      - '{"id":2, "name": "${body-name}"}'

  API-POST-FROM:
    method: POST
    url: ${test-domain}/api2/new
    headers:
      Content-Type: application/x-www-form-urlencoded   # FORM POST Body Header
      token: ${token}
    body:
      - 'id=1&name=${body-name}'

  API-GET:
    method: GET
    url: ${test-domain}/api3
    headers:
      token: ${token}
```
