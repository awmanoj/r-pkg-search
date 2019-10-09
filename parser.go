package main 

import(
	"github.com/stapelberg/godebiancontrol"
	"strings"
	"log"
	//"fmt"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// SAMPLE 
// 
//////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
Package: BioInstaller
Version: 0.3.7
Depends: R (>= 3.3.0)
Imports: stringr (>= 1.2.0), futile.logger (>= 1.4.1), configr (>=
        0.3.3), jsonlite, git2r (>= 0.0.3), R.utils (>= 2.5.0), RCurl
        (>= 1.95-4.8), rvest (>= 0.3.2), devtools (>= 1.13.2), stringi
        (>= 1.1.5), shiny, liteq
Suggests: knitr, rmarkdown, testthat, prettydoc, DT
License: MIT + file LICENSE
MD5sum: 51edfde45ebbcc8186ee9f457bbd7ddf
NeedsCompilation: no

Package: biolink
Version: 0.1.6
Imports: rentrez, xml2, DBI, RMySQL, glue, memoise
Suggests: testthat, lintr, httr, covr
License: MIT + file LICENSE
MD5sum: 5474ac2b1785440d0e1b9a1fcd7050bb
NeedsCompilation: no

Package: Biolinv
Version: 0.1-2
Depends: R (>= 3.2.4)
Imports: raster (>= 2.5-2), fields (>= 8.3-6), spatstat (>= 1.48-0), sp
        (>= 1.2-4), grDevices (>= 3.3.2), stats (>= 3.3.2), classInt
        (>= 0.1-23)
License: GPL-3
MD5sum: 0524f36857a07013683d2902f1061f53
NeedsCompilation: no
*/


func parse(data string) ([]Package, error) {
	packages := []Package{}

	paragraphs, err := godebiancontrol.Parse(strings.NewReader(data))
	if err != nil {
		log.Println("err", "problem parsing the data", err)
		return packages, err
	}

	/*
		Name string 
	Version string 
	PublicationDate string 
	Title string
	Description string
	Authors []Person 
	Maintainer []Person

	*/

	// Print a list of which source package uses which package format.
	for _, pkg := range paragraphs {
		//fmt.Printf("%s %s\n", pkg["Package"], pkg["Version"])
		//fmt.Println(pkg["Title"])
		packages = append(packages, Package{
			Name: pkg["Package"],
			Version: pkg["Version"], 
			PublicationDate: pkg["DatePublication"],
			Title: pkg["Title"],
			Description: pkg["Description"],
			Authors: pkg["Author"],
			Maintainer: pkg["Maintainer"],
		})
	}

	return packages, nil 
}