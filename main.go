package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"code.google.com/p/gofpdf"
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
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 15)
	pdf.SetTopMargin(20)
	pdf.SetLeftMargin(20)
	pdf.SetRightMargin(20)
	pageW, _ := pdf.GetPageSize()
	left, top, right, _ := pdf.GetMargins()
	width := pageW - left - right

	initialColWidths := []float64{0.5 * width, 0.25 * width, 0.12 * width, 0.13 * width}
	initialColYs := []float64{
		left,
		left + initialColWidths[0],
		left + initialColWidths[0] + initialColWidths[1],
		left + initialColWidths[0] + initialColWidths[1] + initialColWidths[2]}

	pdf.Image("ComplyLogo.png", left, top, 0.4*width, 0, false, "", 0, "")
	pdf.Text(initialColYs[1], top+10, "QUALITY")
	_, fontSize := pdf.GetFontSize()
	pdf.SetFontSize(11)
	pdf.Text(initialColYs[1], top+10+fontSize, "Weekly Report")
	pdf.SetXY(initialColYs[2], top)
	pdf.MultiCell(initialColWidths[2], 2, "\nContract #\n\n\n\n\nCustomer\n\n\n\n\n", "1", "T", false)
	pdf.SetXY(initialColYs[3], top)
	pdf.CellFormat(initialColWidths[3], 7, "", "1", 2, "", false, 0, "")
	pdf.CellFormat(initialColWidths[3], 15, "", "1", 1, "", false, 0, "")
	// pdf.Ln(2)
	// pdf.SetFontSize()
	lineHeight := fontSize + 2
	pdf.CellFormat(0.22*width, lineHeight, "Make", "1", 0, "", false, 0, "")
	pdf.CellFormat(0.22*width, lineHeight, "Model", "1", 0, "", false, 0, "")
	pdf.CellFormat(0.22*width, lineHeight, "Serial #", "1", 0, "", false, 0, "")
	pdf.CellFormat(0.34*width, lineHeight, "Equipment #", "1", 1, "", false, 0, "")

	pdf.SetFontSize(8)
	pdf.Ln(1)
	pdf.Cell(0, 10, "Inspection: Check if acceptable in accordance with Compliance Services")
	pdf.Ln(lineHeight + 2)

	pdf.SetFontSize(11)
	pdf.CellFormat(0.5*width, lineHeight, "Week #:  ", "1", 0, "R", false, 0, "")
	numDivs := 12

	printBlanks := func(numDivs int, wdth float64) {
		for i := 0; i < numDivs; i++ {
			pdf.CellFormat(wdth/float64(numDivs), lineHeight, "", "1", 0, "", false, 0, "")
		}
	}

	printBlanks(numDivs, 0.5*width)
	pdf.Ln(lineHeight + 1)

	pdf.SetFontSize(10)
	printCheckOffLine := func(label string) {
		pdf.CellFormat(0.5*width, lineHeight, label, "1", 0, "", false, 0, "")
		printBlanks(numDivs, 0.5*width)
		pdf.Ln(lineHeight)
	}

	labels := []string{
		"All functions operate properly",
		"Structure and components secure and undamaged",
		"Safety devices intact and functional",
		"Operational and safety decals present and legible",
		"Batteries Charged",
		"Electrical cord/plug",
		"Containment tank\tQty Pumps:",
		"Appearance acceptable",
		"Fire Extinguisher\tQty:",
		"First Aid",
	}

	for _, label := range labels {
		printCheckOffLine(label)
	}

	pdf.SetFontSize(11)
	pdf.CellFormat(0.5*width, 6, "Manufacture/Model", "", 1, "BC", false, 0, "")
	pdf.Ln(1)

	pdf.SetFontSize(10)
	printCheckOffLine("AED")

	for _, label := range []string{"S/N", "Pad Exp.", "Battery Exp.", "Placed in Service"} {
		pdf.CellFormat(0.25*width, lineHeight, label, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(lineHeight)
	for i := 0; i < 4; i++ {
		pdf.CellFormat(0.25*width, lineHeight, "", "1", 0, "C", false, 0, "")
	}

	pdf.Ln(lineHeight + 4)

	pdf.SetFontSize(11)
	pdf.CellFormat(0.9*width, lineHeight-2, "Weekly Notes:", "", 0, "", false, 0, "")
	pdf.CellFormat(0.1*width, lineHeight-2, "Initial", "1", 1, "C", false, 0, "")

	for i := 1; i <= 12; i++ {
		pdf.CellFormat(0.1*width, lineHeight, fmt.Sprintf("Week %d", i), "1", 0, "", false, 0, "")
		pdf.CellFormat(0.8*width, lineHeight, "", "1", 0, "", false, 0, "")
		pdf.CellFormat(0.1*width, lineHeight, "", "1", 1, "", false, 0, "")
	}

	// pdf.MultiCell(initialColWidths[3], 2, "\nContract #\n\n\n\n\nCustomer\n\n\n\n\n", "1", "T", false)
	// pdf.Cell(initialColWidths[1],10,"Quality")
	// pdf.Cell(initialColWidths[1],10,"Weekly Report")
	// pdf.Text(50, 20, "logo.png")
	// pdf.Image(imageFile("logo.gif"), 10, 40, 30, 0, false, "", 0, "")
	// pdf.Text(50, 50, "logo.gif")
	// pdf.Image(imageFile("logo-gray.png"), 10, 70, 30, 0, false, "", 0, "")
	// pdf.Text(50, 80, "logo-gray.png")
	// pdf.Image(imageFile("logo-rgb.png"), 10, 100, 30, 0, false, "", 0, "")
	// pdf.Text(50, 110, "logo-rgb.png")
	// pdf.Image(imageFile("logo.jpg"), 10, 130, 30, 0, false, "", 0, "")
	// pdf.Text(50, 140, "logo.jpg")
	pdf.OutputFileAndClose("example.pdf")
}

func generatePdf2() {
	pdf := gofpdf.New("P", "mm", "A4", "../font")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello")
}
