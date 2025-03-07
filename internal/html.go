package html

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/geophpherie/hoops2-ini/internal/hoops2"
)

type NavRowEntryData struct {
	Name string
	Href string
}

func CreateMainSite(files []string) {
	tmpl, err := template.ParseFiles("./internal/templates/index.html", "./internal/templates/nav.html")
	if err != nil {
		panic(err)
	}

	// need to prep for the table
	navRowEntries := make([]NavRowEntryData, len(files))
	for i, file := range files {
		name := strings.TrimSuffix(strings.TrimPrefix(file, "/"), "/index.html")
		navRowEntries[i] = NavRowEntryData{
			Name: name,
			Href: file,
		}
	}
	// need to sort (based on the Name value
	slices.SortFunc(navRowEntries, func(a, b NavRowEntryData) int {
		yearA, _ := strconv.Atoi(a.Name) // Convert string to int
		yearB, _ := strconv.Atoi(b.Name) // Convert string to int
		return yearA - yearB
	})

	// split evently
	top, bot := splitArray(navRowEntries)

	data := struct {
		TopRow []NavRowEntryData
		BotRow []NavRowEntryData
	}{
		TopRow: top,
		BotRow: bot,
	}

	file, _ := os.Create("./web/index.html")
	defer file.Close()

	tmpl.Execute(file, data)
}

func CreatePoolFiles(pool hoops2.Pool) []string {
	tmpl, err := template.ParseFiles("./internal/templates/table.html")
	if err != nil {
		panic(err)
	}

	folder := fmt.Sprintf("./web/%v", pool.Name)
	_ = os.Mkdir(folder, 0755)

	indexFile := filepath.Join(folder, "index.html")
	sortNameFile := filepath.Join(folder, "sortName.html")
	sortRound1File := filepath.Join(folder, "sortRound1.html")
	sortRound2File := filepath.Join(folder, "sortRound2.html")
	sortRound3File := filepath.Join(folder, "sortRound3.html")
	sortRound4File := filepath.Join(folder, "sortRound4.html")
	sortRound5File := filepath.Join(folder, "sortRound5.html")
	sortRound6File := filepath.Join(folder, "sortRound6.html")
	sortRegion1File := filepath.Join(folder, "sortRegion1.html")
	sortRegion2File := filepath.Join(folder, "sortRegion2.html")
	sortRegion3File := filepath.Join(folder, "sortRegion3.html")
	sortRegion4File := filepath.Join(folder, "sortRegion4.html")

	a := []struct {
		active string
		file   string
		sort   hoops2.SortMethod
	}{
		{"pts", indexFile, hoops2.SORT_TOTAL_POINTS},
		{"name", sortNameFile, hoops2.SORT_NAME},
		{"rd1", sortRound1File, hoops2.SORT_ROUND_1},
		{"rd2", sortRound2File, hoops2.SORT_ROUND_2},
		{"rd3", sortRound3File, hoops2.SORT_ROUND_3},
		{"rd4", sortRound4File, hoops2.SORT_ROUND_4},
		{"rd5", sortRound5File, hoops2.SORT_ROUND_5},
		{"rd6", sortRound6File, hoops2.SORT_ROUND_6},
		{"rg1", sortRegion1File, hoops2.SORT_REGION_1},
		{"rg2", sortRegion2File, hoops2.SORT_REGION_2},
		{"rg3", sortRegion3File, hoops2.SORT_REGION_3},
		{"rg4", sortRegion4File, hoops2.SORT_REGION_4},
	}

	for _, b := range a {
		file, _ := os.Create(b.file)

		pool.SortResults(b.sort)

		tmpl.Execute(file, struct {
			hoops2.Pool
			ActiveName string
			PtsHref    string
			NameHref   string
			Rd1Href    string
			Rd2Href    string
			Rd3Href    string
			Rd4Href    string
			Rd5Href    string
			Rd6Href    string
			Rg1Href    string
			Rg2Href    string
			Rg3Href    string
			Rg4Href    string
		}{
			pool,
			b.active,
			strings.TrimPrefix(indexFile, "web"),
			strings.TrimPrefix(sortNameFile, "web"),
			strings.TrimPrefix(sortRound1File, "web"),
			strings.TrimPrefix(sortRound2File, "web"),
			strings.TrimPrefix(sortRound3File, "web"),
			strings.TrimPrefix(sortRound4File, "web"),
			strings.TrimPrefix(sortRound5File, "web"),
			strings.TrimPrefix(sortRound6File, "web"),
			strings.TrimPrefix(sortRegion1File, "web"),
			strings.TrimPrefix(sortRegion2File, "web"),
			strings.TrimPrefix(sortRegion3File, "web"),
			strings.TrimPrefix(sortRegion4File, "web"),
		})

		file.Close()
	}

	// we need an array of htmlFile and sort value
	// template will need to take pool data as well as _all_ the links

	return []string{indexFile}
}

func splitArray[T any](arr []T) ([]T, []T) {
	// Calculate the index to split the array at
	mid := (len(arr) + 1) / 2 // This ensures the extra element goes to the first part if the length is odd

	// Slice the array into two parts
	firstPart := arr[:mid]
	secondPart := arr[mid:]

	return firstPart, secondPart
}
