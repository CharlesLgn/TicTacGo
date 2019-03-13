function play(x, y) {
    $.get('/play?x=' + x + "&y=" + y, function (data) {
        if (data.length > 1) {
            restart();
        } else {
            let button = document.getElementById(x + "" + y);
            button.style.color = data === "X" ? "red" : "blue";
            editButton(button, true, data);
        }
    });
}
function score() {
    $.get('/score', function (data) {
        document.getElementById("res").innerHTML = data;
    });
}
function restart() {
    $.get('/restart', function () {
        for (let i = 0; i < 3; i++) {
            for (let j = 0; j < 3; j++) {
                editButton(document.getElementById(i + "" + j), false, "");
            }
        }
        score();
    });
}

function createGrid() {
    $.get('/start', function (data) {
        document.getElementById("main").innerHTML = data;
    });
}

function editButton(button, disable, txt){
    button.disabled = disable;
    button.innerText = txt;
}

$.get('/title', function (data) {
    $('#image').attr('src', data);
});