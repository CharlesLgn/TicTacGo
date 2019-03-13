function play(x, y) {
    $.get('/play?x=' + x + "&y=" + y, function (data) {
        let button = document.getElementById(x + "" + y);
        if (data.length > 1) {
            restart();
        } else {
            button.style.color = data === "X" ? "red" : "blue";
            button.disabled = true;
            button.innerText = data;
        }
    });
}
function score() {
    $.get('/score', function (data) {
        document.getElementById("res").innerHTML = data;
    });
}
function start() {
    window.location.href = "/start";
}
function restart() {
    $.get('/restart', function () {
        for (let i = 0; i < 3; i++) {
            for (let j = 0; j < 3; j++) {
                let button = document.getElementById(i + "" + j);
                button.innerText = "";
                button.disabled = false;
            }
        }
        score();
    });
}