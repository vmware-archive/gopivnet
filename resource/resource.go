package resource

import "strings"

type Product struct {
	Releases []Release `json:"releases"`
}

type Release struct {
	Id              int    `json:"id"`
	Version         string `json:"version"`
	ReleaseType     string `json:"release_type"`
	ReleaseDate     string `json:"release_date"`
	ReleaseNotesUrl string `json:"release_notes_url"`
	Availability    string `json:"availability"`
	Description     string `json:"description"`
	Eula            Eula   `json:"eula"`
	Links           Links  `json:"_links"`
}

type Eula struct {
	Id   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Links map[string]Link

type Link struct {
	Url string `json:"href"`
}

type ProductFiles struct {
	Files []ProductFile `json:"product_files"`
}

type ProductFile struct {
	Id           int    `json:"id"`
	AwsObjectKey string `json:"aws_object_key"`
	FileVersion  string `json:"file_version"`
	Links        Links  `json:"_links"`
}

func (p *ProductFile) Name() string {
	tokens := strings.Split(p.AwsObjectKey, "/")
	return tokens[len(tokens)-1]
}

type EulaMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Links   Links  `json:"_links"`
}
