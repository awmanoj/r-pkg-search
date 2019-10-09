package main 

type Package struct {
	Name string 
	Version string 
	PublicationDate string 
	Title string
	Description string
	Authors string // this can be further parsed into array (or names, emails) but skipping for now
	Maintainer string // this can be further parsed into array (or names, emails) but skipping for nows
}

