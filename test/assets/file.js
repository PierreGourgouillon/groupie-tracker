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

function petitTest() {
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

//window.onload = () => {
    //var concertLocations = ["paris-france", "tours-france", "lyon-france"];

    



    // for(location in concertLocations) {
    //     // Conversion des lieux des concerts en coordonnées GPS (latitude, longitude)
    //     ajaxGet(`https://nominatim.openstreetmap.org/search?q=washington-usa&format=json&addressdetails=1&limit=1&polygon_svg=1`).then(reponse => {
    //         var coords = JSON.parse(reponse)
    //         var pin = L.marker([coords[0].lat, coords[0].lon], {icon: pinImage}).addTo(map);
    //     })
    // }

//}
