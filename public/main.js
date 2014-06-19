(function($) {
    // How long to have graph data
    var bufferedTime = 1 * 60 * 1000 // 1m * 60s/m * 1000 ms/s

    // Create the graph
    var plot = $.plot("#placeholder", [], {
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
    var data = {}, fetched = {};
    var dataArray = [];

    function refreshPlot() {
        plot.setData(dataArray);
        plot.setupGrid();
        plot.draw();
    }
    // Populate the graph with initial data
    function onJSONFetch(allSeries) {
        for (var i = 0; i < allSeries.length; i++) {
            series = allSeries[i]
            if (!fetched[series.label]) {
                fetched[series.label] = true;
                data[series.label] = series;
            }
        }
        dataArray = allSeries
        refreshPlot()
    }

    // Fetch the initial data
    $.ajax({
        url: "/api/graphs?canid=" + 0x402 + "&field=VehicleVelocity&canid=" + 0x403 + "&field=BusVoltage&field=BusCurrent&time=1m",
        type: "GET",
        dataType: "json",
        success: onJSONFetch
    });

    $(function() {

        // Update the graph with the given point
        function update(name, point) {
            if (fetched[name]) {
                series = data[name].data;
                series.push(point);
                while (series.length > 1 && point[0] - series[0][0] > bufferedTime) {
                    series.shift();
                }

                // TODO(stvn): Support multiple graphs
                refreshPlot()
            }
        }

        // Create a websocket connection and do live updates of the data
        var ws = new WebSocket("ws://" + window.location.host + "/ws");
        ws.onmessage = function(e) {
            var data = JSON.parse(e.data);

            if (data.time) {
                data.time = (new Date(data.time)).getTime();
            }
            if (data.CAN) {
                for (var key in data.CAN) {
                    if (data.CAN.hasOwnProperty(key)) {
                        var val = data.CAN[key]
                        if ($("#" + key).length) {
                            $("#" + key).text(val.toFixed(1));
                        }
                        update("0x" + parseInt(data.canID, 10).toString(16) + " - " + key, [data.time, val]);
                    }
                }
            }
        }
    });
})(jQuery);
