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

## How to generate Plan.yml using csv data:
```bash
go run ez-plan-gen.go abstract-plan.yml
```

## Abstract-plan.yml Template:
```YAML
rate: 1          # Requests per second
duration: 5s     # Duration in seconds of the test [0 = forever]

# Default Values
defaults:
  test-domain: http://test.example.com
  token: b82d30f3f1fc4e43b3f427ba3d7b9a50
  qid: 12345
  body-name: somebody

# Test Requests
requests:
  API-POST-JSON-FROM-CSV:
    method: POST
    url: ${test-domain}/api1/new?qid=${qid}
    headers:
      Content-Type: application/json
      token: ${token}
    body: # Load id,name from "from.csv" into body iteratively
      - '{"id":${from.csv.id}, "name": "${from.csv.name}"}'   
  API-POST-JSON-DYNAMIC-ID:
    method: POST
    url: ${test-domain}/api1/new?qid=${qid}
    headers:
      Content-Type: application/json
      token: ${token}
    body: # Generate body with id from "test1" to id "test100"
      - '{"id":${test[1:100]}, "name": "test"}'
  API-GET-FROM-CSV: # Generate requests iteratively from "from.csv"
    method: GET
    url: ${test-domain}/api1/${from.csv.id}?name=${from.csv.name}
    headers:
      token: ${token}
  API-GET-DYNAMIC-ID: # Generate requests with url from "/api1/test1" to "/api/test100" 
    method: GET
    url: ${test-domain}/api1/${test[1:100]}
    headers:
      token: ${token}      
```