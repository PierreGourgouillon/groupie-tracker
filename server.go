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

type pageA struct {
	Data   API
	Number int
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

func JSON() {
	urlArtists := "https://groupietrackers.herokuapp.com/api/artists"
	urlLocations := "https://groupietrackers.herokuapp.com/api/locations"
	urlDates := "https://groupietrackers.herokuapp.com/api/dates"
	urlRelation := "https://groupietrackers.herokuapp.com/api/relation"

	ParseJSON(urlArtists, &Tracker.Artists)
	ParseJSON(urlLocations, &Tracker.Locations)
	ParseJSON(urlDates, &Tracker.Dates)
	ParseJSON(urlRelation, &Tracker.Relation)
	fmt.Println("Server UP")
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

	for i, b := range Tracker.Artists {

		//si city est check
		if filterAPI.checkCity == "cityIsCheck" {
			isFilterCity = filterCity(filterAPI.city, i)

			if !isFilterCity {
				continue
			}
		}

		//si Members est check
		if filterAPI.checkMembers == "membersIsCheck" {
			isFilterMembers = filterMembers(filterAPI.members, b.Members)

			if !isFilterMembers {
				continue
			}
		}

		//si creation est check
		if filterAPI.checkCreation == "creationIsCheck" {
			isFilterCreation = filterCreation(b.CreationDate, filterAPI.creationDate)
			if !isFilterCreation {
				continue
			}
		}

		//si seulement album est check
		if filterAPI.checkAlbum == "albumIsCheck" {
			isFilterAlbum = filterAlbum(filterAPI.FirstAlbum, b.FirstAlbum)
			if !isFilterAlbum {
				continue
			}
		}

		table = append(table, Tracker.Artists[i])
	}

	if table == nil {
		notFoundArtist = true
	}

	test.Artists = table
	return test, notFoundArtist
}

func filterCity(filterCity string, index int) bool {
	number, _ := strconv.Atoi(filterCity)
	for j, b := range Tracker.Locations.Index {
		if j == index {
			if len(b.Locations) == number {
				return true
			}
		}
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
	if err2 != nil || nbrPath < 0 || nbrPath > 52 {
		errorHandler(w, r)
		return
	}

	Tracker.ID = nbrPath - 1

	tmpl.Execute(w, Tracker)
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
	http.HandleFunc("/", Menu)
	http.HandleFunc("/groupie/", groupiePage)
	http.HandleFunc("/artist/", artistPage)
	http.ListenAndServe(":8080", nil)
}

func Menu(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		errorHandler(w, r)
		return
	}
	tmpl, _ := template.ParseFiles("./templates/index.html")

	tmpl.Execute(w, nil)
}
