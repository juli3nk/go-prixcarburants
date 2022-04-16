package prixcarburants

import (
	"fmt"
	"strings"
)

type Carburant struct {
	Id string
	Code string
	Label string
}

func Carburants() []*Carburant {
	var carburants []*Carburant

	carburants = append(carburants, &Carburant{Id: "1", Code: "B7", Label: "Gazole"})
	carburants = append(carburants, &Carburant{Id: "2", Code: "E5", Label: "SP95"})
	carburants = append(carburants, &Carburant{Id: "3", Code: "E85", Label: "E85"})
	carburants = append(carburants, &Carburant{Id: "4", Code: "LPG", Label: "GPLc"})
	carburants = append(carburants, &Carburant{Id: "5", Code: "E10", Label: "SP95-E10"})
	carburants = append(carburants, &Carburant{Id: "6", Code: "E5", Label: "SP98"})

	return carburants
}

func CarburantByLabel(label string) *Carburant {
	labelLowercase := strings.ToLower(label)

	for _, c := range Carburants() {
		if strings.ToLower(c.Label) == labelLowercase {
			return c
		}
	}

	return nil
}

func CarburantIcon(id string) string {
	return fmt.Sprintf("%s/images/carburants/%s.svg", URL, id)
}
