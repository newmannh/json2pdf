package main

import (
	"code.google.com/p/gofpdf"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var bob = "bob"

func main() {
	fmt.Printf("Hello, world.\n")
	parseJson("example.json")
	generatePdf()
}

func parseJson(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Unable to read file", filename, "due to", err.Error())
		return
	}

	var v interface{}
	if err = json.Unmarshal(bytes, &v); err != nil {
		fmt.Printf("Unable to unmarshal JSON file due to", err.Error())
		return
	}

	fmt.Println("Unmarshaling succesful; value obtained:\n", v)

}

func generatePdf() {
	pdf := gofpdf.New("P", "mm", "A4", "../font")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, whirrrl")
	pdf.OutputFileAndClose("example.pdf")
}
