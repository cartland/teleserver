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
        $('.nav li').removeClass('active')
        $('#' + newTab + '-nav').addClass('active')
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
    var tabs = document.querySelector('.nav a');

    if (!window.location.hash) {
        window.location.hash = 'home';
    }
    tabs.selected = window.location.hash.slice(1);
    switchTabs(tabs.selected);
    $(window).on('hashchange', function() {
        tabs.selected = window.location.hash.slice(1);
        switchTabs(tabs.selected);
    });


    tabs.addEventListener('core-select', function() {
        window.location.hash = tabs.selected;
        switchTabs(tabs.selected);
    });

    var graphs = new Graphs({
            'home-graph': [
                [0x402, 'BusVoltage', 'Voltage'],
                [0x402, 'BusCurrent', 'Current'],
                [0x403, 'VehicleVelocity', 'Velocity']
            ],
            'array-current-graph': [
                [0x600, 'ArrayCurrent', 'Front Right Array Current'],
                [0x601, 'ArrayCurrent', 'Front Left Array Current'],
                [0x602, 'ArrayCurrent', 'Back Right Array Current'],
                [0x603, 'ArrayCurrent', 'Back Left Array Current']
            ],
            'array-voltage-graph': [
                [0x600, 'ArrayVoltage', 'Front Right Array Voltage'],
                [0x601, 'ArrayVoltage', 'Front Left Array Voltage'],
                [0x602, 'ArrayVoltage', 'Back Right Array Voltage'],
                [0x603, 'ArrayVoltage', 'Back Left Array Voltage']
            ],
            'array-battery-voltage-graph': [
                [0x600, 'BatteryVoltage', 'Front Right Battery Voltage'],
                [0x601, 'BatteryVoltage', 'Front Left Battery Voltage'],
                [0x602, 'BatteryVoltage', 'Back Right Battery Voltage'],
                [0x603, 'BatteryVoltage', 'Back Left Battery Voltage']
            ],
            'array-temperature-graph': [
                [0x600, 'Temperature', 'Front Right Array Temperature'],
                [0x601, 'Temperature', 'Front Left Array Temperature'],
                [0x602, 'Temperature', 'Back Right Array Temperature'],
                [0x603, 'Temperature', 'Back Left Array Temperature']
            ],
            'dash-graph': [
                [0x402, 'BusVoltage', 'Voltage'],
                [0x402, 'BusCurrent', 'Current'],
                [0x403, 'VehicleVelocity', 'Velocity']
            ],
            'battery-cells-graph-low': [
                [0x130, 'Voltage0', 'Cell 00'],
                [0x130, 'Voltage1', 'Cell 01'],
                [0x130, 'Voltage2', 'Cell 02'],
                [0x130, 'Voltage3', 'Cell 03'],
                [0x131, 'Voltage0', 'Cell 04'],
                [0x131, 'Voltage1', 'Cell 05'],
                [0x131, 'Voltage2', 'Cell 06'],
                [0x131, 'Voltage3', 'Cell 07']
            ],
            'battery-cells-graph-mid': [
                [0x140, 'Voltage0', 'Cell 08'],
                [0x140, 'Voltage1', 'Cell 09'],
                [0x140, 'Voltage2', 'Cell 10'],
                [0x140, 'Voltage3', 'Cell 11'],
                [0x141, 'Voltage0', 'Cell 12'],
                [0x141, 'Voltage1', 'Cell 13'],
                [0x141, 'Voltage2', 'Cell 14'],
                [0x141, 'Voltage3', 'Cell 15'],
                [0x142, 'Voltage0', 'Cell 16'],
                [0x142, 'Voltage1', 'Cell 17'],
                [0x142, 'Voltage2', 'Cell 18'],
                [0x142, 'Voltage3', 'Cell 19']
            ],
            'battery-cells-graph-high': [
                [0x150, 'Voltage0', 'Cell 20'],
                [0x150, 'Voltage1', 'Cell 21'],
                [0x150, 'Voltage2', 'Cell 22'],
                [0x150, 'Voltage3', 'Cell 23'],
                [0x151, 'Voltage0', 'Cell 24'],
                [0x151, 'Voltage1', 'Cell 25'],
                [0x151, 'Voltage2', 'Cell 26'],
                [0x151, 'Voltage3', 'Cell 27']
            ],
            'battery-current-graph': [
                [0x124, 'Current', 'Battery Current']
            ]
        },
        '2m');

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
                            // Cheating to scale voltage
                            if (id.indexOf("Voltage") > -1) {
                                $('#' + id).text((val / 10000).toFixed(3));
                            } else {
                                $('#' + id).text(val.toFixed(1));
                            }
                        }
                        graphs.update(id, [data.time, val]);
                    }
                }
            }

            // Update the total current from the array
            var sum = 0;
            $('.array-current-val').each(function() {
                sum += getNum($(this).text());
            });
            $('#array-current-total').text(sum);

            // Update the total voltage of the batteries
            sum = 0;
            $('.battery-cell-voltage').each(function() {
                sum += getNum($(this).text());
            });
            $('#battery-voltage-total').text(sum.toFixed(3));
        };
    });
})(jQuery);
