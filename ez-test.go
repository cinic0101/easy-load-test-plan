package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/http"
	"regexp"
	"os"
	"strings"

	"./ez"
	"github.com/tsenart/vegeta/lib"
)

func main() {
	if len(os.Args) == 1 {
		panic("No yaml file")
	}

	fileName := os.Args[1]
	p := ez.NewTestPlan(fileName)

	folderName := strings.Replace(os.Args[1], ".yml", "", 1)
	os.MkdirAll("./" + folderName, os.ModePerm)

	var csvRows [][]string
	for requestName, r := range p.Requests {
		rate := p.Rate
		duration := p.Duration

		header := http.Header{}
		for hk := range r.Headers {
			h := r.Headers[hk]
			ReplaceFromDefaults(&h, p.Defaults)
			header.Add(hk, h)
		}

		ReplaceFromDefaults(&r.URL, p.Defaults)

		var targets []vegeta.Target
		if len(r.Body) > 0 {
			for _, b := range r.Body {
				ReplaceFromDefaults(&b, p.Defaults)
				targets = append(targets, vegeta.Target{
					Method: r.Method,
					URL:    r.URL,
					Header: header,
					Body:   []byte(b),
				})
			}
		} else {
			targets = append(targets, vegeta.Target{
				Method: r.Method,
				URL:    r.URL,
				Header: header,
			})
		}

		target := vegeta.NewStaticTargeter(targets...)
		attacker := vegeta.NewAttacker()

		if p.Timeout != nil {
			vegeta.Timeout(*p.Timeout)(attacker)
		}

		if p.Workers != nil {
			vegeta.Workers(*p.Workers)(attacker)
		}

		var metrics vegeta.Metrics
		var results vegeta.Results
		for res := range attacker.Attack(target, rate, duration) {
			metrics.Add(res)
			if p.Result.Plot {
				results.Add(res)
			}
		}
		metrics.Close()

		if p.Result.Stdout {
			Stdout(requestName, metrics)
		}

		if p.Result.CSV {
			csvRows = append(csvRows, []string{
				requestName,
				r.URL,
				fmt.Sprintf("%.1f", metrics.Success * 100),
				fmt.Sprintf("%.3f", metrics.Latencies.Mean.Seconds()),
				fmt.Sprintf("%.3f", metrics.Latencies.Max.Seconds()),
			})
		}

		if p.Result.Plot {
			DrawPlot(folderName, requestName, results)
		}
	}

	if p.Result.CSV {
		header := [][]string{{ "Name", "URL", "Success(%)", "Mean(s)", "Max(s)", "Remark" }}
		csvRows = append(header, csvRows...)
		WriteToCSV(folderName, csvRows)
	}
}

func ReplaceFromDefaults(text *string, defaults map[string]string) {
	regex := regexp.MustCompile(`(?P<key>\${[A-Za-z0-9_-]+})+`)
	match := regex.FindAllString(*text, -1)

	for _, m := range match {
		dKey := strings.Replace(strings.Replace(m, "${", "", -1), "}", "", -1)
		*text = strings.Replace(*text, m, defaults[dKey], -1)
	}
}

func Stdout(requestName string, metrics vegeta.Metrics)  {
	fmt.Printf("========================\n")
	fmt.Printf("Name: %v\n", requestName)
	fmt.Printf("Requests: %v\n", metrics.Requests)
	fmt.Printf("Success: %v\n", metrics.Success)
	fmt.Printf("StatusCodes: %v\n", metrics.StatusCodes)
	fmt.Printf("ByteIn: %v\n", metrics.BytesIn.Mean)
	fmt.Printf("ByteOut: %v\n", metrics.BytesOut.Mean)
	fmt.Printf("Latencies[Mean,P95,P99,Max]: [%v, %v, %v, %v]\n",
		metrics.Latencies.Mean, metrics.Latencies.P95 , metrics.Latencies.P99, metrics.Latencies.Max)

	if len(metrics.Errors) > 0 {
		fmt.Printf("Errors:\n")
		errors := make(map[string]int)
		for _, e := range metrics.Errors {
			errors[e] =  errors[e] + 1
		}
		for k, v := range errors {
			fmt.Printf("%v: %v\n", k, v)
		}
	}
}

func FormatCSVName(folderName string) string {
	return fmt.Sprintf("%v/results.csv", folderName)
}

func WriteToCSV(folderName string, data [][]string) {
	csvFileName := FormatCSVName(folderName)

	file, err := os.Create(csvFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}
}

func FormatPlotName(folderName string, plotName string) string {
	return fmt.Sprintf("%v/%v.html", folderName, plotName)
}

func DrawPlot(folderName string, plotName string, results vegeta.Results) {
	fileName := FormatPlotName(folderName, plotName)

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	defer w.Flush()

	report := vegeta.NewPlotReporter(plotName, &results)
	report.Report(w)
}