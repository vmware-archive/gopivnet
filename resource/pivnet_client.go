package resource

import (
	"fmt"
	"net/http"
)

type pivnetClient struct {
	token string
}

func (p *pivnetClient) Do(req *http.Request) (resp *http.Response, err error) {
	p.setPivnetHeaders(req)

	return http.DefaultClient.Do(req)
}

func (p *pivnetClient) DoWithoutRedirect(req *http.Request) (resp *http.Response, err error) {
	p.setPivnetHeaders(req)
	return http.DefaultTransport.RoundTrip(req)
}

func (p *pivnetClient) setPivnetHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Token "+p.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("gopivnet %s", Version))
}
