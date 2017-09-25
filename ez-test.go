package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tsenart/vegeta/lib"
	"gopkg.in/yaml.v2"
	"regexp"
	"bufio"
)

type TestPlan struct {
	Rate uint64
	Duration time.Duration
	Result struct {
		Stdout bool
		Csv bool
		Plot bool
	}
	Defaults map[string]string
	Requests map[string]Request
}

type Request struct {
	Method string
	Url string
	Headers map[string]string
	Body []string
	Remark string
}

func main() {
	if len(os.Args) == 1 {
		panic("No yaml file")
	}

	planName := os.Args[1]
	buf, err := ioutil.ReadFile(planName)
	if err != nil {
		panic(err)
	}

	p := TestPlan{}
	if err = yaml.Unmarshal(buf, &p); err != nil {
		panic(err)
	}

	var csvRows [][]string
	for requestName, r := range p.Requests {
		rate := p.Rate
		duration := p.Duration * time.Second

		header := http.Header{}
		for hk := range r.Headers {
			value := r.Headers[hk]
			length := len(value)
			first := string(value[0])
			last := string(value[length-1])

			if first == last && last == "$" {
				dKey := value[1:length-1]

				if val, ok := p.Defaults[dKey]; ok {
					header.Add(hk, val)
				} else {
					panic(fmt.Sprintf("Can not find key \"%s\" in defaults.", dKey))
				}
			} else {
				header.Add(hk, r.Headers[hk])
			}
		}

		ReplaceFromDefaults(&r.Url, p.Defaults)

		var targets []vegeta.Target
		if len(r.Body) > 0 {
			for _, b := range r.Body {
				ReplaceFromDefaults(&b, p.Defaults)
				targets = append(targets, vegeta.Target{
					Method: r.Method,
					URL:    r.Url,
					Header: header,
					Body:   []byte(b),
				})
			}
		} else {
			targets = append(targets, vegeta.Target{
				Method: r.Method,
				URL:    r.Url,
				Header: header,
			})
		}

		target := vegeta.NewStaticTargeter(targets...)
		attacker := vegeta.NewAttacker()

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

		if p.Result.Csv {
			csvRows = append(csvRows, []string{
				requestName,
				r.Url,
				fmt.Sprintf("%.1f", metrics.Success * 100),
				fmt.Sprintf("%.3f", metrics.Latencies.Mean.Seconds()),
				fmt.Sprintf("%.3f", metrics.Latencies.Max.Seconds()),
			})
		}

		if p.Result.Plot {
			DrawPlot(planName, requestName, results)
		}
	}

	if p.Result.Csv {
		header := [][]string{{ "Name", "Url", "Success(%)", "Mean(s)", "Max(s)", "Remark" }}
		csvRows = append(header, csvRows...)
		WriteToCsv(planName, csvRows)
	}
}

func ReplaceFromDefaults(text *string, defaults map[string]string) {
	regex := regexp.MustCompile(`(?P<key>\$[A-Za-z0-9_-]+\$)+`)
	match := regex.FindAllString(*text, -1)

	for _, m := range match {
		dKey := strings.Replace(m, "$", "", -1)
		*text = strings.Replace(*text, m, defaults[dKey], -1)
	}
}

func Stdout(requestName string, metrics vegeta.Metrics)  {
	fmt.Printf("========================\n")
	fmt.Printf("Name: %v\n", requestName)
	fmt.Printf("Requests: %v\n", metrics.Requests)
	fmt.Printf("Success: %v\n", metrics.Success)
	fmt.Printf("StatusCodes: %v\n", metrics.StatusCodes)
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

func WriteToCsv(planName string, data [][]string) {
	csvFileName := fmt.Sprintf("%v.csv", strings.Replace(planName, ".yml", "", 1))

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

func DrawPlot(planName string, plotName string, results vegeta.Results) {
	fileName := fmt.Sprintf("%v_%v.html", strings.Replace(planName, ".yml", "", 1), plotName)

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	defer w.Flush()

	report := vegeta.NewPlotReporter(plotName, &results)
	report.Report(w)
}