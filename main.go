package main

import (
	"code.google.com/p/gofpdf"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//TODO: JSON
/*

{
   "trailer":{
      "companyNumber":"1",
      "DOTNumber":"123",
      "serialNumber":"123",
      "make":"Black Ford 2013",
      "location":"Winter Pad",
      "operator":"Operator Inc.",
      "fracCompany":"Fracking Co."
   },
   "inspections":[
      {
         "date":"(parse date format, don’t know)",
         "notes":"I just made a sandwich i didn’t do my job",
         "by":"Logan Spears"
      }
   ],
   "AED":{
      "serialNumber":"aaa",
      "padExpirationDate":"(parse date format, don’t know)",
      "batteryExpirationDate":"(parse date format, don’t know)",
      "inServiceDate":"(parse date format, don’t know)"
   }
}

*/

type trailerObj struct {
	CompanyNumber string
	DOTNumber     string
	SerialNumber  string
	Make          string
	Location      string
	Operator      string
	FracCompany   string
}

type inspectionObj struct {
	Date  string
	Notes string
	By    string
}

type aedObj struct {
	SerialNumber          string
	PadExpirationDate     string
	BatteryExpirationDate string
	InServiceDate         string
}

type FormData1 struct {
	Trailer     trailerObj
	Inspections []inspectionObj
	AED         aedObj
}

func main() {
	fmt.Printf("Hello, world.\n")
	parseJson("example.json")
	generatePdf1()
}

func parseJson(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Unable to read file", filename, "due to", err.Error())
		return
	}

	var doc1 FormData1
	if err = json.Unmarshal(bytes, &doc1); err != nil {
		fmt.Printf("Unable to unmarshal JSON file due to", err.Error())
		return
	}

	fmt.Println("Unmarshaling succesful; value obtained:\n", doc1)

}

func generatePdf1() {
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
	pdf.Text(initialColYs[1], top+10+fontSize, "Inspection Report")
	pdf.SetXY(initialColYs[2], top)
	pdf.MultiCell(initialColWidths[2], 2, "\nFrac Co.\n\n\n\n\nOperator\n\n\n\n\n", "1", "T", false)
	pdf.SetXY(initialColYs[3], top)
	pdf.CellFormat(initialColWidths[3], 7, "", "1", 2, "", false, 0, "")
	pdf.CellFormat(initialColWidths[3], 15, "", "1", 1, "", false, 0, "")
	// pdf.Ln(2)
	// pdf.SetFontSize()
	lineHeight := fontSize + 2
	pdf.CellFormat(0.22*width, lineHeight, "Make:", "1", 0, "", false, 0, "")
	pdf.CellFormat(0.22*width, lineHeight, "Model:", "1", 0, "", false, 0, "")
	pdf.CellFormat(0.22*width, lineHeight, "Serial #:", "1", 0, "", false, 0, "")
	pdf.CellFormat(0.34*width, lineHeight, "Equipment #:", "1", 1, "", false, 0, "")

	pdf.SetFontSize(8)
	pdf.Ln(1)
	pdf.Cell(0, 10, "Inspection: Check if acceptable in accordance with Compliance Services")
	pdf.Ln(lineHeight + 2)

	pdf.SetFontSize(11)
	pdf.CellFormat(0.5*width, lineHeight, "Inspection #:  ", "1", 0, "R", false, 0, "")
	numDivs := 12

	printCheckBoxes := func(numDivs, numChecked int, wdth float64, useNumbersInsteadOfChecks bool) {
		for i := 0; i < numDivs; i++ {
			str := ""
			if i < numChecked {
				if useNumbersInsteadOfChecks {
					str = fmt.Sprintf("%d", i+1)
				} else {
					pdf.Image("check_mark.png", pdf.GetX(), pdf.GetY(), 0.5*width/float64(numDivs), 0, false, "", 0, "")
				}
			}
			pdf.CellFormat(wdth/float64(numDivs), lineHeight, str, "1", 0, "C", false, 0, "")
		}
	}

	printCheckBoxes(numDivs, 12, 0.5*width, true)
	pdf.Ln(lineHeight + 1)

	pdf.SetFontSize(10)
	printCheckOffLine := func(label string) {
		pdf.CellFormat(0.5*width, lineHeight, label, "1", 0, "", false, 0, "")
		printCheckBoxes(numDivs, 12, 0.5*width, false)
		pdf.Ln(lineHeight)
	}

	labels := []string{
		"All functions operate properly",
		"Structure and components secure and undamaged",
		"Shower and eyewash intact and functional",
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
	pdf.Ln(2)

	for _, label := range []string{"S/N", "Pad Exp.", "Battery Exp.", "Placed in Service"} {
		pdf.CellFormat(0.25*width, lineHeight, label, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(lineHeight)
	for i := 0; i < 4; i++ {
		pdf.CellFormat(0.25*width, lineHeight, "", "1", 0, "C", false, 0, "")
	}

	pdf.Ln(lineHeight + 4)

	pdf.SetFontSize(11)
	pdf.CellFormat(0.78*width, lineHeight-2, "Inspection Notes:", "", 0, "", false, 0, "")
	pdf.CellFormat(0.11*width, lineHeight-2, "Employee", "1", 0, "C", false, 0, "")
	pdf.CellFormat(0.11*width, lineHeight-2, "Date", "1", 1, "C", false, 0, "")

	for i := 1; i <= 12; i++ {
		pdf.CellFormat(0.17*width, lineHeight, "Inspection #: ", "1", 0, "", false, 0, "")
		pdf.CellFormat(0.61*width, lineHeight, "", "1", 0, "", false, 0, "")
		pdf.CellFormat(0.11*width, lineHeight, "", "1", 0, "", false, 0, "")
		pdf.CellFormat(0.11*width, lineHeight, "", "1", 1, "", false, 0, "")
	}

	pdf.OutputFileAndClose("example1.pdf")
}

func generatePdf2() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 8)
	pdf.SetTopMargin(10)
	pdf.SetLeftMargin(10)
	pdf.SetRightMargin(10)
	pdf.SetAutoPageBreak(true, 10)

	pageW, _ := pdf.GetPageSize()
	left, top, right, _ := pdf.GetMargins()
	width := pageW - left - right

	initialColWidths := []float64{0.5 * width, 0.25 * width, 0.25 * width}
	initialColXs := []float64{
		left,
		left + initialColWidths[0],
		left + initialColWidths[0] + initialColWidths[1]}

	_, fontSize := pdf.GetFontSize()
	pdf.SetXY(initialColXs[2], top)
	lineHeight := 5 * fontSize
	pdf.CellFormat(initialColWidths[2], lineHeight, "Pad Name", "1", 2, "LM", false, 0, "")
	pdf.CellFormat(initialColWidths[2], lineHeight, "Operator", "1", 1, "LM", false, 0, "")

	pdf.Image("ComplyLogo.png", left, top, 0.4*width, 0, false, "", 0, "")
	pdf.SetXY(initialColXs[1], top)
	pdf.SetFontSize(14)
	pdf.CellFormat(initialColWidths[2], lineHeight, "QUALITY", "", 2, "CB", false, 0, "")
	pdf.SetFontSize(10)
	pdf.CellFormat(initialColWidths[2], lineHeight, "CONDITION REPORT", "", 1, "L", false, 0, "")

	pdf.SetFontSize(8)
	lineHeight = fontSize * 3

	pdf.CellFormat(0.25*width, lineHeight, "Make", "1", 0, "L", false, 0, "")
	pdf.CellFormat(0.25*width, lineHeight, "Model", "1", 0, "L", false, 0, "")
	pdf.CellFormat(0.25*width, lineHeight, "Serial #", "1", 0, "L", false, 0, "")
	pdf.CellFormat(0.25*width, lineHeight, "Equip #", "1", 1, "L", false, 0, "")

	lineHeight = fontSize * 2
	pdf.CellFormat(width, lineHeight, "MARK CLEARLY ALL DAMAGE BY SYMBOL        \"C\"=CUT    \"H\"=HOLE    \"D\"=DENT    \"P\"=PAINT DAMAGE", "1", 1, "C", false, 0, "")

	colWidth := 0.49 * width
	colXs := []float64{left, left + width - colWidth}

	lineHeight = fontSize * 2
	pdf.SetFontSize(7)
	pdf.CellFormat(colWidth, lineHeight, "RENTAL OUT - Diagram & Identify any existing damages or strike-out box", "1", 0, "C", false, 0, "")
	pdf.SetX(colXs[1])
	pdf.CellFormat(colWidth, lineHeight, "RENTAL IN - Diagram & Identify any damages upon return", "1", 1, "C", false, 0, "")

	pdf.CellFormat(colWidth, lineHeight*8, "", "1", 0, "C", false, 0, "")
	pdf.SetX(colXs[1])
	pdf.CellFormat(colWidth, lineHeight*8, "", "1", 1, "C", false, 0, "")

	// lineHeight = lineHeight * 1.2

	checkboxWidth := 0.04 * width
	pdf.CellFormat(colWidth-2*checkboxWidth, lineHeight, "Inspection: Check if acceptable in accordance with Company spec", "1", 0, "", false, 0, "")
	pdf.CellFormat(checkboxWidth, lineHeight, "OUT", "1", 0, "C", false, 0, "")
	pdf.CellFormat(checkboxWidth, lineHeight, "IN", "1", 0, "C", false, 0, "")
	pdf.SetX(colXs[1])
	pdf.CellFormat(colWidth, lineHeight, "Comments - IN (if damaged, provide details of occurence)", "1", 1, "", false, 0, "")

	items := []string{
		"Structure and components secure and undamaged",
		"All functions operate properly",
		"Safety devices intact and functional",
		"Operational and safety decals present and legible",
		"Batteries charged",
		"Electrical cord/plug",
		"Containment tank      Qty pumps:",
		"Appearance Acceptable - note and mark exceptions in diagram",
		"Fire extinguisher        Qty:",
		"First Aid",
		"AED",
		"",
		"",
	}

	for _, item := range items {
		pdf.CellFormat(colWidth-2*checkboxWidth, lineHeight, item, "1", 0, "", false, 0, "")
		pdf.CellFormat(checkboxWidth, lineHeight, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(checkboxWidth, lineHeight, "", "1", 0, "C", false, 0, "")
		pdf.SetX(colXs[1])
		pdf.CellFormat(colWidth, lineHeight, "", "1", 1, "", false, 0, "")
	}

	weirdLine := func(col1Text, col2Text string, lnHeight float64) {
		pdf.CellFormat(colWidth, lnHeight, col1Text, "1", 0, "", false, 0, "")
		pdf.SetX(colXs[1])
		pdf.CellFormat(colWidth, lnHeight, col2Text, "1", 1, "", false, 0, "")
	}

	weirdLine("Towable equipment-Record license plate of towing vehicle:", "", lineHeight)
	weirdLine("Comments - OUT ( list additional equipment )", "", lineHeight)
	for i := 0; i < 8; i++ {
		weirdLine("", "", lineHeight)
	}

	tab := func(wdth float64) string {
		spaceW := pdf.GetStringWidth(" ")
		numSpaces := wdth / spaceW
		tabStr := ""
		for i := 0.0; i < numSpaces; i += 1 {
			tabStr += " "
		}
		return tabStr
	}

	weirdLine(
		"Inspected and condtion OUT confirmed as specified above:",
		"Inspected and condition IN confirmed as specified above:", lineHeight)

	twoTabbedWeirdLine := func(col1First, col1Sec, col2First, col2Sec string, lnHeight float64) {
		col1FirstW := pdf.GetStringWidth(col1First)
		col2FirstW := pdf.GetStringWidth(col2First)

		tabStart := 0.55 * colWidth
		col1TabWidth := tabStart - col1FirstW
		col2TabWidth := tabStart - col2FirstW

		weirdLine(
			col1First+tab(col1TabWidth)+col1Sec,
			col2First+tab(col2TabWidth)+col2Sec,
			lnHeight)
	}

	twoTabbedWeirdLine("Compliance Services by", "Print Name", "Compliance Services by", "Print Name", lineHeight)
	twoTabbedWeirdLine("Date", "Date Delivered", "Date", "Date Received", lineHeight)

	// _, multiCellLineHeight := pdf.GetFontSize()
	// pdf.CellFormat(width, lineHeight,
	// 	"CUSTOMER ACKNOWLEDGEMENT: THE SAFETY AND PERFORMANCE OF THIS EQUIPMENT"+
	// 		" HAS BEEN VERIFIED. AS USER OF THIS EQUIPMENT, I UNDERSTAND THE CORRECT"+
	// 		" OPERATION AND FUNCTION OF THE CONTROLS AND CONFIRM THAT I HAVE RECEIVED"+
	// 		" ADEQUATE INSTRUCTION AND HAVE ADHERED TO THE SAFETY SHEET, THUS ENABLING MYSELF"+
	// 		" AND/OR MY CREW TO USE THE EQUIPMENT IN A SAFE AND PROPER MANNER WITHOUT RISK OF INJURY.\n",
	// 	"1", 1, "", false, 0, "")
	pdf.MultiCell(width, 3,
		"\nCUSTOMER ACKNOWLEDGEMENT: THE SAFETY AND PERFORMANCE OF THIS EQUIPMENT"+
			" HAS BEEN VERIFIED. AS USER OF THIS EQUIPMENT, I UNDERSTAND THE CORRECT"+
			" OPERATION AND FUNCTION OF THE CONTROLS AND CONFIRM THAT I HAVE RECEIVED"+
			" ADEQUATE INSTRUCTION AND HAVE ADHERED TO THE SAFETY SHEET, THUS ENABLING MYSELF"+
			" AND/OR MY CREW TO USE THE EQUIPMENT IN A SAFE AND PROPER MANNER WITHOUT RISK OF INJURY.\n\n",
		"1", "", false)

	weirdLine("Customer by", "Customer by", lineHeight)
	twoTabbedWeirdLine("Print Name", "Date", "Print Name", "Date", lineHeight)

	pdf.OutputFileAndClose("example2.pdf")
}
