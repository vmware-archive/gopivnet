package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var RequireEula = 451

type ReleaseRequester interface {
	GetProduct(productName string) (*Product, error)
	GetProductFiles(release Release) (*ProductFiles, error)
	GetProductDownloadUrl(productFile *ProductFile) (string, error)
}

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
	DoWithoutRedirect(req *http.Request) (resp *http.Response, err error)
}

func NewRequester(url string, token string) ReleaseRequester {
	return &PivnetRequester{
		pivnetUrl: url,
		client:    &pivnetClient{token: token},
	}
}

type PivnetRequester struct {
	pivnetUrl string
	client    HttpClient
}

func (p *PivnetRequester) getProductRequest(productName string) *http.Request {
	requestUrl := fmt.Sprintf("%s/api/v2/products/%s/releases", p.pivnetUrl, productName)

	req, _ := http.NewRequest("GET", requestUrl, nil)
	return req
}

func (p *PivnetRequester) GetProduct(productName string) (*Product, error) {
	req := p.getProductRequest(productName)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("bad status code from server")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	product := Product{}

	err = json.Unmarshal(body, &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *PivnetRequester) GetProductFiles(release Release) (*ProductFiles, error) {
	productFilesLink, ok := release.Links["product_files"]
	if !ok {
		return nil, errors.New("Unable to get product files")
	}

	req, _ := http.NewRequest("GET", productFilesLink.Url, nil)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("bad status code from server")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	productFiles := ProductFiles{}

	err = json.Unmarshal(body, &productFiles)
	if err != nil {
		return nil, err
	}

	return &productFiles, nil
}

func (p *PivnetRequester) GetProductDownloadUrl(productFile *ProductFile) (string, error) {
	downloadLink, ok := productFile.Links["download"]
	if !ok {
		return "", errors.New("Unable to get product files")
	}

	req, _ := http.NewRequest("POST", downloadLink.Url, nil)
	resp, err := p.client.DoWithoutRedirect(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 451 {
		body, _ := ioutil.ReadAll(resp.Body)
		eula := EulaMessage{}
		json.Unmarshal(body, &eula)

		err = p.acceptEula(eula.Links["eula_agreement"].Url)
		if err != nil {
			return "", err
		}
		resp, err = p.client.DoWithoutRedirect(req)
		if err != nil {
			return "", err
		}
	}

	if resp.StatusCode != http.StatusFound {
		return "", errors.New("bad status code from server")
	}

	downloadUrl := resp.Header.Get("Location")
	return downloadUrl, nil
}

func (p *PivnetRequester) acceptEula(url string) error {
	req, _ := http.NewRequest("POST", url, nil)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Unable to accept eula")
	}
	return nil
}
