function randomArtist() {
    var random = Math.floor(Math.random()*52);
    document.location.href="http://localhost:8080/artist/"+random;
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
        let lat = position.coords.latitude;
        let lon = position.coords.longitude;
        var map = L.map('mapID').setView([lat, lon], 2);

        L.tileLayer('https://{s}.tile.openstreetmap.fr/osmfr/{z}/{x}/{y}.png', {
            // Il est toujours bien de laisser le lien vers la source des données
            attribution: 'données © OpenStreetMap/ODbL - rendu OSM France',
            minZoom: 1,
            maxZoom: 20
        }).addTo(map);

        var pinDesign = L.icon({
            iconUrl: "/static/images/pin.svg",
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
                    L.marker([coords[0].lat, coords[0].lon], {icon: pinDesign}).addTo(map);
                    console.log(coords);
                }
            };
            request.open('GET', url, true);
            request.send();
        }
    });
}