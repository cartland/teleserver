$(function() {
    var data = [],
        maxPoints = 300;
    var plot = $.plot("#speed-plot", [data], {
        legend: {
            show: true
        },
        series: {
            shadowSize: 0 // Drawing is faster without shadows
        },
        xaxis: {
            mode: "time",
            timeformat: "%H:%M:%S",
            timezone: "browser"
        }
    });

    function update(point) {
        data.push([Date.now(), point])
        while (data.length > maxPoints) {
            data.shift()
        }
        plot.setData([{
            label: "Speed",
            data: data
        }]);
        plot.setupGrid()
        plot.draw();
    }

    var ws = new WebSocket("ws://" + window.location.host + "/ws");
    ws.onmessage = function(e) {
        var data = JSON.parse(e.data)
        if (data.type) {
            $("#" + data.type).text(data.value.toFixed(1))
        }
        if (data.type === "speed") {
            update(data.value)
        }
    }
});
