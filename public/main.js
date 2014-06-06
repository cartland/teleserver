(function($) {

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

    // Populate the graph with initial data
    function onJSONFetch(series) {
        if (!fetched[series.label]) {
            fetched[series.label] = true;
            data[series.label] = series;
            console.log(series)
            plot.setData([series]);
            plot.setupGrid();
            plot.draw();
        }
    }

    // Fetch the initial data
    $.ajax({
        url: "/data/speed.json",
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
                series.shift();

                // TODO(stvn): Support multiple graphs and base stuff on time
                plot.setData([data[name]]);
                plot.setupGrid();
                plot.draw();
            }
        }

        // Create a websocket connection and do live updates of the data
        var ws = new WebSocket("ws://" + window.location.host + "/ws");
        ws.onmessage = function(e) {
            var data = JSON.parse(e.data);
            if (data.type) {
                $("#" + data.type).text(data.value.toFixed(1));
                update(data.type, [data.time, data.value]);
            }
        }
    });
})(jQuery);
