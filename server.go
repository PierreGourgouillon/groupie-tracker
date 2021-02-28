package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
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
	Deezer      DeezerAPI
}

//DeezerAPI -> Structure les données Deezer de l'artiste
type DeezerAPI struct {
	DeezerArtist ArtistDeezer
	ListSong     ListSong
	AlbumInfo    Album
}

// ArtistAPI -> structure Artiste spécial pour la page /artist contient un seul artist, ses lieux de concert et leurs dates
type ArtistAPI struct {
	ID        int
	Artists   Artists
	Locations []string
	Relation  map[string][]string
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
	City        string
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

	transformAPILocation()
	transformAPIRelation()

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

// staticFile -> créer le serveur de fichiers destinées au site (css, js, images)
func staticFile() {
	fileserver := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))
}

// groupiePage -> charge la page /groupie
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
		structTest, notFound := filters(filterAPI)
		if notFound {
			errorHandler(w, r)
			return
		}

		tmpl.Execute(w, structTest)
		return
	}

	tmpl.Execute(w, Artist)
}

// filters -> filtre en fonction des données reçues
func filters(filterAPI filter) (pageArtist2, bool) {
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

// filterMembers -> filtre les artistes en fonction du nombre de membre
func filterMembers(filterMembers string, Members []string) bool {
	intFilterMembers, _ := strconv.Atoi(filterMembers)

	if len(Members) == intFilterMembers {
		return true
	}
	return false
}

// filterAlbum -> filtre les artistes en fonction de la date de sortie du premier album
func filterAlbum(filterAlbum, albumArtist string) bool {
	plageAlbum, _ := strconv.Atoi(filterAlbum)
	dateString := albumArtist[6:10]
	date, _ := strconv.Atoi(dateString)

	if date == plageAlbum {
		return true
	}

	return false
}

// filterCreation -> filtre les artistes en fonction de la date de création
func filterCreation(creationDate int, filterCreation string) bool {
	plageCreation, _ := strconv.Atoi(filterCreation)

	if creationDate == plageCreation {
		return true
	}
	return false
}

// filterCitySearch ->
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

// search ->
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

// artistPage -> charge la page /artist
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

	var ArtistTracker ArtistAPI

	var tableArtist Artists
	tableArtist = append(tableArtist, Tracker.Artists[nbrPath-1])
	ArtistTracker.Artists = tableArtist

	var tableLocations []string
	for _, location := range Tracker.Locations.Index[nbrPath-1].Locations {
		tableLocations = append(tableLocations, location)
	}
	ArtistTracker.Locations = tableLocations

	var tableRelation map[string][]string
	tableRelation = Tracker.Relation.Index[nbrPath-1].DatesLocations
	ArtistTracker.Relation = tableRelation

	var TableFlag []string
	TableFlag = flagCountryFilter(ArtistTracker.Locations)

	ArtistTracker.Locations = deleteCountry(ArtistTracker.Locations)

	ArtistTracker.Relation = deleteCountryMap(ArtistTracker.Relation)

	selectedArtist := pageArtist{
		Data:        Tracker,
		SpecialData: ArtistTracker,
		Flag:        TableFlag,
	}
	var selectedArtistPointer *pageArtist
	selectedArtistPointer = &selectedArtist
	Deezer(selectedArtistPointer)

	tmpl.Execute(w, selectedArtistPointer)
}

// deleteCountry -> supprime le pays des lieux de concert
func deleteCountry(table []string) []string {
	for i, location := range table {
		idxPipe := strings.Index(location, "|")
		table[i] = location[:idxPipe-1]
	}
	return table
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

// concertLocationPage -> charge la page /concertLoaction
func concertLocationPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/concertLocation.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var LocationsTracker ConcertAPI
	var TableFlag []string

	LocationsTracker = locationsConcertFilter()
	TableFlag = flagCountryFilter(allLocations)
	LocationsTracker.Locations = deleteCountry(LocationsTracker.Locations)

	locationsConcert := pageConcert{
		Data:        Tracker,
		SpecialData: LocationsTracker,
		Flag:        TableFlag,
	}
	tmpl.Execute(w, locationsConcert)
}

// flagCountryFilter -> créer un tableau contenant les codes des pays (de l'API flagcdn) dans l'ordre des lieux de concert
func flagCountryFilter(locations []string) []string {
	var TableFlag []string
	for _, location := range locations {
		idxPipe := strings.Index(location, "|")
		if idxPipe == -1 {
			continue
		}
		location = location[idxPipe+2:]
		if location == "usa" {
			TableFlag = append(TableFlag, "us")
			continue
		} else if location == "us" {
			TableFlag = append(TableFlag, "us")
			continue
		} else if location == "uk" {
			TableFlag = append(TableFlag, "gb")
			continue
		} else if location == "netherlands antilles" {
			TableFlag = append(TableFlag, "nl")
			continue
		} else if location == "czech republic" {
			TableFlag = append(TableFlag, "cz")
			continue
		} else if location == "brasil" {
			TableFlag = append(TableFlag, "br")
			continue
		} else if location == "philippine" {
			TableFlag = append(TableFlag, "ph")
			continue
		} else if location == "korea" {
			TableFlag = append(TableFlag, "kr")
			continue
		}
		for countryCode, value := range flagCountry {
			value = strings.ToLower(value)
			if value == location {
				TableFlag = append(TableFlag, countryCode)
				break
			}
		}
	}
	return TableFlag
}

// locationsConcertFilter -> créer un tableau contenant les lieux des concert en un seul exemplaire, l'insert dans la variable globale allLocations est dans la structure API destinée à la page /concertLocations
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
	sort.Strings(tableLocations)
	allLocations = tableLocations
	API.Locations = tableLocations
	return API
}

// locationIn -> regarde si un lieu de concert est déjà présent dans le tableau
func locationIn(locations []string, selectedLocation string) bool {
	for i := 0; i < len(locations); i++ {
		if locations[i] == selectedLocation {
			return false
		}
	}
	return true
}

// cityConcertPage -> charge la page /cityConcert
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
		City:        city,
	}
	tmpl.Execute(w, selectedCity)
}

// cityConcertFilter -> filtre les artists en fonction de leurs lieux de concert
func cityConcertFilter(city string) CityAPI {
	var tableArtist Artists
	for index, api := range Tracker.Locations.Index {
		for i := 0; i < len(api.Locations); i++ {
			idxPipe := strings.Index(Tracker.Locations.Index[index].Locations[i], "|")
			if Tracker.Locations.Index[index].Locations[i][:idxPipe-1] == city {
				tableArtist = append(tableArtist, Tracker.Artists[index])
				continue
			}
		}
	}
	var API CityAPI
	API.Artists = tableArtist
	return API
}

// errorHandler -> charge la page /error
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
	staticFile()
	serverJSON()
	http.HandleFunc("/", menuPage)
	http.HandleFunc("/groupie/", groupiePage)
	http.HandleFunc("/artist/", artistPage)
	http.HandleFunc("/concertLocation/", concertLocationPage)
	http.HandleFunc("/cityConcert/", cityConcertPage)
	http.ListenAndServe(":8080", nil)
}

/***************DEEZER API***************/

func Deezer(selectedArtist *pageArtist) {
	nameArtist := SearchNameID(selectedArtist.SpecialData.Artists[0].ID)
	URLTracklist := SearchArtistDeezer(nameArtist, selectedArtist)
	unmarshallJSON(URLTracklist, &selectedArtist.Deezer.ListSong) //Range dans la structure ListSong, tous les sons
}

//SearchNameID -> Trouve le nom pour l'id et le modifie
func SearchNameID(ID int) string {
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

	unmarshallJSON(url, &StructPageArtist.Deezer.DeezerArtist)
	unmarshallJSON(url, &StructPageArtist.Deezer.DeezerArtist)

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
	CoverURL       string `json:"cover_xl"`
	TrackListAlbum string `json:"tracklist"`
}
