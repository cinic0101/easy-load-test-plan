package main

import (
	"os"
	"encoding/csv"
	"github.com/nu7hatch/gouuid"
)

func main() {
	file, err := os.Create("temp6.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"temp_id"})
	for i := 0; i < 300 * 60; i++ {
		u, _ := uuid.NewV4()

		if err := writer.Write([]string{u.String()}); err != nil {
			panic(err)
		}
	}
}
