/* Javascript : Choisi un nombre au hasard et renvoie sur la page de l'artidte correspondant */
function randomArtist() {
    /* Choisi un nombre entier entre 0 et 52 */
    var random = Math.floor(Math.random()*52);
    /* Renvoie sur la page de l'artiste correspondant à l'id */
    document.location.href="http://localhost:8080/artist/"+random;
}

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
    document.location.href="http://localhost:8080/error/";
}

/* Javascript : Transmorme les données de la datalist pour les recherches tapéees à la main
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
    var map = L.map('mapID').setView([lat, lon], 16);

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
        iconAnchor: [25, 50]
    });

    /* Récupère les données se trouvant dans les objet HTML ayant la classe geocities */
    var geo = document.getElementsByClassName('geocities');
    for(var i = 0; i < geo.length; i++){
        /* Récupère le texte se trouvant dans la donnée geo à l'index i */
        var geoLocation = transformText(geo[i].innerHTML);
        /* URL de l'API concernant le bon lieu de concert */
        var url = 'https://nominatim.openstreetmap.org/search.php?q='+geoLocation+'&polygon_geojson=1&format=jsonv2';
        /* Création et envoie de la requête API */
        var request = new XMLHttpRequest();
        request.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                /* Récupère les données de la requête */
                var coords = JSON.parse(this.responseText);
                /* Ajoute à la map le marqueur personnalisé aux lieux des concerts */
                L.marker([coords[0].lat, coords[0].lon], {icon: pinDesign}).addTo(map);
            }
        };
        request.open('GET', url, true);
        request.send();
    }
}

/* Javascript : Transforme l'orthographe des lieux des concerts afin qu'ils soient conformnt pour être utilisés par l'API */
function transformText(text) {
    var newText = "";
    for(i = 0; i < text.length; i++) {
        if(text[i] == "_") {
            newText = newText + "-";
        } else {
            newText = newText + text[i];
        }
    }
    return newText;
}