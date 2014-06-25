(function($) {
    // Create the graph
    var plot = $.plot("#placeholder", [], plotDefaults);
    $("#placeholder").bind("plothover", tooltip());


    // Populate the graph with initial data
    function onJSONFetch(allSeries) {
        plot.setData(allSeries);
        plot.setupGrid();
        plot.draw();
    }

    // Fetch the initial data
    var canids = [0x130, 0x131, 0x132, 0x140, 0x141, 0x142].join("&canid=");
    var fields = ["Voltage0", "Voltage1", "Voltage2", "Voltage3"].join("&field=");
    $.ajax({
        url: "/api/graphs?canid=" + canids + "&field=" + fields + "&time=1m",
        type: "GET",
        dataType: "json",
        success: onJSONFetch
    });

    $(function() {
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
                        key = "0x" + parseInt(data.canID, 10).toString(16) + key
                        if ($("#" + key).length) {
                            $("#" + key).text(val);
                        }
                    }
                }
            }

        };
    });
})(jQuery);
