var ws = new WebSocket("ws://" + window.location.host + "/ws");
ws.onmessage = function(e) {
    var data = JSON.parse(e.data)
    if (data.type) {
        document.getElementById(data.type).innerText = data.value.toFixed(1)
    }
}
