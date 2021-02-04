function randomArtist() {
    var random = Math.floor(Math.random()*52);
    document.location.href="http://localhost:8080/artist/"+random;
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

