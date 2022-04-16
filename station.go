package prixcarburants

import (
	"net/url"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/juli3nk/stack/client"
)

const URL = "https://www.prix-carburants.gouv.fr"

type Station struct {
	Name string
	Address string
	Price string
	Date string
}

func Stations(carbu, departement, localisation string) ([]Station, error) {
	var stations []Station

	if len(departement) > 0 && len(localisation) > 0 {
		departement = ""
	}

	u, err := client.ParseUrl(URL)
	if err != nil {
		return nil, err
	}

	// Get form token
	ccR1 := &client.Config{
		Scheme: u.Scheme,
		Host:   u.Host,
		Port:   u.Port,
		Path:   "/",
		Cookie: client.Cookie{Enabled: true},
	}

	reqR1, err := client.New(ccR1)
	if err != nil {
		return nil, err
	}

	resultR1 := reqR1.Get()

	docIndex := soup.HTMLParse(string(resultR1.Body))
	token := docIndex.Find("input", "name", "rechercher[_token]").Attrs()["value"]

	// Post form
	ccR2 := &client.Config{
		Scheme: u.Scheme,
		Host:   u.Host,
		Port:   u.Port,
		Path:   "/",
		Cookie: client.Cookie{Enabled: true, Jar: resultR1.CookieJar},
	}

	reqR2, err := client.New(ccR2)
	if err != nil {
		return nil, err
	}

	reqR2.HeaderAdd("Content-Type", "application/x-www-form-urlencoded")

	v := url.Values{}
	v.Add("rechercher[choix_carbu][]", carbu)
	v.Add("rechercher[geolocalisation_long]", "")
	v.Add("rechercher[geolocalisation_lat]", "")
	v.Add("rechercher[departement]", departement)
	v.Add("rechercher[localisation]", localisation)
	v.Add("rechercher[type_enseigne]", "")
	v.Add("rechercher[_token]", token)

	data := v.Encode()

	resultR2 := reqR2.Post(strings.NewReader(data))
	if resultR2.Error != nil {
		return nil, resultR2.Error
	}

	// GET /recherche/
	ccR3 := &client.Config{
		Scheme: u.Scheme,
		Host:   u.Host,
		Port:   u.Port,
		Path:   "/recherche/",
		Cookie: client.Cookie{Enabled: true, Jar: resultR2.CookieJar},
	}

	reqR3, err := client.New(ccR3)
	if err != nil {
		return nil, err
	}

	reqR3.ValueAdd("page", "1")
	reqR3.ValueAdd("limit", "100")

	resultR3 := reqR3.Get()

	doc := soup.HTMLParse(string(resultR3.Body))

	//number := doc.Find("div", "id", "sectionNombreResultats").Find("h3").Text()

	rows := doc.Find("table", "id", "tab_resultat").FindAll("tr", "class", "data")

	for _, r := range rows {
		station := Station{}

		desc := r.Find("div", "class", "pdv-description")

		// Name
		station.Name = desc.Find("h4", "class", "title").Find("strong").Text()

		// Address
		var address []string

		addressTmp := desc.FindAll("span")
		for _, a := range addressTmp {
			address = append(address, a.Text())
		}
		station.Address = strings.Join(address, " ")

		// Price - Date
		var date []string

		chiffres := r.Find("td", "class", "chiffres").FindAll("span")
		for i, c := range chiffres {
			if i == 0 {
				station.Price = c.Find("strong").Text()

				continue
			}

			date = append(date, c.Text())
		}
		station.Date = strings.Join(date, " ")

		stations = append(stations, station)
	}

	return stations, nil
}
