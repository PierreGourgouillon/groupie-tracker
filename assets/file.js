/************* Random Artist *************/

/* Javascript : Choisi un nombre au hasard et renvoie sur la page de l'artidte correspondant */
function randomArtist() {
    /* Choisi un nombre entier entre 1 et 52 */
    var random = Math.floor(Math.random()*52)+1;
    /* Renvoie sur la page de l'artiste correspondant à l'id */
    document.location.href="http://localhost:8080/artist/"+random;
}

/************* Search Bar *************/

/* Javacript : Compare la valeur d'entrée avec les valeurs de la datalist pour renvoyer l'utilisateur sur la bonne page */
function searchBar() {
    /* Récupère la valeur de l'objet HTML ayant pour id "search" */
    var valInput = document.getElementById('search').value;
    /* Récupère les options liées à l'objet HTML ayant pour id "searchs" */
    var valData = document.getElementById('searchs').options;
    for(let i = 0; i < valData.length; i++) {
        var valTransformedData = transformSearch(valData[i].value.toLowerCase());
        if(valData[i].value.toLowerCase() == valInput.toLowerCase() || valTransformedData == valInput.toLowerCase()) {
            /* Récupère le data-id de l'option correspondante */
            var ID = valData[i].dataset.id;
            /* Renvoie sur la page de l'artiste correspondant à l'id */
            document.location.href="http://localhost:8080/artist/"+ID;
            return;
        }
    }
    /* Renvoie sur la page erreur 404 */
    document.location.href="http://localhost:8080/deezer/"+ valInput.toLowerCase();
}

function searchBarArtistdeezer(){
    var valData = document.getElementById('ArtistDeezerInput').value
    document.location.href="http://localhost:8080/deezer/"+ valData.toLowerCase();
}

/* Javascript : Transforme les données de la datalist pour les recherches tapéees à la main
(Supprime la partie à partir du |) */
function transformSearch(text) {
    var newText = "";
    for(i = 0; i < text.length; i++) {
        if(text[i] == " " && text[i+1] == "|") {
            console.log(newText);
            return newText;
        } else {
            newText = newText + text[i];
        }
    }
}

/************* Map sur la page artist.html *************/

/* Javascript : Centrage de la map en fonction des options de l'utilisateur */
function mapConcert() {
    let lat = 0;
    let lon = 0;
    /* Si la geolocalisation est activée pour le navigateur utilisé */
    if('geolocation' in navigator) {
        /* Prend la position de l'appareil utilisé et charge la map centrée sur cette position */
        navigator.geolocation.getCurrentPosition(position => {
            loadMap(position.coords.latitude, position.coords.longitude)
        });
    /* Sinon charge la map centrée sur Paris */
    } else {
        loadMap(48.85341, 2.34880);
    }
}

/* Javascript : Création et chargement de la map */
function loadMap(lat, lon) {
    console.log(lat, lon);
    /* Créer la map */
    let map = L.map('mapID').setView([lat, lon], 3);

    /* Ajoute à la map les cartes sur openstreetmap */
    L.tileLayer('https://{s}.tile.openstreetmap.fr/osmfr/{z}/{x}/{y}.png', {
        // Il est toujours bien de laisser le lien vers la source des données
        attribution: 'données © OpenStreetMap/ODbL - rendu OSM France',
        minZoom: 1,
        maxZoom: 20
    }).addTo(map);

    /* créer un pin personnalisé à mettre sur la carte */
    var pinDesign = L.icon({
        iconUrl: "/static/images/pin.svg",
        iconSize: [50, 50],
        iconAnchor: [25, 50],
        popupAnchor: [-3, -76]
    });

    /* Récupère les données se trouvant dans les objet HTML ayant la classe geocities */
    var geo = document.getElementsByClassName('geocities');
    for(var i = 0; i < geo.length; i++){
        /* Récupère le texte se trouvant dans la donnée geo à l'index i */
        let geoLocation = transformText(geo[i].innerHTML);
        /* URL de l'API concernant le bon lieu de concert */
        var url = 'https://nominatim.openstreetmap.org/search.php?q='+geoLocation+'&polygon_geojson=1&format=jsonv2';
        /* Création et envoie de la requête API */
        var request = new XMLHttpRequest();
        request.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                /* Récupère les données de la requête */
                var coords = JSON.parse(this.responseText);
                /* Ajoute à la map le marqueur personnalisé aux lieux des concerts */
                let pin = L.marker([coords[0].lat, coords[0].lon], {icon: pinDesign}).addTo(map);
                var textPopup = transformTextPopup(geoLocation)
                pin.bindPopup(textPopup);

                pin.on('click', function(e){
                    map.setView(e.latlng, 13);
                });
            }
        };
        request.open('GET', url, true);
        request.send();
    }

    map.on('popupclose', function(e) {
        map.setView([48.85341, 2.34880], 3);
    });
}

/************* Map sur la page cityConcert.html *************/
/* Javascript : Création et chargement de la map */
function loadMapCity(city) {
    /* Créer la map */
    var map = L.map('mapcity').setView([0, 0], 12);

    /* Ajoute à la map les cartes sur openstreetmap */
    L.tileLayer('https://{s}.tile.openstreetmap.fr/osmfr/{z}/{x}/{y}.png', {
        // Il est toujours bien de laisser le lien vers la source des données
        attribution: 'données © OpenStreetMap/ODbL - rendu OSM France',
        minZoom: 1,
        maxZoom: 20
    }).addTo(map);

    console.log(city)

    var geoLocation = transformText(city);
    /* URL de l'API concernant le bon lieu de concert */
    var url = 'https://nominatim.openstreetmap.org/search.php?city='+geoLocation+'&polygon_geojson=1&format=jsonv2'
    /* Création et envoie de la requête API */
    var request = new XMLHttpRequest();
    request.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            /* Récupère les données de la requête */
            var coords = JSON.parse(this.responseText);
            map.panTo([coords[0].lat, coords[0].lon])
            var geojson = L.geoJSON(coords[0].geojson, {
                style: {
                    "color": '#3A1757',
                    "opacity": 1,
                    "weight": 1,
                    "fillColor": '#F98718',
                    "fillOpacity": 0.5
                }
            }).addTo(map);
        }
    };
    request.open('GET', url, true);
    request.send();
}

/* Javascript : Transforme la syntaxe des lieux des concerts afin qu'ils soient conforment pour être utilisés par l'API */
function transformText(text) {
    var newText = "";
    for(i = 0; i < text.length; i++) {
        if(text[i] == " ") {
            newText = newText + "-";
        } else if(text[i] == "|")  {
            i++
        } else {
            newText = newText + text[i];
        }
    }
    return newText;
}

/* Javascript : Transforme la syntaxe des lieux des concerts pour qu'ils soent plus propres dans les popups */
function transformTextPopup(text) {
    table = text.split("-");
    for (let i = 0; i < table.length; i++) {
        table[i] = table[i][0].toUpperCase() + table[i].substr(1);
    }
    return table.join(" ");
}

//Javascript : évite d'appuyer sur enter dans le formulaire de filtre
function checkEnter(event){
    if (event.keyCode == 13){
    return false;
    }
}

//Javascript: Rends les élèments invisibles lorsque l'on clique sur le boutton
function printFilter(){
    document.getElementById("form-filters").hidden = false;
    document.getElementById("buton-hidden").hidden = true;
    document.getElementById("random-artist").hidden = true;
}

function checkboxHiddenMembers(){
    var a = document.getElementById("idmembers").hidden

    if (a == true){
        document.getElementById("idmembers").hidden = false;
        document.getElementById("membersSlider").hidden = false;
    }else {
        document.getElementById("idmembers").hidden = true;
        document.getElementById("membersSlider").hidden = true;
    }
}

function checkboxHiddencity(){
    var a = document.getElementById("idcity").hidden

    if (a == true){
        document.getElementById("idcity").hidden = false;
        document.getElementById("citySlider").hidden = false;
    }else {
        document.getElementById("idcity").hidden = true;
        document.getElementById("citySlider").hidden = true;
    }
}

/************* concertLocation *************/
function goToConcert(location) {
    console.log(location);
    document.location.href="http://localhost:8080/artist/"+location;
}

/************** filter *********************/
function displayFilter(thingId) {
    let targetElement;
    let flexbox;
    targetElement = document.getElementById(thingId) ;
    flexbox = document.getElementsByClassName('flex-box');
    if (targetElement.style.display == "none")
    {
        window.scrollTo(0,0);
        targetElement.style.display = "" ;
        flexbox[0].style.marginTop = "5%";
    } else {
        targetElement.style.display = "none" ;
        flexbox[0].style.marginTop = "3%";
    }
}

/************** player deezer ****************/
function playMusic(ID) {
    var table = []
    var f = table.unshift(ID);
    DZ.player.playTracks(table);
}

function playAlbum(ID) {
    IDAlbum = parseInt(ID);
    DZ.player.playAlbum(IDAlbum);
}