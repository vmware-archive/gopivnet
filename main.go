package main

import (
	"flag"
	"log"

	"github.com/cfmobile/gopivnet/api"
	"github.com/cfmobile/gopivnet/resource"
)

var productName = flag.String("product", "", "product to download")

var version = flag.String("version", "", "version of the product. If missing download the latest version")

var token = flag.String("token", "", "pivnet token")

var file = flag.String("file", "", "filename where to save the pivotal product")

var fileType = flag.String("fileType", "", "type of file.  Defaults to 'pivotal' tile.")

func main() {
	flag.Parse()

	if *productName == "" {
		log.Fatal("Need a product name")
	}

	if *token == "" {
		log.Fatal("Need a pivnet token")
	}

    if *fileType == "" {
        *fileType = "pivotal"
    }

	pivnetApi := api.New(*token)

	var pivotalProduct *resource.ProductFile
	var err error
	if *version != "" {
		pivotalProduct, err = pivnetApi.GetProductFileForVersion(*productName, *version, *fileType)
	} else {
		pivotalProduct, err = pivnetApi.GetLatestProductFile(*productName, *fileType)
	}
	if err != nil {
		log.Fatal(err)
	}

	fileName := *file
	if fileName == "" {
		fileName = pivotalProduct.Name()
	}

	pivnetApi.Download(pivotalProduct, fileName)
}
