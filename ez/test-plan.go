package ez

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type TestPlan struct {
	Rate uint64
	Duration time.Duration
	Result struct {
		Stdout bool
		CSV    bool
		Plot   bool
	}
	Defaults map[string]string
	Requests map[string]Request
}

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []string
	Remark  string
}

func NewTestPlan(file string) *TestPlan {
	p := &TestPlan {
		Result: struct {
			Stdout bool
			CSV    bool
			Plot   bool
		}{Stdout: true, CSV: true, Plot: true},
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(buf, &p); err != nil {
		panic(err)
	}

	return p
}