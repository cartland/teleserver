(function($) {

    // Handle tab switching
    var home = $('#home-container');
    var solar = $('#solar-container');
    var dash = $('#dash-container');
    var batt = $('#batt-container');

    function switchTabs(newTab) {
        home.hide();
        solar.hide();
        dash.hide();
        batt.hide();
        switch (newTab) {
            case 'home':
                home.show();
                break;
            case 'solar':
                solar.show();
                break;
            case 'dash':
                dash.show();
                break;
            case 'batt':
                batt.show();
                break;
        }
    }
    var tabs = document.querySelector('paper-tabs');

    if (!window.location.hash) {
        window.location.hash = 'home';
    }
    tabs.selected = window.location.hash.slice(1);
    switchTabs(tabs.selected);

    tabs.addEventListener('core-select', function() {
        window.location.hash = tabs.selected;
        switchTabs(tabs.selected);
    });

    var graphs = new Graphs({
        'placeholder': [
            [0x402, 'BusVoltage', 'Voltage'],
            [0x402, 'BusCurrent', 'Current'],
            [0x403, 'VehicleVelocity', 'Velocity']
        ]
    }, '2m');

    $(function() {

        // Update the graph with the given point
        function update(name, point) {
            if (fetched[name]) {
                series = data[name].data;
                series.push(point);
                while (series.length > 1 &&
                    point[0] - series[0][0] > bufferedTime) {
                    series.shift();
                }

                // TODO(stvn): Support multiple graphs
                refreshPlot();
            }
        }

        // Create a websocket connection and do live updates of the data
        var ws = new WebSocket('ws://' + window.location.host + '/ws');
        ws.onmessage = function(e) {
            var data = JSON.parse(e.data);

            if (data.time) {
                data.time = (new Date(data.time)).getTime();
            }
            if (data.CAN) {
                for (var key in data.CAN) {
                    if (data.CAN.hasOwnProperty(key)) {
                        var val = data.CAN[key];
                        var id = getIdForMsg(data.canID, key);

                        if ($('#' + id).length) {
                            $('#' + id).text(val.toFixed(1));
                        }
                        graphs.update(id, [data.time, val]);
                    }
                }
            }
        };
    });
})(jQuery);
