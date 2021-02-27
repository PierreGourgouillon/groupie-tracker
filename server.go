package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	Deezer DeezerAPI
}

//DeezerAPI -> Structure les données Deezer de l'artiste
type DeezerAPI struct {
	DeezerArtist ArtistDeezer
	ListSong     ListSong
	AlbumInfo    Album
}

type pageArtist2 struct {
	Data   API
	number int
	Cities []citySearch
}

//ConcertAPI -> Structure pour la page concertLocation
type ConcertAPI struct {
	ID        int
	Locations []string
}

type pageConcert struct {
	Data        API
	SpecialData ConcertAPI
}

//CityAPI -> Structure pour la page cityConcert
type CityAPI struct {
	ID      int
	Artists Artists
}

type pageCity struct {
	Data        API
	SpecialData CityAPI
}

type filter struct {
	FirstAlbum       string
	creationDate     string
	checkMembers     string
	members          string
	checkCity        string
	city             string
	citySearchFilter string
}

type citySearch struct {
	ID   string
	City string
}

var Tracker API
var Artist pageArtist2
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

	structCity := filterCitySearch()

	Artist.Data.Artists = Tracker.Artists
	Artist.Data.Locations = Tracker.Locations
	Artist.Data.Dates = Tracker.Dates
	Artist.Data.Relation = Tracker.Relation
	Artist.Cities = structCity

	filterAPI := filter{
		FirstAlbum:       r.FormValue("firstAlbum"),
		creationDate:     r.FormValue("creationDate"),
		checkMembers:     r.FormValue("checkMembers"),
		members:          r.FormValue("members"),
		checkCity:        r.FormValue("checkCity"),
		city:             r.FormValue("city"),
		citySearchFilter: r.FormValue("citySearchFilter"),
	}

	if filterAPI.FirstAlbum != "" || filterAPI.creationDate != "" || filterAPI.checkMembers != "" || filterAPI.checkCity != "" || filterAPI.citySearchFilter != "" {
		structTest, notFound := Filters(filterAPI)
		if notFound {
			errorHandler(w, r)
			return
		} else {
			tmpl.Execute(w, structTest)
			return
		}
	}

	tmpl.Execute(w, Artist)
}

func Filters(filterAPI filter) (pageArtist2, bool) {
	var isFilterCreation bool
	var isFilterAlbum bool
	var isFilterMembers bool
	var isFilterCity bool
	var isFilterCitySearch bool
	var notFoundArtist bool = false
	var table Artists
	var test pageArtist2

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
		if filterAPI.creationDate != "" {
			isFilterCreation = filterCreation(b.CreationDate, filterAPI.creationDate)
			if !isFilterCreation {
				continue
			}
		}

		//si album est check
		if filterAPI.FirstAlbum != "" {
			isFilterAlbum = filterAlbum(filterAPI.FirstAlbum, b.FirstAlbum)
			if !isFilterAlbum {
				continue
			}
		}

		if filterAPI.citySearchFilter != "" {
			isFilterCitySearch = filterSearchCity(b.ID, filterAPI.citySearchFilter)

			if !isFilterCitySearch {
				continue
			}
		}

		table = append(table, Tracker.Artists[i])
	}

	if table == nil {
		notFoundArtist = true
	}

	test.Data.Artists = table
	return test, notFoundArtist
}

func filterSearchCity(ID int, stringID string) bool {
	a := ""

	for i := range stringID {

		if stringID[i] != ',' {
			a += string(stringID[i])
		} else {
			number, _ := strconv.Atoi(a)
			if number == ID {
				return true
			}
			a = ""
			number = 0
		}
	}
	return false
}

//Nombre de villes
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
	plageAlbum, _ := strconv.Atoi(filterAlbum)
	dateString := albumArtist[6:10]
	date, _ := strconv.Atoi(dateString)

	if date == plageAlbum {
		return true
	}

	return false
}

func filterCreation(creationDate int, filterCreation string) bool {
	plageCreation, _ := strconv.Atoi(filterCreation)

	if creationDate == plageCreation {
		return true
	}
	return false
}

func filterCitySearch() []citySearch {
	var pb []citySearch

	a := false
	for _, b := range Tracker.Locations.Index {
		for _, city := range b.Locations {
			a, pb = search(city, b.ID, pb, a)
		}
	}
	return pb
}

func search(ville string, id int, pb []citySearch, a bool) (bool, []citySearch) {
	if id == 1 && !a {
		item1 := citySearch{
			ID:   strconv.Itoa(id) + ",",
			City: ville,
		}
		pb = append(pb, item1)
		a = true
		return a, pb
	}

	for i, k := range pb {
		if k.City == ville {
			o := pb[i]
			o.ID += strconv.Itoa(id) + ","
			pb[i] = o
			return a, pb
		}

		if i == len(pb)-1 {
			item := citySearch{
				ID:   strconv.Itoa(id) + ",",
				City: ville,
			}
			pb = append(pb, item)
		}
	}

	return a, pb
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
	var selectedArtistPointer *pageArtist
	selectedArtistPointer = &selectedArtist
	Deezer(selectedArtistPointer)

	tmpl.Execute(w, selectedArtistPointer)
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

func menuPage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		errorHandler(w, r)
		return
	}
	tmpl, _ := template.ParseFiles("./templates/index.html")

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

/***************DEEZER API***************/

func Deezer(selectedArtist *pageArtist) {
	nameArtist := SearchNameID(selectedArtist.Number)
	URLTracklist := SearchArtistDeezer(nameArtist, selectedArtist)
	ParseJSON(URLTracklist, &selectedArtist.Deezer.ListSong) //Range dans la structure ListSong, tous les sons
	for i, b := range selectedArtist.Deezer.DeezerArtist.Data {
		if i == 0 {
			fmt.Println(b.Picture)
		}
	}
}

//SearchNameID -> Trouve le nom pour l'id et le modifie
func SearchNameID(ID int) string {
	ID++
	for _, data := range Tracker.Artists {
		if data.ID == ID {
			nameModif := strings.Replace(data.Name, " ", "%20", -1)
			return nameModif
		}
	}
	return ""
}

//SearchArtistDeezer -> Trouve l'artiste dans l'API Deezer
func SearchArtistDeezer(name string, StructPageArtist *pageArtist) string {
	var urlPart1 string = "https://api.deezer.com/search/artist/?q="

	var url = urlPart1 + name

	ParseJSON(url, &StructPageArtist.Deezer.DeezerArtist)
	ParseJSON(url, &StructPageArtist.Deezer.DeezerArtist)

	URLTracklist := ""

	for i, b := range StructPageArtist.Deezer.DeezerArtist.Data {
		if i == 0 {
			URLTracklist = b.TracklistSongURL
		}
	}

	return URLTracklist
}

//ArtistDeezer -> Information sur l'artiste
type ArtistDeezer struct {
	Data []struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Link             string `json:"link"`
		Picture          string `json:"picture_big"`
		NbAlbum          int    `json:"nb_album"`
		NbFan            int    `json:"nb_fan"`
		Radio            bool   `json:"radio"`
		TracklistSongURL string `json:"tracklist"`
	} `json:"data"`
}

//ListSong -> Liste tous les sons de l'artiste
type ListSong struct {
	Data []struct {
		ID             int    `json:"id"`
		Readable       bool   `json:"readable"`
		Title          string `json:"title"`
		TitleShort     string `json:"title_short"`
		LinkURL        string `json:"link"`
		Duration       int    `json:"duration"`
		Rank           int    `json:"rank"`
		ExplicitLyrics bool   `json:"explicit_lyrics"`
		Preview        string `json:"preview"`
		Album          Album  `json:"album"`
	} `json:"data"`
}

//Album -> Information sur l'album
type Album struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	CoverURL       string `json:"cover"`
	TrackListAlbum string `json:"tracklist"`
}
