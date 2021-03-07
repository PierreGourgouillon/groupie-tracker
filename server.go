package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

/******************************************TOUTES LES STRUCTURES*****************************************/

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

/******************************************STRUCTURE POUR LA PAGE ARTIST******************************************/

// pageArtists -> structures pour envoyer les informations à la page /artist
type pageArtist struct {
	Data         API
	SpecialData  ArtistAPI
	Flag         []string
	AllLocations []string
	Deezer       DeezerAPI
}

/******************************************STRUCTURE API DEEZER******************************************/

//DeezerAPI -> Structure les données Deezer de l'artiste
type DeezerAPI struct {
	DeezerArtist ArtistDeezer
	ListSong     ListSong
	AlbumInfo    Album
	ListAlbum    []Album
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
		DurationMinute string
		Album          Album `json:"album"`
	} `json:"data"`
}

//Album -> Information sur l'album
type Album struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	CoverURL       string `json:"cover_big"`
	TrackListAlbum string `json:"tracklist"`
}

/******************************************FIN STRUCTURE API DEEZER******************************************/

// ArtistAPI -> structure Artiste spécial pour la page /artist contient un seul artist, ses lieux de concert et leurs dates
type ArtistAPI struct {
	ID        int
	Artists   Artists
	Locations []string
	Relation  map[string][]string
}

type pageGroupie struct {
	Data         API
	Cities       []citySearch
	AllLocations []string
}

type pageFilter struct {
	Data         API
	SpecialData  ArtistAPI
	FilterEmpty  bool
	AllLocations []string
}

/******************************************STRUCTURE POUR LA PAGE CONCERTLOCATION******************************************/

// pageConcert -> structure pour envoyer les informations à la page /concertLocations
type pageConcert struct {
	Data         API
	SpecialData  ConcertAPI
	Flag         []string
	AllLocations []string
}

// ConcertAPI -> structure spécial pour la page /concertLocations contient les lieux de concert
type ConcertAPI struct {
	ID        int
	Locations []string
}

/******************************************STRUCTURE POUR LA PAGE CITYCONCERT******************************************/

// pageCity -> structure pour envoyer les informations à la page /cityConcert
type pageCity struct {
	Data         API
	SpecialData  CityAPI
	City         string
	AllLocations []string
}

// CityAPI -> structure spécial pour la page /cityConcert des artistes
type CityAPI struct {
	ID      int
	Artists Artists
}

// filter -> structure pour recevoir les informations du filtre
type filter struct {
	typeArtist        string
	FirstAlbum        string
	creationDate      string
	checkMembers      string
	members           string
	checkCity         string
	city              string
	citySearchFilter0 string
	citySearchFilter1 string
	citySearchFilter2 string
}

// citySearch ->
type citySearch struct {
	ID   string
	City string
}

// Variable Globale

// Tracker -> Variable contenant toutes les informatiions de l'API groupie-tracker
var Tracker API

// allLocations -> Variable contenant tous les lieux de concert en un seul exemplaire
var allLocations []string

// serachBarallLocations -> Variable contenant tous les lieux de concert en un seul exemplaire pour la searchbar
var searchBarAllLocations []string

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
	allLocationsFilter()

	fmt.Println("Server UP")
}

// allLocationsFilter -> créer un tableau contenant les lieux des concert en un seul exemplaire pour la search bar
func allLocationsFilter() {
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
	searchBarAllLocations = tableLocations
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

	allInformations := pageGroupie{
		Data:         Tracker,
		AllLocations: searchBarAllLocations,
		Cities:       structCity,
	}

	filterAPI := filter{
		typeArtist:        r.FormValue("typeArtist"),
		FirstAlbum:        r.FormValue("firstAlbum"),
		creationDate:      r.FormValue("creationDate"),
		checkMembers:      r.FormValue("checkMembers"),
		members:           r.FormValue("members"),
		checkCity:         r.FormValue("checkCity"),
		city:              r.FormValue("city"),
		citySearchFilter0: r.FormValue("citySearchFilter"),
		citySearchFilter1: r.FormValue("citySearchFilter1"),
		citySearchFilter2: r.FormValue("citySearchFilter2"),
	}

	if filterAPI.citySearchFilter0 == filterAPI.citySearchFilter1 {
		filterAPI.citySearchFilter1 = ""
	}
	if filterAPI.citySearchFilter0 == filterAPI.citySearchFilter2 {
		filterAPI.citySearchFilter2 = ""
	}
	if filterAPI.citySearchFilter1 == filterAPI.citySearchFilter2 {
		filterAPI.citySearchFilter2 = ""
	}

	if filterAPI.FirstAlbum != "" || filterAPI.creationDate != "" || filterAPI.checkMembers != "" || filterAPI.checkCity != "" || filterAPI.typeArtist != "" || filterAPI.citySearchFilter0 != "" || filterAPI.citySearchFilter1 != "" || filterAPI.citySearchFilter2 != "" {
		structTest := filters(filterAPI)
		template, _ := template.ParseFiles("./templates/filter.html")
		structTest.AllLocations = searchBarAllLocations
		template.Execute(w, structTest)
		return
	}
	tmpl.Execute(w, allInformations)
}

// filters -> filtre en fonction des données reçues
func filters(filterAPI filter) pageFilter {
	var isFilterCreation bool
	var isFilterAlbum bool
	var isFilterMembers bool
	var isFilterCity bool
	var isFilterCitySearch bool
	var isArtistOrBand bool
	var table Artists
	var apiFilter pageFilter

	for i, b := range Tracker.Artists {
		if filterAPI.typeArtist == "Artiste" {
			isArtistOrBand = filterArtistOrBand(filterAPI.typeArtist, len(b.Members))
			if !isArtistOrBand {
				continue
			}
		} else if filterAPI.typeArtist == "Groupe" {
			isArtistOrBand = filterArtistOrBand(filterAPI.typeArtist, len(b.Members))
			if !isArtistOrBand {
				continue
			}
		}

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

		if filterAPI.citySearchFilter0 != "" || filterAPI.citySearchFilter1 != "" || filterAPI.citySearchFilter2 != "" {
			isFilterCitySearch = filterSearchCity(b.ID, filterAPI)

			if !isFilterCitySearch {
				continue
			}
		}
		table = append(table, Tracker.Artists[i])
	}

	if table == nil {
		apiFilter.FilterEmpty = true
	}

	apiFilter.SpecialData.Artists = table
	return apiFilter
}

func filterSearchCity(ID int, API filter) bool {
	var table []string
	a := ""

	if API.citySearchFilter0 != "" {
		table = append(table, API.citySearchFilter0)
	}

	if API.citySearchFilter1 != "" {
		table = append(table, API.citySearchFilter1)
	}

	if API.citySearchFilter2 != "" {
		table = append(table, API.citySearchFilter2)
	}

	for i := range table {
		stringID := table[i]
		for j := range stringID {
			if stringID[j] != ',' {
				a += string(stringID[j])
			} else {
				number, _ := strconv.Atoi(a)
				if number == ID {
					return true
				}
				a = ""
				number = 0
			}
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

func filterArtistOrBand(s string, nbr int) bool {
	if s == "Artiste" {
		if nbr == 1 {
			return true
		}
		return false
	} else if s == "Groupe" {
		if nbr > 1 {
			return true
		}
		return false
	}
	return true
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
	if err2 != nil || nbrPath > 52 {
		if r.URL.Path[8:12] == "city" {
			nbrPathCity, err3 := strconv.Atoi(r.URL.Path[12:])
			if err3 != nil {
				errorHandler(w, r)
				return
			}

			city := deleteCountry([]string{searchBarAllLocations[nbrPathCity]})
			CityTracker := cityConcertFilter(city[0])

			selectedCity := pageCity{
				Data:         Tracker,
				SpecialData:  CityTracker,
				City:         city[0],
				AllLocations: searchBarAllLocations,
			}

			tmplCity, err4 := template.ParseFiles("./templates/cityConcert.html")
			if err4 != nil {
				fmt.Println(err4)
				os.Exit(1)
			}

			tmplCity.Execute(w, selectedCity)
			return

		}

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

	selectedArtist := pageArtist{
		Data:         Tracker,
		SpecialData:  ArtistTracker,
		AllLocations: searchBarAllLocations,
	}
	var selectedArtistPointer *pageArtist
	selectedArtistPointer = &selectedArtist
	deezer(selectedArtistPointer, "")

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
		Data:         Tracker,
		SpecialData:  LocationsTracker,
		Flag:         TableFlag,
		AllLocations: searchBarAllLocations,
	}
	tmpl.Execute(w, locationsConcert)
}

// allLocationsFilter -> créer un tableau contenant les lieux des concert en un seul exemplaire, l'insert dans la variable globale allLocations
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

// flagCountryFilter -> créer un tableau contenant les codes des pays (de l'API flagcdn) dans l'ordre des lieux de concert
func flagCountryFilter(locations []string) []string {
	var TableFlag []string
	for _, location := range locations {
		idxPipe := strings.Index(location, "|")
		if idxPipe == -1 {
			continue
		}
		location = location[idxPipe+2:]
		switch location {
		case "usa":
			TableFlag = append(TableFlag, "us")
		case "us":
			TableFlag = append(TableFlag, "us")
		case "uk":
			TableFlag = append(TableFlag, "gb")
		case "netherlands antilles":
			TableFlag = append(TableFlag, "nl")
		case "czech republic":
			TableFlag = append(TableFlag, "cz")
		case "brasil":
			TableFlag = append(TableFlag, "br")
		case "philippine":
			TableFlag = append(TableFlag, "ph")
		case "korea":
			TableFlag = append(TableFlag, "kr")
		default:
			for countryCode, value := range flagCountry {
				value = strings.ToLower(value)
				if value == location {
					TableFlag = append(TableFlag, countryCode)
					break
				}
			}
		}
	}
	return TableFlag
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
		Data:         Tracker,
		SpecialData:  CityTracker,
		City:         city,
		AllLocations: searchBarAllLocations,
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
	http.HandleFunc("/deezer/", deezerPage)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println(port)
	http.ListenAndServe(":"+port, nil)
}

/******************************************TRANSFORMATION SYNTHAXIQUE******************************************/

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

//transformLocation -> remplace les "_" par un " " et les "-" par un "|"
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

//conversionDeezer -> Convertie les durées de son de secondes à minutes
func conversionDeezer(selectedArtist *pageArtist) {
	var conversion string
	durationSong := 0
	minutes := 0
	secondes := 0
	for i := 0; i < len(selectedArtist.Deezer.ListSong.Data); i++ {
		durationSong = selectedArtist.Deezer.ListSong.Data[i].Duration
		minutes = durationSong / 60
		secondes = durationSong % 60
		if secondes < 10 {
			conversion = strconv.Itoa(minutes) + ":" + "0" + strconv.Itoa(secondes)
			selectedArtist.Deezer.ListSong.Data[i].DurationMinute = conversion

		} else {
			conversion = strconv.Itoa(minutes) + ":" + strconv.Itoa(secondes)
			selectedArtist.Deezer.ListSong.Data[i].DurationMinute = conversion
		}
	}
}

/******************************************DEEZER API******************************************/

//deezerPage -> Page pour les artistes deezer
func deezerPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/deezer.html")

	var artistInDeezer *pageArtist
	var structTemporary pageArtist

	structTemporary.Data = Tracker
	structTemporary.AllLocations = searchBarAllLocations
	artistInDeezer = &structTemporary

	deezerArtist := r.URL.Path[8:]

	deezer(artistInDeezer, deezerArtist)
	conversionDeezer(artistInDeezer)

	tmpl.Execute(w, artistInDeezer)
}

//SearchNameID -> Trouve le nom pour l'id et le modifie
func searchNameID(ID int) string {
	for _, data := range Tracker.Artists {
		if data.ID == ID {
			nameModif := strings.Replace(data.Name, " ", "%20", -1)
			return nameModif
		}
	}
	return ""
}

//SearchArtistDeezer -> Trouve l'artiste dans l'API Deezer
func searchArtistDeezer(name string, StructPageArtist *pageArtist) string {
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

// deezer -> Crée l'artiste via l'API Deezer
func deezer(selectedArtist *pageArtist, artistInSearchBar string) {
	var nameArtist string

	if artistInSearchBar != "" {
		nameArtist = strings.Replace(artistInSearchBar, " ", "%20", -1)
	} else {
		nameArtist = searchNameID(selectedArtist.SpecialData.Artists[0].ID)
	}

	URLTracklist := searchArtistDeezer(nameArtist, selectedArtist)
	unmarshallJSON(URLTracklist, &selectedArtist.Deezer.ListSong) //Range dans la structure ListSong, tous les sons
	listAlbum(selectedArtist)
}

//listAlbum -> Crée la liste des albums de l'artiste
func listAlbum(selectedArtist *pageArtist) {
	var table []string
	var tableau []Album
	var isFirstAlbum bool = true

	for i, b := range selectedArtist.Deezer.ListSong.Data {
		if i == 0 && isFirstAlbum {
			table = append(table, b.Album.Title)
			isFirstAlbum = false
		}

		for _, k := range table {
			if b.Album.Title == k {
				isFirstAlbum = true
				break
			} else {
				isFirstAlbum = false
			}
		}

		if !isFirstAlbum {
			table = append(table, b.Album.Title)
			tableau = append(tableau, b.Album)
		}
	}
	selectedArtist.Deezer.ListAlbum = tableau
}
