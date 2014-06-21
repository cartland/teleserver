(function($) {
    // How long to have graph data
    var bufferedTime = 60 * 1000 // 20s * 1000 ms/s

    // These metrics must match up with both the graph html ids and the json field
    // names.
    var metrics = ["ArrayVoltage", "ArrayCurrent", "BatteryVoltage", "Temperature"]
    // Create maps to hold the plots, the data for each plot, and the data for each line
    var plots = {}, dataArrays = {}, plotData = {};

    // Start following again if we're looking into the future
    function checkFuture(e, plot) {
        var time = (new Date()).getTime()
        if (time < plot.getOptions().xaxes[0].max) {
            plot.getOptions().xaxes[0].max = null;
        }
    }

    // Refresh the plots 5 times a second.
    window.setInterval(
        function() {
            for (var i = 0; i < metrics.length; i++) {
                var plotname = metrics[i];
                if (plots[plotname] && dataArrays[plotname]) {
                    var plot = plots[plotname];
                    var array = dataArrays[plotname]
                    var time = (new Date()).getTime();
                    var min = time;

                    plot.setData(array);
                    for (var j = 0; j < array.length; j++) {
                        if (array[j].data[0][0] < min) {
                            min = array[j].data[0][0];
                        }
                    }
                    if (min > plot.getOptions().xaxes[0].min) {
                        plot.getOptions().xaxes[0].min = null;
                    }
                    plot.getOptions().xaxes[0].panRange = [time - bufferedTime, time + 1000];
                    plot.setupGrid();
                    plot.draw();
                }
            }
        }, 250);

    // Populate the initial values of the plots by getting historical data from
    // ajax queries. Also fill in the dataArrays and plotData with this initial
    // information.
    var canids = [0x600, 0x601, 0x602, 0x603].join("&canid=");
    for (var i = 0; i < metrics.length; i++) {
        (function() {
            var metric = metrics[i];
            var time = (new Date()).getTime();
            var plot = $.plot("#" + metric, [], plotDefaults);
            $("#" + metric).bind("plothover", tooltip());
            $("#" + metric).bind("plotpan", checkFuture);
            $("#" + metric).bind("plotzoom", checkFuture);
            plot.getOptions().xaxes[0].panRange[0] = time - bufferedTime;
            plots[metric] = plot
            $.ajax({
                url: "/api/graphs?time=1m&canid=" + canids + "&field=" + metric,
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
            if (data.CAN && data.canID >= 0x600 && data.canID < 0x620) {
                // Update any graphs and ids that match.
                for (var i = 0; i < metrics.length; i++) {
                    var metric = metrics[i]
                    var name = "0x" + parseInt(data.canID, 10).toString(16) + " - " + metric
                    $("#" + data.CAN.ArrayLocation + metric).text(data.CAN[metric])
                    if (data.CAN[metric] && plotData[name]) {
                        var series = plotData[name];
                        var time = data.time.getTime()
                        series.push([time, data.CAN[metric]]);
                        while (series.length > 1 && time - series[0][0] > bufferedTime) {
                            series.shift();
                        }
                    }
                }

                // Update the total current from the array
                var sum = 0;
                $('.ArrayCurrentVal').each(function() {
                    sum += getNum($(this).text());
                });
                $("#total").text(sum);
            }

        };
    });
})(jQuery);
