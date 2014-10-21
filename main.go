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

type FormData1 struct {
	Trailer struct {
		FracCompany     string `json:"fracCompany"`
		Operator        string `json:"operator"`
		Location        string `json:"location"`
		Make            string `json:"make"`
		Model           string `json:"model"`
		SerialNumber    string `json:"serialNumber"`
		EquipmentNumber string `json:"equipmentNumber"`
	} `json:"trailer"`
	Inspections []struct {
		Date  string `json:"date"`
		Notes string `json:"notes"`
		By    string `json:"by"`
	} `json:"inspections"`
	AED struct {
		SerialNumber          string `json:"serialNumber"`
		PadExpirationDate     string `json:"padExpirationDate"`
		BatteryExpirationDate string `json:"batteryExpirationDate"`
		InServiceDate         string `json:"inServiceDate"`
	} `json:"AED"`
}

func main() {
	fmt.Printf("Hello, world.\n")
	// generatePdf1(parseJson("example.json", 1).(FormData1))
	generatePdf3(parseJson("form3data.json", 3).(Form3Data))
}

func parseJson(filename string, form int) interface{} {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Unable to read file", filename, "due to", err.Error())
		return FormData1{}
	}

	switch form {
	case 1:
		var data FormData1
		if err = json.Unmarshal(bytes, &data); err != nil {
			fmt.Printf("Unable to unmarshal JSON file due to", err.Error())
			return FormData1{}
		}
		return data
	case 3:
		var data Form3Data
		if err = json.Unmarshal(bytes, &data); err != nil {
			fmt.Printf("Unable to unmarshal JSON file due to", err.Error())
			return Form3Data{}
		}
		return data
	default:
		fmt.Println("An error occurred: unknown form type (", form, ")")
		return nil
	}
}

func writeText(txt string, width, height float64, placementAfter int, box bool, pdf *gofpdf.Fpdf) {
	initialFontSize, _ := pdf.GetFontSize()
	for fontSize := initialFontSize; pdf.GetStringWidth(txt) >= width-1.5; fontSize = fontSize - 0.1 {
		pdf.SetFontSize(fontSize)
	}
	boxStr := ""
	if box {
		boxStr = "1"
	}
	pdf.CellFormat(width, height, txt, boxStr, placementAfter, "C", false, 0, "")
	pdf.SetFontSize(initialFontSize)
}

func generatePdf1(data FormData1) {
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
	pdf.CellFormat(initialColWidths[2], 7, "Frac Co.", "1", 2, "C", false, 0, "")
	pdf.CellFormat(initialColWidths[2], 7, "Operator", "1", 2, "C", false, 0, "")
	pdf.CellFormat(initialColWidths[2], 7, "Location", "1", 1, "C", false, 0, "")
	pdf.SetXY(initialColYs[3], top)

	writeText(data.Trailer.FracCompany, initialColWidths[3], 7, 2, true, pdf)
	writeText(data.Trailer.Operator, initialColWidths[3], 7, 2, true, pdf)
	writeText(data.Trailer.Location, initialColWidths[3], 7, 1, true, pdf)

	lineHeight := fontSize + 2

	type MMSE struct {
		Label string
		Width float64
		Value string
	}

	mmseLine := []MMSE{
		{Label: "Make:", Width: 0.22 * width, Value: data.Trailer.Make},
		{Label: "Model:", Width: 0.22 * width, Value: data.Trailer.Model},
		{Label: "Serial #:", Width: 0.22 * width, Value: data.Trailer.SerialNumber},
		{Label: "Equipment #:", Width: 0.34 * width, Value: data.Trailer.EquipmentNumber},
	}

	for index, mmseElement := range mmseLine {
		labelWidth := pdf.GetStringWidth(mmseElement.Label)
		valueStart := pdf.GetX() + labelWidth
		valueWidth := mmseElement.Width - labelWidth

		pdf.CellFormat(mmseElement.Width, lineHeight, mmseElement.Label, "1", 0, "", false, 0, "")

		placementAfter := 0
		if index >= len(mmseLine)-1 {
			placementAfter = 1
		}
		pdf.SetX(valueStart)
		writeText(mmseElement.Value, valueWidth, lineHeight, placementAfter, false, pdf)
	}

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
		printCheckBoxes(numDivs, len(data.Inspections), 0.5*width, false)
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
	for _, value := range []string{data.AED.SerialNumber, data.AED.PadExpirationDate, data.AED.BatteryExpirationDate, data.AED.InServiceDate} {
		writeText(value, 0.25*width, lineHeight, 0, true, pdf)
		// pdf.CellFormat(0.25*width, lineHeight, "", "1", 0, "C", false, 0, "")
	}

	pdf.Ln(lineHeight + 4)

	pdf.SetFontSize(11)
	pdf.CellFormat(0.78*width, lineHeight-2, "Inspection Notes:", "", 0, "", false, 0, "")
	pdf.CellFormat(0.11*width, lineHeight-2, "Employee", "1", 0, "C", false, 0, "")
	pdf.CellFormat(0.11*width, lineHeight-2, "Date", "1", 1, "C", false, 0, "")

	for index, inspection := range data.Inspections {
		str := "Inspection #:"
		strWidth := pdf.GetStringWidth(str)
		cell1Width := 0.17 * width
		numberWidth := cell1Width - strWidth
		numberStart := pdf.GetX() + strWidth

		pdf.CellFormat(cell1Width, lineHeight, str, "1", 0, "", false, 0, "")
		pdf.SetX(numberStart)
		writeText(fmt.Sprintf("%d", index+1), numberWidth, lineHeight, 0, false, pdf)
		writeText(inspection.Notes, 0.61*width, lineHeight, 0, true, pdf)
		writeText(inspection.By, 0.11*width, lineHeight, 0, true, pdf)
		writeText(inspection.Date, 0.11*width, lineHeight, 1, true, pdf)
	}
	for i := len(data.Inspections) + 1; i <= 12; i++ {
		pdf.CellFormat(0.17*width, lineHeight, "Inspection #: ", "1", 0, "", false, 0, "")
		pdf.CellFormat(0.61*width, lineHeight, "", "1", 0, "", false, 0, "")
		pdf.CellFormat(0.11*width, lineHeight, "", "1", 0, "", false, 0, "")
		pdf.CellFormat(0.11*width, lineHeight, "", "1", 1, "", false, 0, "")
	}

	pdf.OutputFileAndClose("example1.pdf")
}

/*



THE SECOND PDF



*/
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

/*


Form 3

*/

type Form3Data struct {
	Trailer struct {
		CompanyNumber string `json:"companyNumber"`
		DOTNumber     string `json:"DOTNumber"`
		SerialNumber  string `json:"serialNumber"`
		Make          string `json:"make"`
		Location      string `json:"location"`
		Operator      string `json:"operator"`
		FracCompany   string `json:"fracCompany"`
	} `json:"trailer"`
	TruckNumber        string          `json:"truckNumber"`
	Odometer           int             `json:"odometer"`
	Remarks            string          `json:"remarks"`
	DriverSignatureUrl string          `json:"driverSignatureUrl"`
	Date               string          `json:"date"`
	Checklist          map[string]bool `json:"checklist"`
}

func generatePdf3(data Form3Data) {

	type TextPoint struct {
		Value string
		X     float64
		Y     float64
		Width float64
	}
	textPoints := map[string]TextPoint{
		"companyNumber": {data.Trailer.CompanyNumber, 50, 113.15, 137}, // trailer number(s)?
		// "dotNumber":     {data.Trailer.DOTNumber, 50, 116.5, 30},       // trailer number(s)?
		// "serialNumber":  {data.Trailer.SerialNumber, 50, 116.5, 30},    // trailer number(s)?

		// "make":          {data.Trailer.Make, 10, 10, 30},
		// "location":      {data.Trailer.Location, 10, 10, 30},
		// "operator":      {data.Trailer.Operator, 10, 10, 30},
		// "fracCompany":   {data.Trailer.FracCompany, 10, 10, 30},

		"truckNumber": {data.TruckNumber, 56, 62.7, 78},
		"odometer":    {fmt.Sprintf("%d", data.Odometer), 156, 62.7, 32},
		// "remarks":       {data.Remarks, 39, 142.5, 149},
		// "remarks2":      {"Remarks2 Remarks 2 Remarks2 Remarks 2 Remarks2 Remarks 2 Remarks2 Remarks 2r", 19, 151, 170},
		// "remarks3":      {"Remarks3 Remarks 3 Remarks3 Remarks 3 Remarks3 Remarks 3 Remarks3 Remarks 3rr", 19, 159, 170},
		"driverSigUrl1": {data.DriverSignatureUrl, 19, 238.5, 89},
		"driverSigUrl2": {data.DriverSignatureUrl, 19, 171.5, 89},
		"date":          {data.Date, 130, 238.5, 54.6},
		"date2":         {data.Date, 130, 171.5, 54.6},
	}

	writeRemarks := func(remarks string, pdf *gofpdf.Fpdf) {
		remarksLen := pdf.GetStringWidth(remarks)
		const (
			line1Len float64 = 147
			line2Len float64 = 168
			line3Len float64 = 168
			total    float64 = line1Len + line2Len + line3Len
		)

		oldFontSize, lineHeight := pdf.GetFontSize()
		for fontSize := oldFontSize; remarksLen > total; fontSize = fontSize - 1 {
			pdf.SetFontSize(fontSize)
			remarksLen = pdf.GetStringWidth(remarks)
		}

		breakIndex1 := int((line1Len / remarksLen) * float64(len(remarks)))
		breakIndex2 := int(((line1Len + line2Len) / remarksLen) * float64(len(remarks)))

		pdf.SetXY(39, 142.5)
		if remarksLen > line1Len+line2Len {
			//Need 3 lines
			remarks1 := remarks[:breakIndex1]
			remarks2 := remarks[breakIndex1:breakIndex2]
			remarks3 := remarks[breakIndex2:]
			pdf.CellFormat(line1Len, lineHeight, remarks1, "", 0, "", false, 0, "")
			pdf.SetXY(19, 151)
			pdf.CellFormat(line2Len, lineHeight, remarks2, "", 0, "", false, 0, "")
			pdf.SetXY(19, 159)
			pdf.CellFormat(line3Len, lineHeight, remarks3, "", 0, "", false, 0, "")
		} else if remarksLen > line1Len {
			//Need 2 lines
			remarks1 := remarks[:breakIndex1]
			remarks2 := remarks[breakIndex1:]
			pdf.CellFormat(line1Len, lineHeight, remarks1, "", 0, "", false, 0, "")
			pdf.SetXY(19, 151)
			pdf.CellFormat(line2Len, lineHeight, remarks2, "", 0, "", false, 0, "")
		} else {
			//Only need one line
			pdf.CellFormat(line1Len, lineHeight, remarks, "", 0, "", false, 0, "")
		}

		pdf.SetFontSize(oldFontSize)
	}

	type Point struct {
		X float64
		Y float64
	}
	checks := map[string]Point{
		"airCompressorTruck":   {18.5, 68},   //Air compressor
		"airLinesTruck":        {18.5, 72.5}, //Air lines
		"parkingBrakesTruck":   {18.5, 77},   //Brakes: Parking
		"serviceBrakesTruck":   {18.5, 81.5}, //Brakes: Service Brakes
		"batteryTruck":         {18.5, 86},   //Battery
		"couplingDevicesTruck": {18.5, 90.5}, //Coupling devices
		"emergencyTruck":       {18.5, 95},   //Emergency equipment
		"fuelTanksTruck":       {18.5, 99.5}, //Fuel tanks
		"hornTruck":            {18.5, 104},  //Horn

		"lightingTruck":    {132.5, 68.1}, //Lighting devices
		"mirrorsTruck":     {132.5, 72.6}, //Mirrors
		"oilPressureTruck": {132.5, 77.1}, //Oil Pressure
		"steeringTruck":    {132.5, 81.6}, //Steering mech
		"tiresTruck":       {132.5, 86},   //tires
		"transTruck":       {132.5, 90.5}, //Transmission
		"wheelsTruck":      {132.5, 95},   //Wheels and Rims
		"windshieldTruck":  {132.5, 99.5}, //Windshield wipers
		"otherTruck":       {132.5, 104},  //Other

		"brakesTrailer":          {18.5, 118.5}, //Brakes
		"couplingDevicesTrailer": {18.5, 123},   //Coupling devices
		"couplingPinTrailer":     {18.5, 127.5}, //Coupling pin
		"doorsTrailer":           {18.5, 132},   //Doors
		"hitchTrailer":           {18.5, 136.5}, //Hitch

		"landingGearTrailer": {88.1, 118.7}, //Landing gear
		"lightsTrailer":      {88.1, 123.2}, //Lights
		"roofTrailer":        {88.1, 127.7}, //Roof
		"tarpTrailer":        {88.1, 132.2}, //Tarpaulin
		"tiresTrailer":       {88.1, 136.7}, //Tires

		"springTrailer": {158, 118.7}, //Springs
		"wheelsTrailer": {158, 123.2}, //Wheels
		"otherTrailer":  {158, 127.7}, //Other

		"conditionSat": {18, 162.3}, //Condition of the above vehicle is satisfactory

		"defectsCorrected":   {18, 207},   //Above defects corrected
		"defectsUncorrected": {18, 212.5}, //Above defects need not be corrected
		"noDefects":          {18, 218},   //No defects to report

	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	pageW, pageH := pdf.GetPageSize()
	pdf.Image("form3raw.png", 0, 0, pageW, 0, false, "", 0, "")

	_, lineH := pdf.GetFontSize()

	for _, point := range textPoints {
		pdf.SetXY(point.X, point.Y)
		writeText(point.Value, point.Width, lineH, 0, false, pdf)
		// pdf.CellFormat(point.Width, lineH, point.Value, "1", 0, "C", false, 0, "")
	}

	writeRemarks(data.Remarks, pdf)
	pdf.SetFontSize(6)
	for x := 0.0; x <= pageW; x = x + 10.0 {
		for y := 0.0; y <= pageH; y = y + 10.0 {
			// pdf.Text(x, y, fmt.Sprintf("(%.f,%.f)", x, y))
		}
	}

	pdf.SetFontSize(8)
	for key, value := range data.Checklist {
		if value {
			mark := checks[key]
			pdf.Image("check_mark.png", mark.X, mark.Y, 5, 0, false, "", 0, "")
		}
	}
	// for _, point := range checks {
	// 	pdf.Image("check_mark.png", point.X, point.Y, 5, 0, false, "", 0, "")
	// }

	pdf.OutputFileAndClose("example3.pdf")
}
