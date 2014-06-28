// Mapping will look like
// {
//   '#graphid': [
//     [0x123, Field, 'Awesome Field'
//   ]
// }
var Graphs = function() {
    function onAjax(plotData) {
        return function(points) {
            for (var i = 0; i < points.length; i++) {
                var data = points[i].data;
                for (var j = 0; j < data.length; j++) {
                    plotData[points[i].label][j] = data[j];
                }
            }
        }
    }

    function makeGraph(graph) {
        var graphid = graph[0];
        var graphData = graph[2];
        return function() {
            if ($('#' + graphid).is(':visible')) {
                graph[1] = $.plot('#' + graphid, graphData, plotDefaults);
                $('#' + graphid).bind('plothover', tooltip());
            }
        }
    }

    function Graphs(mapping, time) {
        // Set of all can ids
        var canids = [];
        // Set of all fields
        var fields = [];
        // Map from message id to plot data
        var plotData = {};
        // [["graphid", plot, graphdata], ["graphid2", plot, graphdata2]]
        var graphs = [];

        // Populate each graph
        for (var graphid in mapping) {
            if (mapping.hasOwnProperty(graphid)) {
                var labels = mapping[graphid];
                var graphData = [];

                // Link graphData to plotData and override labels.
                for (var i = 0; i < labels.length; i++) {
                    var label = labels[i];
                    canids.push(label[0]);
                    fields.push(label[1]);
                    var msgid = getIdForMsg(label[0], label[1]);
                    if (!(msgid in plotData)) {
                        plotData[msgid] = [];
                    }
                    graphData.push({
                        label: label[2],
                        data: plotData[msgid]
                    });
                }

                // Add graph to  list of graphs
                var graph = [graphid, null, graphData];
                graphs.push(graph);

                // Create plot
                f = makeGraph(graph);
                $(window).resize(f);
                $(window).on('hashchange', f);
                f();
                window.setTimeout(f, 10);
            }
        }
        $.ajax({
            url: '/api/graphs?time=' + time + '&canid=' +
                canids.join('&canid=') + '&field=' + fields.join('&field='),
            type: 'GET',
            dataType: 'json',
            success: onAjax(plotData)
        });

        // Refresh the plots 5 times a second.
        window.setInterval(
            function() {
                for (var i = 0; i < graphs.length; i++) {
                    var graphid = graph[0];
                    var plot = graph[1];
                    var graphData = graph[2];
                    if (plot && $('#' + graphid).is(':visible')) {
                        plot.setData(graphData);
                        plot.setupGrid();
                        plot.draw();
                    }
                }
            }, 200);

        this.plotData = plotData;

        // update takes in a message id and a [time, val] pair.
        var bufferedMillis = parseDuration(time);
        this.update = function(id, point) {
            var series = plotData[id];
            if (series !== undefined && series.length) {
                series.push(point);
                while (series.length > 1 &&
                    point[0] - series[0][0] > bufferedMillis) {
                    series.shift();
                }
            }
        };
    }


    return Graphs;
}();
