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

type Filter struct {
	creation string
	album    string
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

	filter := Filter{
		creation: r.FormValue("creation"),
		album:    r.FormValue("album"),
	}

	if filter.album != "" { // A coché Albu

	} else if filter.creation != "" { // A coché Creation
		fmt.Println("salut")
	}

	tmpl.Execute(w, Tracker)
}

func artistPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err1 := template.ParseFiles("./templates/artist.html")
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	nbrPath, err2 := strconv.Atoi(r.URL.Path[8:])
	if err2 != nil || nbrPath < 0 || nbrPath > 51 {
		errorHandler(w, r)
		return
	}
	Tracker.ID = nbrPath

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
