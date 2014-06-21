(function($) {
    // How long to have graph data
    // This constant should be kept in sync with metrics.go
    var bufferedTime = 20 * 1000 // 20s * 1000 ms/s

    // Create the graphs
    var plots = {
        "array-voltage": $.plot("#array-voltage", [], plotDefaults),
        "array-current": $.plot("#array-current", [], plotDefaults),
        "battery-voltage": $.plot("#battery-voltage", [], plotDefaults),
        "temperature": $.plot("#temperature", [], plotDefaults),
    }

    var data = {}, fetched = {};
    var dataArray = {};
    $("#array-current").bind("plothover", tooltip());

    function refreshPlot(plotname) {
        var plot = plots[plotname];
        debugger;
        plot.setData(dataArray[plotname]);
        plot.setupGrid();
        plot.draw();
    }
    // Populate the graph with initial data
    function onFetch(plotname) {
        return function(allSeries) {
            dataArray[plotname] = allSeries
            refreshPlot(plotname)
        }
    }

    // Fetch the initial data
    var canids = [0x600, 0x601, 0x602, 0x603].join("&canid=")
    $.ajax({
        url: "/api/graphs?canid=" + canids + "&field=ArrayCurrent",
        type: "GET",
        dataType: "json",
        success: onFetch("array-current"),
    });
    $.ajax({
        url: "/api/graphs?canid=" + canids + "&field=ArrayVoltage",
        type: "GET",
        dataType: "json",
        success: onFetch("array-voltage"),
    });
    $.ajax({
        url: "/api/graphs?canid=" + canids + "&field=BatteryVoltage",
        type: "GET",
        dataType: "json",
        success: onFetch("battery-voltage"),
    });
    $.ajax({
        url: "/api/graphs?canid=" + canids + "&field=Temperature",
        type: "GET",
        dataType: "json",
        success: onFetch("temperature"),
    });

    $(function() {

        // Update the graph with the given point
        function update(name, point) {
            name = "0x" + parseInt(data.canID, 10).toString(16) + " - " + name
            if (fetched[name]) {
                series = data[name].data;
                series.push(point);
                while (series.length > 1 && point[0] - series[0][0] > bufferedTime) {
                    series.shift();
                }

                // TODO(stvn): Support multiple graphs
                refreshPlot(plots["temperature"])
            }
        }

        // Create a websocket connection and do live updates of the data
        var ws = new WebSocket("ws://" + window.location.host + "/ws");
        ws.onmessage = function(e) {
            var data = JSON.parse(e.data);

            if (data.time) {
                data.time = (new Date(data.time)).getTime();
            }
            if (data.canID >= 0x600 && data.canID <= 0x604) {
                update(data.CAN.ArrayLocation, [data.time, data.CAN]);
            }
        }
    });
})(jQuery);
