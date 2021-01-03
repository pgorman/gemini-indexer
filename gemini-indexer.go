// Copyright 2020 Paul Gorman. Licensed under the GPL.

// Gemini Indexer generates an index.gmi file.
//
// See the Project Gemini documentation and spec at:
// https://gemini.circumlunar.space/docs/
// gemini://gemini.circumlunar.space/docs/

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var reDate *regexp.Regexp

type link struct {
	File  string
	Date  string
	Label string
}

type templateData struct {
	Title       string
	GeminiLinks []link
	OtherFiles  []string
}

// extractLabel returns the file name without a date or extension.
func extractLabel(file string) string {
	var s string
	s = strings.TrimSuffix(file, ".gmi")
	s = strings.TrimSuffix(s, ".gemini")
	s = strings.ReplaceAll(s, "_", " ")
	s = reDate.ReplaceAllString(s, "")
	return s
}

func main() {
	var err error
	var optDotFiles bool
	var optIgnore string
	var optInputDir string
	var optOutputFile string
	var optTemplate string
	var optTitle string
	var tmpl *template.Template

	flag.BoolVar(&optDotFiles, "dotfiles", false, "include dotfiles, like '..' and '.git' in the index")
	flag.StringVar(&optIgnore, "ignore", "", "comma-separated list of files to not include in the index")
	flag.StringVar(&optInputDir, "indir", "", "path to the directory of files to index (default: current directory")
	flag.StringVar(&optOutputFile, "outfile", "", "where to write the index file (default: stdout)")
	flag.StringVar(&optTemplate, "template", "", "template file for the index (see https://pkg.go.dev/text/template)")
	flag.StringVar(&optTitle, "title", "Gemini Index", `title displayed on index page, like "Jane Smith's Gemlog"`)
	flag.Parse()

	if optInputDir == "" {
		optInputDir, err = os.Getwd()
		if err != nil {
			log.Fatal("main: unable to get the current working directory (try setting --indir):", err)
		}
	}

	if optTemplate == "" {
		tmpl, err = template.New("index").Parse(`# {{.Title}}
{{if .GeminiLinks}}{{range $l := .GeminiLinks}}
=> {{$l.File}} {{$l.Date}} {{$l.Label}}{{end}}
{{end}}{{if .OtherFiles}}{{range $f := .OtherFiles}}
=> {{$f}}{{end}}{{end}}
`)
	} else {
		tmpl, err = template.ParseFiles(optTemplate)
	}
	if err != nil {
		log.Fatal("main: failed to parse template:", err)
	}

	reDate = regexp.MustCompile(`(?:[-_])?(\d{4}-\d{2}-\d{2})(?:[-_])?`)
	var td templateData
	td.Title = optTitle
	td.GeminiLinks = make([]link, 0, 50)
	td.OtherFiles = make([]string, 0, 50)

	files, err := ioutil.ReadDir(optInputDir)
	if err != nil {
		log.Fatalf("main: failed to open input directory '%s': %v", optInputDir, err)
	}
	for _, f := range files { // Non-Gemini files, sorted alphabetically (well, by increasing value)
		if !optDotFiles && f.Name()[0] == "."[0] {
			continue
		}
		s := strings.Split(f.Name(), ".")
		if (len(s) == 1) || (len(s) > 1 && (s[len(s)-1] != "gmi" && s[len(s)-1] != "gemini")) {
			td.OtherFiles = append(td.OtherFiles, f.Name())
		}
	}
	for i := len(files) - 1; i >= 0; i-- { // Gemini files, sorted newest to oldest
		f := files[i]
		if !optDotFiles && f.Name()[0] == "."[0] {
			continue
		}
		s := strings.Split(f.Name(), ".")
		if len(s) > 1 && (s[len(s)-1] == "gmi" || s[len(s)-1] == "gemini") {
			var lk link
			lk.File = f.Name()
			d := reDate.FindStringSubmatch(f.Name())
			if len(d) > 1 {
				lk.Date = d[1]
			} else {
				lk.Date = f.ModTime().Format("2006-01-02")
			}
			l := extractLabel(f.Name())
			if l != "" {
				lk.Label = l
			} else {
				lk.Label = f.Name()
			}
			td.GeminiLinks = append(td.GeminiLinks, lk)
		}
	}

	if optOutputFile != "" {
		of, err := os.OpenFile(optOutputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("main: failed to open output file to write:", err)
		}
		defer of.Close()
		err = tmpl.Execute(of, td)
	} else {
		err = tmpl.Execute(os.Stdout, td)
	}
	if err != nil {
		log.Fatal("main: failed to default template:", err)
	}
}
