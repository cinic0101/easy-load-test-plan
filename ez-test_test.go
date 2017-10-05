package main

import (
	"os"
	"testing"

	"github.com/tsenart/vegeta/lib"
)

func TestReplaceFromDefaults(t *testing.T) {
	text := "$a$,$b$,$c$"
	defaults := map[string]string{ "a": "1", "b": "2", "c": "3"}

	ReplaceFromDefaults(&text, defaults)

	if text != "1,2,3" {
		t.Error("Replace fail")
	}
}

func TestWriteToCsv(t *testing.T) {
	planName := "test"
	WriteToCsv(planName, [][]string{{"h1", "h2"}, {"r1c1", "r1c2"}})
	os.Remove(FormatCsvName(planName))
}

func TestDrawPlot(t *testing.T) {
	planName := "test"
	plotName := "plot"
	DrawPlot(planName, plotName, vegeta.Results{})
	os.Remove(FormatPlotName(planName, plotName))
}
