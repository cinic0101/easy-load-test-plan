package main
import (
	"io/ioutil"
	"fmt"
	"net/http"
	"time"
	"os"

	"github.com/tsenart/vegeta/lib"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Rate uint64
	Duration time.Duration
	Defaults map[string]string
	Requests map[string]struct {
		Method string
		Url string
		Headers map[string]string
		Body string
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("No yaml file")
		return
	}
	fmt.Println("Parsing file: " + os.Args[1])

	buf, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	data := string(buf)

	c := Config{}
	if err = yaml.Unmarshal([]byte(data), &c); err != nil {
		panic(err)
	}

	for r := range c.Requests {
		rate := c.Rate
		duration := c.Duration * time.Second

		header := http.Header{}
		for k := range c.Requests[r].Headers {
			value := c.Requests[r].Headers[k]
			first := string(value[0])

			if first == "$" {
				dKey := value[1:]
				header.Add(k, c.Defaults[dKey])
			} else {
				header.Add(k, c.Requests[r].Headers[k])
			}
		}

		target := vegeta.NewStaticTargeter(vegeta.Target{
			Method: c.Requests[r].Method,
			URL:    c.Requests[r].Url,
			Header: header,
			Body:   []byte(c.Requests[r].Body),
		})
		attacker := vegeta.NewAttacker()

		var metrics vegeta.Metrics
		for res := range attacker.Attack(target, rate, duration) {
			metrics.Add(res)
		}
		metrics.Close()

		fmt.Printf("\n========================\n")
		fmt.Printf("Name: %s%v%s\n", CLR_R, r, CLR_N)
		fmt.Printf("Requests: %v\n", metrics.Requests)
		fmt.Printf("Success: %s%v%s\n", success(metrics.Success), metrics.Success, CLR_N)
		fmt.Printf("StatusCodes: %v\n", metrics.StatusCodes)
		fmt.Printf("Latencies[Mean,P95,P99,Max]: [%s%v%s, %v, %v, %v]\n",
			CLR_Y, metrics.Latencies.Mean, CLR_N, metrics.Latencies.P95 , metrics.Latencies.P99, metrics.Latencies.Max)
		fmt.Printf("Errors: \n")
		for _, e := range metrics.Errors {
			fmt.Printf("%v\n", e)
		}
	}
}

const CLR_N = "\x1b[0m"
const CLR_R = "\x1b[31;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_C = "\x1b[36;1m"

func success(success float64) string {
	if success == 1 {
		return CLR_C
	}

	return CLR_R
}