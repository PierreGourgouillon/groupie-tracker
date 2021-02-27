package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Structure pour l'API Globale

// API -> structure contenants les structure de l'API groupie-trackers
type API struct {
	ID        int
	Artists   Artists
	Locations Locations
	Dates     Dates
	Relation  Relation
}

// Artists -> structures contenant les informations sur les artistes
type Artists []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// Locations -> structure contenant les villes des concerts
type Locations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

// Dates -> structure contenant les dates des concerts
type Dates struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

// Relation -> structures reliant les dates des concerts à leurs villes
type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

// Structure pour la page artist
// pageArtists -> structures pour envoyer les informations à la page /artist
type pageArtist struct {
	Data        API
	SpecialData ArtistAPI
	Flag        []string
}

// ArtistAPI -> structure Artiste spécial pour la page /artist contient un seul artist, ses lieux de concert et leurs dates
type ArtistAPI struct {
	ID        int
	Artists   Artists
	Locations []string
	Relation  [][]string
}

type pageArtist2 struct {
	Data   API
	number int
	Cities []citySearch
}

// Structure pour la page concertLocation
// pageConcert -> structure pour envoyer les informations à la page /concertLocations
type pageConcert struct {
	Data        API
	SpecialData ConcertAPI
	Flag        []string
}

// ConcertAPI -> structure spécial pour la page /concertLocations contient les lieux de concert
type ConcertAPI struct {
	ID        int
	Locations []string
}

// Structure pour la page cityConcert
// pageCity -> structure pour envoyer les informations à la page /cityConcert
type pageCity struct {
	Data        API
	SpecialData CityAPI
}

// CityAPI -> structure spécial pour la page /cityConcert des artistes
type CityAPI struct {
	ID      int
	Artists Artists
}

// filter -> structure pour recevoir les informations du filtre
type filter struct {
	FirstAlbum       string
	creationDate     string
	checkMembers     string
	members          string
	checkCity        string
	city             string
	citySearchFilter string
}

// citySearch ->
type citySearch struct {
	ID   string
	City string
}

// Variable Globale

// Tracker -> Variable contenant toutes les informatiions de l'API groupie-tracker
var Tracker API

// Artist ->
var Artist pageArtist2

// allLocations -> Variable contenant tous les lieux de concert en un seul exemplaire
var allLocations []string

// flagCountry -> Variable contenant tous les pays et leur code pour pouvoir utiliser l'API flagcdn
var flagCountry map[string]string

// serverJSON -> Stocke les urls des APIs et prévient quand tous les informations ont été parsées
func serverJSON() {
	urlArtists := "https://groupietrackers.herokuapp.com/api/artists"
	urlLocations := "https://groupietrackers.herokuapp.com/api/locations"
	urlDates := "https://groupietrackers.herokuapp.com/api/dates"
	urlRelation := "https://groupietrackers.herokuapp.com/api/relation"
	urlDrapeau := "https://flagcdn.com/en/codes.json"

	unmarshallJSON(urlArtists, &Tracker.Artists)
	unmarshallJSON(urlLocations, &Tracker.Locations)
	unmarshallJSON(urlDates, &Tracker.Dates)
	unmarshallJSON(urlRelation, &Tracker.Relation)
	unmarshallJSON(urlDrapeau, &flagCountry)

	fmt.Println("Before", Tracker.Relation.Index)
	transformAPIRelation()
	fmt.Println("After", Tracker.Relation.Index)

	fmt.Println("Server UP")
}

// transformAPILocation -> change l'orthographe des lieux de concert
func transformAPILocation() {
	for index, api := range Tracker.Locations.Index {
		for i := 0; i < len(api.Locations); i++ {
			locationChange := transformLocation(api.Locations[i])
			Tracker.Locations.Index[index].Locations[i] = locationChange
		}
	}
}

// deleteCountry -> supprime le pays des lieux de concert
func deleteCountry(table []string) []string {
	for i, location := range table {
		idxPipe := strings.Index(location, "|")
		table[i] = location[:idxPipe-1]
	}
	return table
}

// transformAPIRelation -> change l'orthographe des lieux de concert
func transformAPIRelation() {
	for i := 0; i < len(Tracker.Artists); i++ {
		var mapChange = make(map[string][]string)
		for index, api := range Tracker.Relation.Index[i].DatesLocations {
			locationChange := transformLocation(index)
			mapChange[locationChange] = api
		}
		Tracker.Relation.Index[i].DatesLocations = mapChange
	}
}

// deleteCountryMap -> change l'orthographe des lieux de concert
func deleteCountryMap(mapRelation map[string][]string) map[string][]string {
	var mapChange = make(map[string][]string)
	for index, api := range mapRelation {
		locationChange := deleteCountry([]string{index})
		mapChange[locationChange[0]] = api
	}
	return mapChange
}

// transformLocation -> remplace les "_" par un " " et les "-" par un "|"
func transformLocation(text string) string {
	var newText string = ""
	for i := 0; i < len(text); i++ {
		if text[i] == '_' {
			newText = newText + " "
		} else if text[i] == '-' {
			newText = newText + " | "
		} else {
			newText = newText + string(text[i])
		}
	}
	return newText
}

// unmarshallJSON -> parse les données JSON
func unmarshallJSON(url string, API interface{}) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	json.Unmarshal(body, &API)
}

func main() {
	serverJSON()
	salut := deleteCountryMap(Tracker.Relation.Index[0].DatesLocations)
	fmt.Println(salut)
}
