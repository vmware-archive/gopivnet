package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cfmobile/gopivnet/resource"
)

var product = flag.String("product", "", "product to download")

var version = flag.String("version", "", "version of the product")

func main() {
	flag.Parse()

	if *product == "" {
		log.Fatal("Need a product name")
	}
	pr := resource.NewRequester("https://network.pivotal.io", "yonFmSKVP4cFcqW5Khg6")

	prod, err := pr.GetProduct(*product)
	if err != nil {
		log.Fatal(err)
	}

	productFiles, _ := pr.GetProductFiles(prod.Releases[0])

	fmt.Printf("%v\n", productFiles)
	var pivotalProduct *resource.ProductFile
	for index, productFile := range productFiles.Files {
		fmt.Println(productFile.AwsObjectKey)
		if strings.Contains(productFile.AwsObjectKey, ".pivotal") {
			pivotalProduct = &productFiles.Files[index]
			break
		}
	}

	if pivotalProduct == nil {
		log.Fatal("Unable to find a pivotal product")
	}

	url, _ := pr.GetProductDownloadUrl(pivotalProduct)

	fmt.Println(url)
	download(url)
}

func download(url string) {
	outputFile := *product + ".pivotal"

	out, err := os.Create(outputFile)
	defer out.Close()

	resp, err := http.Get(url)
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Unable to write to file")
		return
	}

	fmt.Printf("Written %d bytes to file", n)
}
