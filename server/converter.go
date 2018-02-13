package server

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

/*
	Muuttuja 1	Sarana
	Muuttuja 2	Kätisyys
	Muuttuja 3	Vedin Reikäväli
	Muuttuja 4	Vedin Asento
	Muuttuja 5	Saranan Paikka
	Muuttuja 6	Ovimallin 1/0 muuttuja

	Barcode 1
	Barcode 2	ovikoodi/ohjelma	var 5
	Barcode 3	korkeus
	Barcode 4	leveys
	Barcode 5	kätisyys			var 2
	Barcode 6	sarana				var 1
	Barcode 7	vedin reikäväli		var 3
	Barcode 8	vedin asento		var 4
*/

/*
	TPA\ALBATROS\CNC2000
	;test.TCN;;N;1;1;1;0;0;0;0;0;1955;595;16;0;0;;0;;0;0;;#0=;#1=;#2=;#3=;#4=;#5=;#6=;#7=;#0=;#1=;#2=;#3=;#4=;#5=;#6=;#7=;#2=96;#3=33;#4=h-33;;;;;
	;;;;;0;0;0;;0;0;0;0;;;;0;0;;0;;0;0;;;;;;;;;;;;;;;;;;;;;;;;;
	$=

	First Line CNC2000 List Headline

	1-Program Name		Alphanumeric Line Type 32 Types
	2-Comment			Alphanumeric Line Type 32 Types
	3-Work Area			Only N (Not Used)
	4-Repetition Rep 0	Integer (Program Repetition Number)
	5-Repetition M2		No 0 Default
	6-Repetition M3 	No 0 Default
	7-Repetition M4 	No 0 Default
	8-Repetition M5 	No 0 Default
	9-Repetition M6 	No 0 Default
	10-Repetition M7 	No 0 Default
	11-Repetition M8 	No 0 Default
	12-Length			Alphanumeric Line Type 32 Types
	13-Height			Alphanumeric Line Type 32 Types
	14-Thickness		Alphanumeric Line Type 32 Types
	15-Free				0 Default
	16-Free				0 Default
	17-Status			Empty, reserved column DO NOT FILL IN
	18-Free				0 Default
	19-Unload Type 		0 Default
	20-Free				0 Default
	21-Free 			0 Default
	22-Reserved 		Empty, reserved column DO NOT FILL IN
	23-38				8 more variables possible
	39-> 				Rewritable Variables R passed by the Program preceded by #variablenumber

*/

func Convert(r io.Reader, db *Client) (*bytes.Buffer, error) {
	cData, err := db.GetAll()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	zipWriter := zip.NewWriter(buf)

	scanner := bufio.NewScanner(r)

	// scan header line
	scanner.Scan()

	var w io.Writer
	var order, result string
	// Convert data row by row
	for scanner.Scan() {
		rawRow := scanner.Text()
		if rawRow == "" {
			continue
		}

		row := make([]string, 47)

		// Add defaults
		row[3] = "N"
		row[4] = "1"
		row[5] = "1"
		row[6] = "1"
		row[7] = "0"
		row[8] = "0"
		row[9] = "0"
		row[10] = "0"
		row[11] = "0"
		// 12-14 Length, Height, Thickness
		row[15] = "0"
		row[16] = "0"
		row[17] = ""
		row[18] = "0"
		row[19] = ""
		row[20] = "0"
		row[21] = "0"
		row[22] = ""

		// Init variables
		for i := 0; i < 8; i++ {
			row[23+i] = fmt.Sprintf("#%d=", i)
			row[31+i] = fmt.Sprintf("#%d=", i)
		}

		cols := strings.Split(rawRow, ";")
		if len(cols) != 29 {
			log.Println("Invalid amount of cols:", len(cols), rawRow)
			continue
		}

		finish := func(readyRow []string) {
			result += "\r\n"
			result += strings.Join(readyRow, ";")
		}

		// New ordernumber
		if cols[3] != "" && cols[3] != order {
			if w != nil {
				// Finish order file
				result += "\r\n;;;;0;0;0;;0;0;0;0;;;;0;0;;0;;0;0;;;;;;;;;;;;;;;;;;;;;;;;;\r\n"
				result += "$=\r\n"
				w.Write([]byte(result))
			}

			result = `TPA\ALBATROS\CNC2000`
			order = cols[3]
			if w, err = zipWriter.Create(order + ".lsc"); err != nil {
				return nil, err
			}
		}

		// Make sure we have order to write
		if w == nil {
			log.Println("Row without order:", rawRow)
			continue
		}

		// Length, height, thicknes
		row[12] = cols[7]
		row[13] = cols[8]

		// Door model
		if dm, ok := cData.DoorModels[cols[4]]; ok {
			row[44] += fmt.Sprintf("#6=%d", dm.Var6)
			row[14] = fmt.Sprintf("%d", dm.Depth)
			row[2] = cols[4]
		} else {
			log.Println("Could not found doormodel for:", cols[4], rawRow)
			row[2] += " OM"
		}

		// Program code
		if prog, ok := cData.Programs[cols[5]]; ok {
			row[1] = prog.Program
			row[43] += fmt.Sprintf("#5=%.1f", prog.HingePosition)
		} else {
			log.Println("Could not found program for:", cols[5], rawRow)
			row[2] += " O"
		}

		// Loop each possible barcode position from row
		// and create row per barcode in result set
		found := false
		for _, i := range []int{11, 17, 22, 27} {
			if cols[i] == "" {
				continue
			}

			found = true
			newRow := make([]string, len(row))
			copy(newRow, row)

			// Parse number of pieces
			amount, err := strconv.Atoi(cols[i-1])
			if err != nil {
				log.Println("Invalid number of pieces:", err)
				newRow[2] += " KPL"
			} else {
				newRow[4] = strconv.Itoa(amount)
			}

			if cols[i+1] != "" {
				log.Println("Excluded barcode", cols[i], rawRow)
				newRow[2] += " B!"
			} else if len(cols[i]) != 10 {
				log.Println("Invalid lenght barcode", len(cols[i]), cols[i], rawRow)
				newRow[2] += " BL"
			} else {
				barcode := cols[i]

				if sarana, ok := cData.Hinges[barcode[6]]; ok {
					newRow[39] = fmt.Sprintf("#1=%d", sarana.Var5)
				} else {
					log.Println("Saranointia ei löytynyt:", string(barcode[6]), barcode, rawRow)
					newRow[2] += " S"
				}

				if katisyys, ok := cData.Handednesses[barcode[5]]; ok {
					newRow[40] = fmt.Sprintf("#2=%s", katisyys.Handedness)
				} else {
					log.Println("Kätisyyttä ei löytynyt:", string(barcode[5]), barcode, rawRow)
					newRow[2] += " K"
				}

				if vedin, ok := cData.Handles[barcode[7]]; ok {
					newRow[41] = fmt.Sprintf("#3=%d", vedin.Handle)
				} else {
					log.Println("Reikäväliä ei löytynyt:", string(barcode[7]), barcode, rawRow)
					newRow[2] += " R"
				}

				if asento, ok := cData.HandlePositions[barcode[8]]; ok {
					newRow[42] = fmt.Sprintf("#4=%s", asento.Position)
				} else {
					log.Println("Vetimen asentoa ei löytynyt:", string(barcode[8]), barcode, rawRow)
					newRow[2] += " V"
				}
			}
			finish(newRow)
		}

		if !found {
			log.Println("No barcode found for row:", rawRow)
			row[2] += " B"
			finish(row)
		}
	}

	zipWriter.Flush()
	zipWriter.Close()
	return buf, nil
}
