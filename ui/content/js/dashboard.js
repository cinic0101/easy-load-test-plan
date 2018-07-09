"use strict";

(function () {
    function uuidv4() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
            var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }

    $('#newProjectModal .btn-primary').click(function () {
        var id = uuidv4();
        $('#projectMenu li a, #mainTab div').removeClass('active');
        $('#projectMenu').prepend('<li class="nav-item"><a class="nav-link list-group-item-action active" data-toggle="list" href="#t-' + id + '" role="tab">Test1 <span class="sr-only">(current)</span></a></li>');
        $('#mainTab').append('<div class="tab-pane active" id="t-' + id + '" role="tabpanel">' + id + '</div>');
        $('#newProjectModal').modal('hide');
    });


})();