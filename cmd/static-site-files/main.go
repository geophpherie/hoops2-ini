package main

import (
	"path/filepath"
	"strings"

	html "github.com/geophpherie/hoops2-ini/internal"
	"github.com/geophpherie/hoops2-ini/internal/hoops2"
)

const EXPORT_FOLDER = "./web"
const INI_FOLDER = "./ini"
const RESULTS_FILES_FOLDER = "./results_files"

func main() {
	matches, err := filepath.Glob(filepath.Join(INI_FOLDER, "*.I*"))
	if err != nil {
		panic(err)
	}

	var files []string
	for _, match := range matches {
		pool, err := hoops2.NewPoolFromFile(match)
		if err != nil {
			panic(err)
		}

		htmlFiles := html.CreatePoolFiles(pool)

		for _, file := range htmlFiles {
			files = append(files, strings.TrimPrefix(file, "web"))
		}
	}

	html.CreateMainSite(files)
}
