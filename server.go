package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// Structure pour l'API Globale
type API struct {
	ID        int
	Artists   Artists
	Locations Locations
	Dates     Dates
	Relation  Relation
}

type Artists []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Locations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type Dates struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

// Structure pour la page artist
type pageArtist struct {
	Data   API
	Number int
}

// Structure pour la page concertLocation

type ConcertAPI struct {
	ID        int
	Locations []string
}

type pageConcert struct {
	Data        API
	SpecialData ConcertAPI
}

// Structure pour la page cityConcert

type CityAPI struct {
	ID      int
	Artists Artists
}

type pageCity struct {
	Data        API
	SpecialData CityAPI
}

type filter struct {
	checkAlbum    string
	FirstAlbum    string
	checkCreation string
	creationDate  string
	checkMembers  string
	members       string
	checkCity     string
	city          string
}

var Tracker API
var allLocations []string

func JSON() {
	urlArtists := "https://groupietrackers.herokuapp.com/api/artists"
	urlLocations := "https://groupietrackers.herokuapp.com/api/locations"
	urlDates := "https://groupietrackers.herokuapp.com/api/dates"
	urlRelation := "https://groupietrackers.herokuapp.com/api/relation"

	ParseJSON(urlArtists, &Tracker.Artists)
	ParseJSON(urlLocations, &Tracker.Locations)
	ParseJSON(urlDates, &Tracker.Dates)
	ParseJSON(urlRelation, &Tracker.Relation)

	transformAPILocation()

	fmt.Println("Server UP")
}

func transformAPILocation() {
	for index, api := range Tracker.Locations.Index {
		for i := 0; i < len(api.Locations); i++ {
			locationChange := transformLocation(api.Locations[i])
			Tracker.Locations.Index[index].Locations[i] = locationChange
		}
	}
}

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

func ParseJSON(url string, API interface{}) {
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

func StaticFile() {
	fileserver := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))
}

func groupiePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/groupie/" {
		errorHandler(w, r)
		return
	}
	tmpl, err := template.ParseFiles("./templates/menu.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filterAPI := filter{
		checkAlbum:    r.FormValue("checkAlbum"),
		FirstAlbum:    r.FormValue("firstAlbum"),
		checkCreation: r.FormValue("checkCreation"),
		creationDate:  r.FormValue("creationDate"),
		checkMembers:  r.FormValue("checkMembers"),
		members:       r.FormValue("members"),
		checkCity:     r.FormValue("checkCity"),
		city:          r.FormValue("city"),
	}

	if filterAPI.checkAlbum != "" || filterAPI.checkCreation != "" || filterAPI.checkMembers != "" || filterAPI.checkCity != "" {
		structTest, notFound := Filters(filterAPI)
		if notFound {
			errorHandler(w, r)
			return
		} else {
			tmpl.Execute(w, structTest)
			return
		}
	}

	tmpl.Execute(w, Tracker)
}

func Filters(filterAPI filter) (API, bool) {
	var isFilterCreation bool
	var isFilterAlbum bool
	var isFilterMembers bool
	var isFilterCity bool
	var notFoundArtist bool = false
	var table Artists
	var test API

	//Si city est check
	if filterAPI.checkCity == "cityIsCheck" && filterAPI.checkCreation != "creationIsCheck" && filterAPI.checkAlbum != "albumIsCheck" && filterAPI.checkMembers != "membersIsCheck" {
		for i, b := range Tracker.Locations.Index {
			isFilterCity = filterCity(filterAPI.city, b.Locations)

			if isFilterCity {
				table = append(table, Tracker.Artists[i])
			}
		}
	}

	for i, b := range Tracker.Artists {

		//si creation, album et Members sont check
		if filterAPI.checkCreation == "creationIsCheck" && filterAPI.checkAlbum == "albumIsCheck" && filterAPI.checkMembers == "membersIsCheck" && filterAPI.checkCity != "cityIsCheck" {
			isFilterAlbum = filterAlbum(filterAPI.FirstAlbum, b.FirstAlbum)
			isFilterCreation = filterCreation(b.CreationDate, filterAPI.creationDate)
			isFilterMembers = filterMembers(filterAPI.members, b.Members)

			if isFilterAlbum && isFilterCreation && isFilterMembers {
				table = append(table, Tracker.Artists[i])
			}
		}

		//si creation et album sont check
		if filterAPI.checkCreation == "creationIsCheck" && filterAPI.checkAlbum == "albumIsCheck" && filterAPI.checkMembers != "membersIsCheck" {
			isFilterAlbum = filterAlbum(filterAPI.FirstAlbum, b.FirstAlbum)
			isFilterCreation = filterCreation(b.CreationDate, filterAPI.creationDate)

			if isFilterAlbum && isFilterCreation {
				table = append(table, Tracker.Artists[i])
			}
		}

		//si seulement album est check
		if filterAPI.checkAlbum == "albumIsCheck" && filterAPI.checkCreation != "creationIsCheck" && filterAPI.checkMembers != "membersIsCheck" {
			isFilterAlbum = filterAlbum(filterAPI.FirstAlbum, b.FirstAlbum)

			if isFilterAlbum {
				table = append(table, Tracker.Artists[i])
			}
		}

		//Si seulement creation est check
		if filterAPI.checkCreation == "creationIsCheck" && filterAPI.checkAlbum != "albumIsCheck" && filterAPI.checkMembers != "membersIsCheck" {
			isFilterCreation = filterCreation(b.CreationDate, filterAPI.creationDate)
			if isFilterCreation {
				table = append(table, Tracker.Artists[i])
			}
		}

		//si seulement members est check
		if filterAPI.checkCreation != "creationIsCheck" && filterAPI.checkAlbum != "albumIsCheck" && filterAPI.checkMembers == "membersIsCheck" {
			isFilterMembers = filterMembers(filterAPI.members, b.Members)

			if isFilterMembers {
				table = append(table, Tracker.Artists[i])
			}
		}
	}

	if table == nil {
		notFoundArtist = true
	}

	test.Artists = table
	return test, notFoundArtist
}

func filterCity(filterCity string, Locations []string) bool {
	number, _ := strconv.Atoi(filterCity)
	if len(Locations) == number {
		return true
	}
	return false
}

func filterMembers(filterMembers string, Members []string) bool {
	intFilterMembers, _ := strconv.Atoi(filterMembers)

	if len(Members) == intFilterMembers {
		return true
	}
	return false
}

func filterAlbum(filterAlbum, albumArtist string) bool {
	plageOneAlbum, _ := strconv.Atoi(filterAlbum)
	plageTwoAlbum := plageOneAlbum + 10

	dateString := albumArtist[6:10]
	date, _ := strconv.Atoi(dateString)

	if date >= plageOneAlbum && date <= plageTwoAlbum {
		return true
	}

	return false
}

func filterCreation(creationDate int, filterCreation string) bool {
	plageOneCreation, _ := strconv.Atoi(filterCreation)
	plageTwoCreation := plageOneCreation + 10

	if creationDate >= plageOneCreation && creationDate <= plageTwoCreation {
		return true
	}
	return false
}

func artistPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err1 := template.ParseFiles("./templates/artist.html")
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	nbrPath, err2 := strconv.Atoi(r.URL.Path[8:])
	if err2 != nil || nbrPath < 1 || nbrPath > 52 {
		errorHandler(w, r)
		return
	}

	selectedArtist := pageArtist{
		Data:   Tracker,
		Number: nbrPath - 1,
	}

	tmpl.Execute(w, selectedArtist)
}

func concertLocationPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/concertLocation.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	LocationsTracker := locationsConcertFilter()
	locationsConcert := pageConcert{
		Data:        Tracker,
		SpecialData: LocationsTracker,
	}
	tmpl.Execute(w, locationsConcert)
}

func locationsConcertFilter() ConcertAPI {
	var tableLocations []string
	for index, api := range Tracker.Locations.Index {
		for i := 0; i < len(api.Locations); i++ {
			if len(tableLocations) == 0 {
				tableLocations = append(tableLocations, Tracker.Locations.Index[index].Locations[i])
			} else if locationIn(tableLocations, Tracker.Locations.Index[index].Locations[i]) {
				tableLocations = append(tableLocations, Tracker.Locations.Index[index].Locations[i])
			}
		}
	}
	var API ConcertAPI
	API.Locations = tableLocations
	allLocations = tableLocations
	return API
}

func locationIn(locations []string, selectedLocation string) bool {
	for i := 0; i < len(locations); i++ {
		if locations[i] == selectedLocation {
			return false
		}
	}
	return true
}

func cityConcertPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/cityConcert.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nbrPath, err2 := strconv.Atoi(r.URL.Path[13:])
	if err2 != nil {
		errorHandler(w, r)
		return
	}

	city := allLocations[nbrPath]
	CityTracker := cityConcertFilter(city)

	selectedCity := pageCity{
		Data:        Tracker,
		SpecialData: CityTracker,
	}
	tmpl.Execute(w, selectedCity)
}

func cityConcertFilter(city string) CityAPI {
	var tableArtist Artists
	for index, api := range Tracker.Locations.Index {
		for i := 0; i < len(api.Locations); i++ {
			if Tracker.Locations.Index[index].Locations[i] == city {
				tableArtist = append(tableArtist, Tracker.Artists[index])
				continue
			}
		}
	}
	var API CityAPI
	API.Artists = tableArtist
	return API
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/error.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tmpl.Execute(w, nil)
}

func main() {
	StaticFile()
	JSON()
	http.HandleFunc("/", menuPage)
	http.HandleFunc("/groupie/", groupiePage)
	http.HandleFunc("/artist/", artistPage)
	http.HandleFunc("/concertLocation/", concertLocationPage)
	http.HandleFunc("/cityConcert/", cityConcertPage)
	http.ListenAndServe(":8080", nil)
}

func menuPage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		errorHandler(w, r)
		return
	}
	tmpl, _ := template.ParseFiles("./templates/index.html")

	tmpl.Execute(w, nil)
}
