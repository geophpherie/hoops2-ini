package main

import (
	"fmt"
	"path/filepath"
	"strings"

	html "github.com/geophpherie/hoops2-ini/internal"
	"github.com/geophpherie/hoops2-ini/internal/hoops2"
)

const EXPORT_FOLDER = "./web"
const INI_FOLDER = "./ini"

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
			fmt.Println(file)
			files = append(files, strings.TrimPrefix(file, "web"))
		}
	}

	html.CreateMainSite(files)
}
