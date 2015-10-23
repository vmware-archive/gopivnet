package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cfmobile/gopivnet/resource"
)

type Api interface {
	GetLatestProductFile(productName string) (*resource.ProductFile, error)
	GetProductFileForVersion(productName, version string) (*resource.ProductFile, error)
	Download(productFile *resource.ProductFile, fileName string) error
}

type PivnetApi struct {
	Requester resource.ReleaseRequester
}

func New(token string) Api {
	return &PivnetApi{
		Requester: resource.NewRequester("https://network.pivotal.io", token),
	}
}

func (p *PivnetApi) GetLatestProductFile(productName string) (*resource.ProductFile, error) {
	if productName == "" {
		return nil, errors.New("Must specify a product name")
	}

	prod, err := p.Requester.GetProduct(productName)
	if err != nil {
		return nil, err
	}

	productFiles, err := p.Requester.GetProductFiles(prod.Releases[0])
	if err != nil {
		return nil, err
	}

	pivotalProduct := getPivotalProduct(productFiles)
	if pivotalProduct == nil {
		return nil, errors.New("Unable to fund a pivotal product")
	}

	return pivotalProduct, nil
}

func getPivotalProduct(productFiles *resource.ProductFiles) *resource.ProductFile {
	for index, productFile := range productFiles.Files {
		if strings.Contains(productFile.AwsObjectKey, ".pivotal") {
			return &productFiles.Files[index]
		}
	}

	return nil
}

func (p *PivnetApi) GetProductFileForVersion(productName, version string) (*resource.ProductFile, error) {
	if productName == "" {
		return nil, errors.New("Must specify a product name")
	}

	if version == "" {
		return nil, errors.New("Must specify a product version")
	}

	prod, err := p.Requester.GetProduct(productName)
	if err != nil {
		return nil, err
	}

	matchingRelease := getReleaseForVersion(prod, version)
	if matchingRelease == nil {
		return nil, errors.New("Specified version not found")
	}

	productFiles, err := p.Requester.GetProductFiles(*matchingRelease)
	if err != nil {
		return nil, err
	}

	pivotalProduct := getPivotalProduct(productFiles)
	if pivotalProduct == nil {
		return nil, errors.New("Unable to fund a pivotal product")
	}

	return pivotalProduct, nil
}

func getReleaseForVersion(product *resource.Product, version string) *resource.Release {
	for index, release := range product.Releases {
		if release.Version == version {
			return &product.Releases[index]
		}
	}

	return nil
}

func (p *PivnetApi) Download(productFile *resource.ProductFile, fileName string) error {
	if productFile == nil {
		return errors.New("Nil product passed in")
	}

	url, err := p.Requester.GetProductDownloadUrl(productFile)
	if err != nil {
		return err
	}

	download(url, fileName)

	return nil
}

func download(url, fileName string) {
	out, err := os.Create(fileName)
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
