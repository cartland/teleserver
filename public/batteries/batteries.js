(function($) {

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
