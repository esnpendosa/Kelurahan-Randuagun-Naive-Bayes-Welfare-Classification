package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("klasifikasi naive bayes tambahan data.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	getNames := func(sheet string) map[string][]string {
		names := make(map[string][]string)
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Println("Error reading sheet:", err)
			return names
		}
		for i, r := range rows {
			if i == 0 { continue }
			if len(r) > 1 && strings.TrimSpace(r[1]) != "" {
				name := strings.TrimSpace(r[1])
				names[name] = r[2:39] // class + 36 indicators
			}
		}
		return names
	}

	t1 := getNames("Data Training 1")
	u1 := getNames("Data Uji 1")
	t2 := getNames("Data Training 2")
	u2 := getNames("Data Uji 2")
	all := getNames("Seluruh Data Warga")

	fmt.Printf("Data Training 1 unique names: %d\n", len(t1))
	fmt.Printf("Data Uji 1 unique names: %d\n", len(u1))
	fmt.Printf("Data Training 2 unique names: %d\n", len(t2))
	fmt.Printf("Data Uji 2 unique names: %d\n", len(u2))
	fmt.Printf("Seluruh Data Warga unique names: %d\n", len(all))

	// Check if all names in T1 + U1 are in all
	inAll1 := 0
	for k := range t1 {
		if _, ok := all[k]; ok { inAll1++ } else { fmt.Printf("T1 name not in All: %q\n", k) }
	}
	for k := range u1 {
		if _, ok := all[k]; ok { inAll1++ } else { fmt.Printf("U1 name not in All: %q\n", k) }
	}
	fmt.Printf("T1 + U1 names in All: %d / 114\n", inAll1)

	// Check if all names in T2 + U2 are in all
	inAll2 := 0
	for k := range t2 {
		if _, ok := all[k]; ok { inAll2++ } else { fmt.Printf("T2 name not in All: %q\n", k) }
	}
	for k := range u2 {
		if _, ok := all[k]; ok { inAll2++ } else { fmt.Printf("U2 name not in All: %q\n", k) }
	}
	fmt.Printf("T2 + U2 names in All: %d / 114\n", inAll2)
}
