var plotDefaults = {
    legend: {
        show: true
    },
    series: {
        shadowSize: 0 // Drawing is faster without shadows
    },
    grid: {
        hoverable: true,
        mouseActiveRadius: 30,
    },
    xaxis: {
        mode: "time",
        timeformat: "%H:%M:%S",
        timezone: "browser"
    }
}

var getNum = function(str) {
    var i = parseFloat(str);
    return isNaN(i) ? 0 : i;
}

var tooltip = function() {
    function showTooltip(x, y, contents) {
        $('<div id="tooltip">' + contents + '</div>').css({
            position: 'absolute',
            display: 'none',
            top: y + 5,
            left: x + 5,
            border: '1px solid #ddf',
            padding: '2px',
            'background-color': '#eef',
            opacity: 0.80
        }).appendTo("body").fadeIn(0);
    }

    var previousPoint
    return function(event, pos, item) {
        if (item) {
            if (previousPoint != item.dataIndex) {
                previousPoint = item.dataIndex;

                $("#tooltip").remove();
                var x = item.datapoint[0].toFixed(2),
                    y = item.datapoint[1].toFixed(2);

                showTooltip(item.pageX, item.pageY,
                    item.series.label + ": " + y);
            }
        } else {
            $("#tooltip").remove();
            previousPoint = null;
        }
    }
}
