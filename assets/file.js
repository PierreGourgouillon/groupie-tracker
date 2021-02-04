function randomArtist() {
    var random = Math.floor(Math.random()*52);
    document.location.href="http://localhost:8080/artist/"+random;
}

function displayDate(thingId) {
    var targetElement;
    targetElement = document.getElementById(thingId) ;
    if (targetElement.style.display == "none")
    {
        targetElement.style.display = "" ;
    } else {
        targetElement.style.display = "none" ;
    }
}

function searchBar() {
    var valInput = document.getElementById('search').value;
    var valData = document.getElementById('searchs').options;
    for(i = 0; i < valData.length; i++) {
        if(valData[i].value.toLowerCase() == valInput.toLowerCase()) {
            var ID = valData[i].dataset.id
            document.location.href="http://localhost:8080/artist/"+ID;
            return;
        }
    }
    document.location.href="http://localhost:8080/error/";
}

function mapConcert() {
    let geoActivate;
    if('geolocation' in navigator) {
        geoActivate = true;
    } else {
        geoActivate = false;
    }

    navigator.geolocation.getCurrentPosition(position => {
        if(geoActivate) {
            const lat = position.coords.latitude;
            const lon = position.coords.longitude;
        } else {
            const lat = 48.852969;
            const lon = 2.349903;
        }
        var map = L.map('mapID').setView([lat, lon], 2);

        L.tileLayer('https://{s}.tile.openstreetmap.fr/osmfr/{z}/{x}/{y}.png', {
            // Il est toujours bien de laisser le lien vers la source des données
            attribution: 'données © OpenStreetMap/ODbL - rendu OSM France',
            minZoom: 1,
            maxZoom: 20
        }).addTo(map);

        var pinDesign = L.icon({
            iconUrl: '/static/images/pin.png',
            iconSize: [50, 50],
            iconAnchor: [25, 50]
        });

        var geo = document.getElementsByClassName('geocities');
        for(var i = 0; i < geo.length; i++){
            var url = 'https://nominatim.openstreetmap.org/search.php?q='+geo[i].innerHTML+'&polygon_geojson=1&format=jsonv2';
            var request = new XMLHttpRequest();
            request.onreadystatechange = function() {
                if (this.readyState == 4 && this.status == 200) {
                    var coords = JSON.parse(this.responseText);
                    var pin = L.marker([coords[0].lat, coords[0].lon]).addTo(map);
                    console.log(coords);
                }
            };
            request.open('GET', url, true);
            request.send();
        }
    });
}