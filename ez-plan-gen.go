package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	
	"./ez"
)

var csvs = make(map[string][]map[string]string)
var csvRegex = regexp.MustCompile(`\${(?P<csv>[a-zA-Z1-9]+\.csv)\.(?P<column>[a-zA-Z_]+)}`)
var dynamicIdRegex = regexp.MustCompile(`\${(?P<prefix>[a-zA-Z0-9]+)\[(?P<from>\d+):(?P<to>\d+)\]}`)

func main() {
	if len(os.Args) == 1 {
		panic("No yaml file")
	}

	fileName := os.Args[1]
	ap := ez.NewTestPlan(fileName)
	file, err := os.Create(strings.Replace(fileName, ".yml", fmt.Sprintf("-%vX%v.yml", ap.Rate, ap.Duration.Seconds()), -1))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	WriteStaticFields(file, ap)
	WriteRequests(file, ap)

	println(ap)
}

func WriteStaticFields(file *os.File, ap *ez.TestPlan) {
	file.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	file.WriteString(fmt.Sprintf("rate: %v\n", ap.Rate))
	file.WriteString(fmt.Sprintf("duration: %v\n", ap.Duration))

	file.WriteString("result:\n")
	file.WriteString(fmt.Sprintf("  stdout: %v\n", ap.Result.Stdout))
	file.WriteString(fmt.Sprintf("  csv: %v\n", ap.Result.CSV))
	file.WriteString(fmt.Sprintf("  plot: %v\n\n", ap.Result.Plot))

	if len(ap.Defaults) == 0 {
		return
	}

	file.WriteString("defaults:\n")
	for k, v := range ap.Defaults {
		file.WriteString(fmt.Sprintf("  %v: %v\n", k, v))
	}
	file.WriteString("\n")
}

func WriteRequests(file *os.File, ap *ez.TestPlan) {
	if len(ap.Requests) == 0 {
		return
	}

	file.WriteString("requests:\n")
	for n, req := range ap.Requests {
		csvMatches := csvRegex.FindAllStringSubmatch(req.URL, -1)
		dynamicIdMatches := dynamicIdRegex.FindAllStringSubmatch(req.URL, -1)

		if len(csvMatches) > 0 {
			var rowCount int
			for _, m := range csvMatches {
				fileName := m[1]
				csv, ok := csvs[fileName]
				if !ok {
					csv = ez.CSVReader(fileName)
					csvs[fileName] = csv
				}

				rowCount = len(csv)
			}

			originalURL := req.URL
			for i := 0; i < rowCount; i++ {
				formattedURL := originalURL

				for _, m := range csvMatches {
					formattedURL = strings.Replace(formattedURL, m[0], csvs[m[1]][i][m[2]], -1)
				}

				req.URL = formattedURL
				WriteRequest(file, n + "-" + strconv.Itoa(i), req)
			}
		} else if len(dynamicIdMatches) > 0 {
			origin :=  dynamicIdMatches[0][0]
			prefix := dynamicIdMatches[0][1]
			from, _ := strconv.Atoi(dynamicIdMatches[0][2])
			to, _ := strconv.Atoi(dynamicIdMatches[0][3])

			originalURL := req.URL
			for i := from; i <= to ; i++ {
				formattedURL := originalURL
				formattedURL = strings.Replace(formattedURL, origin, prefix + strconv.Itoa(i), -1)

				req.URL = formattedURL
				WriteRequest(file, n + "-" + strconv.Itoa(i), req)
			}

		} else {
			WriteRequest(file, n, req)
		}
	}
}

func WriteRequest(file *os.File, requestName string, r ez.Request) {
	file.WriteString(fmt.Sprintf("  %v: \n", requestName))
	file.WriteString(fmt.Sprintf("    method: %v\n", r.Method))
	file.WriteString(fmt.Sprintf("    url: %v\n", r.URL))

	if len(r.Headers) > 0 {
		file.WriteString("    headers:\n")

		for hk, hv := range r.Headers {
			file.WriteString(fmt.Sprintf("      %v: %v\n", hk, hv))
		}
	}

	if len(r.Body) > 0 {
		file.WriteString("    body:\n")

		for _, b := range r.Body {
			csvMatches := csvRegex.FindAllStringSubmatch(b, -1)
			dynamicIdMatches := dynamicIdRegex.FindAllStringSubmatch(b, -1)

			if len(csvMatches) > 0 {
				var rowCount int
				for _, m := range csvMatches {
					fileName := m[1]
					csv, ok := csvs[fileName]
					if !ok {
						csv = ez.CSVReader(fileName)
						csvs[fileName] = csv
					}

					rowCount = len(csv)
				}

				originalBody := b
				for i := 0; i < rowCount; i++ {
					formattedBody := originalBody

					for _, m := range csvMatches {
						formattedBody = strings.Replace(formattedBody, m[0], csvs[m[1]][i][m[2]], -1)
					}

					file.WriteString(fmt.Sprintf("      - '%v'\n", formattedBody))
				}
			} else if len(dynamicIdMatches) > 0 {
				origin :=  dynamicIdMatches[0][0]
				prefix := dynamicIdMatches[0][1]
				from, _ := strconv.Atoi(dynamicIdMatches[0][2])
				to, _ := strconv.Atoi(dynamicIdMatches[0][3])

				originalBody := b
				for i := from; i <= to ; i++ {
					formattedBody := strings.Replace(originalBody, origin, prefix + strconv.Itoa(i), -1)
					file.WriteString(fmt.Sprintf("      - '%v'\n", formattedBody))
				}
			} else {
				file.WriteString(fmt.Sprintf("      - '%v'\n", b))
			}
		}
	}

	file.WriteString("\n")
}