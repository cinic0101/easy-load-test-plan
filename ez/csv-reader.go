package ez

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"strings"
)

func CSVReader(fileName string) []map[string]string {
	var c []map[string]string

	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	data := string(buf)

	r := csv.NewReader(strings.NewReader(data))
	headers, err := r.Read()

	for {
		rows, err := r.Read()
		if err == io.EOF {
			break
		}

		m := make(map[string]string)
		for i, key := range headers {
			m[key] = rows[i]
		}

		c = append(c, m)
	}

	return c
}
