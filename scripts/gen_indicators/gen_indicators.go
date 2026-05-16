package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("data skripsi/klasifikasi naive bayes.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows("IM (Indikator miskin)")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("package classifier")
	fmt.Println("\ntype Indicator struct {")
	fmt.Println("\tID          string   `json:\"id\"`")
	fmt.Println("\tLabel       string   `json:\"label\"`")
	fmt.Println("\tOptions     []string `json:\"options\"`")
	fmt.Println("\tSection     string   `json:\"section\"`")
	fmt.Println("}")
	fmt.Println("\nfunc GetIndicators() []Indicator {")
	fmt.Println("\treturn []Indicator{")

	var currentLabel string
	var currentID string
	var options []string

	for i, row := range rows {
		if i == 0 { continue }
		if len(row) == 0 { continue }

		cell0 := strings.TrimSpace(row[0])
		
		if strings.Contains(cell0, "(IM") {
			if currentLabel != "" {
				printIndicator(currentID, currentLabel, options)
			}
			
			labelPart := strings.Split(cell0, "(IM")[0]
			currentLabel = strings.TrimSpace(labelPart)
			
			idPart := strings.Split(strings.Split(cell0, "(IM")[1], ")")[0]
			currentID = "IM" + strings.TrimSpace(idPart)
			options = []string{}
			
			// Option A is often in the same row
			if len(row) > 1 {
				optA := strings.TrimSpace(row[1])
				if strings.Contains(optA, "A ") || strings.Contains(optA, "A\t") {
					options = append(options, cleanOpt(optA))
				}
			}
		} else if len(row) > 1 {
			opt := strings.TrimSpace(row[1])
			if opt != "" {
				options = append(options, cleanOpt(opt))
			}
		}
	}
	if currentLabel != "" {
		printIndicator(currentID, currentLabel, options)
	}

	fmt.Println("\t}")
	fmt.Println("}")
}

func cleanOpt(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 2 && (s[1] == ' ' || s[1] == '\t' || s[1] == '.') {
		return strings.TrimSpace(s[2:])
	}
	// Also handle "A Milik..."
	if strings.HasPrefix(s, "A ") || strings.HasPrefix(s, "B ") || strings.HasPrefix(s, "C ") || strings.HasPrefix(s, "D ") {
		return strings.TrimSpace(s[2:])
	}
	return s
}

func printIndicator(id, label string, options []string) {
	opts := "[]string{\"" + strings.Join(options, "\", \"") + "\"}"
	section := "Lainnya"
	idNum := 0
	fmt.Sscanf(id, "IM%d", &idNum)
	if idNum <= 12 { section = "Kondisi Rumah" } else if idNum <= 24 { section = "Ekonomi Keluarga" } else { section = "Aset & Fasilitas" }
	fmt.Printf("\t\t{ID: \"%s\", Label: \"%s\", Section: \"%s\", Options: %s},\n", id, label, section, opts)
}
