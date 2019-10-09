package main 

import(
	"net/http"
	"time"	
	"log"
	"io/ioutil"
	"archive/tar"
	"compress/gzip"
	"io"
	"errors"
	"fmt"
	"strings"
	"os"
	"github.com/blevesearch/bleve"
)


const packagesMetaURL = "https://cran.r-project.org/src/contrib/PACKAGES"
const packageURLPrefix = "https://cran.r-project.org/src/contrib/"

var httpClient = http.Client{
	Timeout: 30 * time.Second, // high timeout due to potentially large file sizes
}

// RunJob executes the job of parsing the PACKAGES file, reading DESCRIPTION for all packages and then
// updating the data model 
func RunJob() {
	log.Println("info", "executing RunJob")

	data, err := fetchPackagesData()
	if err != nil {
		log.Println("err", "fetching the PACKAGES data information", err)
		return 
	}

	packages, err := parse(data)
	if err != nil {
		log.Println("err", "parsing data", err)
		return
	}

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("searchidx.bleve", mapping)
	if err != nil {
		log.Println("err", "failed indexing the bleve db", err)
		
	}
	
	for _, pkg := range packages {
		tarballURL := packageURLPrefix + pkg.Name + "_" + pkg.Version + ".tar.gz"
		
		description, err := readDescriptionFile(tarballURL)
		if err != nil {
			log.Println("err", "read description file, continue with next", err)
			continue
		}

		pkgDetail, err := parse(description)

		fmt.Println(pkgDetail[0].PublicationDate)

		message := struct{
			Id   string
			Name string
			Version string
			PublicationDate string
			Title string
			Description string
			Authors string
			Maintainer string

		}{
			Id: pkgDetail[0].Name + "_" + pkgDetail[0].Version,
			Name: pkgDetail[0].Name, 
			Version: pkgDetail[0].Version, 
			PublicationDate: pkgDetail[0].PublicationDate, 
			Title: pkgDetail[0].Title, 
			Description: pkgDetail[0].Description, 
			Authors: pkgDetail[0].Authors, 
			Maintainer: pkgDetail[0].Maintainer,
		}

		index.Index(message.Id, message)
	}
	os.Exit(1)
}

func fetchPackagesData() (string, error) {
	req, err := http.NewRequest("GET", packagesMetaURL, nil)
	if err != nil {
		log.Println("err", "fetching PACKAGES files from the URL", err)
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("err", "executing the GET request", err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err", "reading the response from PACKAGES URL")
		return "", err
	}

	data := string(body)
	return data, nil 
}


func readDescriptionFile(tarballURL string) (string, error) {
	req, err := http.NewRequest("GET", tarballURL, nil)
	if err != nil {
		log.Println("err", "fetching tarballURL files from the URL", err)
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("err", "executing the GET request for tarball", err)
		return "", err
	}

	defer resp.Body.Close()

	uncompressedStream, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("ExtractTarGz: Next() failed: %s\n", err.Error())
			return "", errors.New("next failed")
		}

		switch header.Typeflag {
		case tar.TypeDir:
			break
		case tar.TypeReg:
			if strings.Contains(header.Name, "/DESCRIPTION") {
				bs, _ := ioutil.ReadAll(tarReader)
				return string(bs), nil
			}
		default:
		}
	}
	return "", errors.New("no DESCRIPTION files")
}

