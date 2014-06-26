(function($) {
    // How long to have graph data
    var bufferedTime = "2m"
    if (window.location.hash) {
        bufferedTime = window.location.hash.slice(1);
    }
    var bufferedMillis = parseDuration(bufferedTime);

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
                if (plots[plotname] && dataArrays[plotname]) {
                    var plot = plots[plotname];

                    plot.setData(dataArrays[plotname]);
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
            plot.getOptions().xaxes[0].panRange[0] = time - bufferedMillis;
            plots[metric] = plot;
            $(window).on("hashchange", function() {
                bufferedTime = window.location.hash.slice(1);
                $.ajax({
                    url: "/api/graphs?time=" + bufferedTime + "&canid=" + canids + "&field=" + metric,
                    type: "GET",
                    dataType: "json",
                    success: function(points) {
                        for (var i = 0; i < points.length; i++) {
                            plotData[points[i].label] = points[i].data
                        }
                        dataArrays[metric] = points
                    },
                });
            });
        })();
    }
    $(window).trigger("hashchange")

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
                        while (series.length > 1 && time - series[0][0] > bufferedMillis) {
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
