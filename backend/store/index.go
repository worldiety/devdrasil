package store

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

const storeIndexUrl = "http://store.devdrasil.worldiety.com/index"

type Index struct {
	Plugins []*Plugin `json:"plugins"`
	Vendors []*Vendor `json:"vendors"`
	Version int       `json:"version"`
}

type Plugin struct {
	Vendor           string   `json:"vendor"`
	Name             string   `json:"name"`
	Id               string   `json:"id"`
	Sales            []*Sale  `json:"sales"`
	Source           *Source  `json:"source"`
	Images           []string `json:"images"`
	Videos           []string `json:"videos"`
	DescriptionUrl   string   `json:"descriptionUrl"`
	PrivacyPolicyUrl string   `json:"privacyPolicyUrl"`
	TermsUrl         string   `json:"termsUrl"`
	Keywords         []string `json:"keywords"`
	Categories       []string `json:"categories"`
	IconUrl          string   `json:"iconUrl"`
	Rating           *Rating  `json:"rating"`
}

type Vendor struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	State     string `json:"state"`
	City      string `json:"city"`
	Zip       string `json:"zip"`
	TaxNumber string `json:"taxNumber"`
	Court     string `json:"court"`
	Manager   string `json:"manager"`
	Phone     string `json:"phone"`
	Web       string `json:"web"`
	EMail     string `json:"email"`
}

type Sale struct {
	Type        string `json:"Type"`
	Price       int    `json:"Price"`
	Currency    string `json:"currency"`
	TaxIncluded bool   `json:"taxIncluded"`
	Href        string `json:"href"`
}

type Source struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Rating struct {
	Stars5   int        `json:"s5"`
	Stars4   int        `json:"s4"`
	Stars3   int        `json:"s3"`
	Stars2   int        `json:"s2"`
	Stars1   int        `json:"s1"`
	Comments []*Comment `json:"comments"`
}

type Comment struct {
	From  string `json:"from"`
	Stars int    `json:"stars"`
	Date  int64  `json:"date"`
	Text  string `json:"text"`
}

type Store struct {
}

func (r *Store) GetIndex() (*Index, error) {
	resp, err := http.Get(storeIndexUrl)
	if err != nil {
		return nil, err
	}

	tmp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	idx := &Index{}
	err = json.Unmarshal(tmp, idx)
	if err != nil {
		return nil, err
	}

	return idx, nil
}

func (i *Index) GetPlugin(id string) *Plugin {
	for _, p := range i.Plugins {
		if p.Id == id {
			return p
		}
	}
	return nil
}
