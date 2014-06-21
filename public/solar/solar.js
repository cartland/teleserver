(function($) {
    // How long to have graph data
    // This constant should be kept in sync with metrics.go
    var bufferedTime = 20 * 1000 // 20s * 1000 ms/s

    // These metrics must match up with both the graph html ids and the json field
    // names.
    var metrics = ["ArrayVoltage", "ArrayCurrent", "BatteryVoltage", "Temperature"]
    // Create maps to hold the plots, the data for each plot, and the data for each line
    var plots = {}, dataArrays = {}, plotData = {};


    // Refresh the plots 5 times a second.
    window.setInterval(
        function() {
            for (var i = 0; i < metrics.length; i++) {
                var plotname = metrics[i];
                if (plots[plotname] && dataArrays[plotname])
                    var plot = plots[plotname];
                plot.setData(dataArrays[plotname]);
                plot.setupGrid();
                plot.draw();
            }
        }, 200);

    // Populate the initial values of the plots by getting historical data from
    // ajax queries. Also fill in the dataArrays and plotData with this initial
    // information.
    var canids = [0x600, 0x601, 0x602, 0x603].join("&canid=");
    for (var i = 0; i < metrics.length; i++) {
        (function() {
            var metric = metrics[i];
            plots[metric] = $.plot("#" + metric, [], plotDefaults);
            $("#" + metric).bind("plothover", tooltip());
            $.ajax({
                url: "/api/graphs?time=20s&canid=" + canids + "&field=" + metric,
                type: "GET",
                dataType: "json",
                success: function(points) {
                    for (var i = 0; i < points.length; i++) {
                        plotData[points[i].label] = points[i].data
                    }
                    dataArrays[metric] = points
                    // refreshPlot(metric)
                },
            });
        })();
    }

    $(function() {
        var ws = new WebSocket("ws://" + window.location.host + "/ws");
        ws.onmessage = function(e) {
            var data = JSON.parse(e.data);

            if (data.time) {
                data.time = (new Date(data.time));
            }
            if (data.CAN) {
                // Update any graphs that match.
                for (var i = 0; i < metrics.length; i++) {
                    var metric = metrics[i]
                    var name = "0x" + parseInt(data.canID, 10).toString(16) + " - " + metric
                    if (data.CAN[metric] && plotData[name]) {
                        var series = plotData[name];
                        var time = data.time.getTime()
                        series.push([time, data.CAN[metric]]);
                        while (series.length > 1 && time - series[0][0] > bufferedTime) {
                            series.shift();
                        }
                        // refreshPlot(metric)
                    }
                }
            }

        };
    });
})(jQuery);
