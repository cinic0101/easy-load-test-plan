"use strict";

(function () {
    function uuidv4() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
            var r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }

    function ajaxAddNewProject(id, name, desc, callback) {
        $.ajax({
            method: "POST",
            url: "/service/project",
            data: {
                'id': id,
                'name': name,
                'desc': desc
            },
            success: callback
        });
    }

    $('#newProjectModal .btn-primary').click(function () {
        var id = uuidv4();
        var $modal = $('#newProjectModal');
        var $name = $modal.find("#name");
        var $desc = $modal.find("#description");
        var name = $name.val();
        var desc = $desc.val();

        ajaxAddNewProject(id, name, desc, function (result) {
            if(!result) {
                alert("Server Error");
            }

            $('#projectMenu li a, #mainTab div').removeClass('active');
            $('#projectMenu').prepend('<li class="nav-item"><a class="nav-link list-group-item-action active" data-toggle="list" href="#t-' + id + '" role="tab">' + name + ' <span class="sr-only">(current)</span></a></li>');
            $('#mainTab').append('<div class="tab-pane active" id="t-' + id + '" role="tabpanel">' + name + '</div>');
            $modal.modal('hide');
            $name.val('');
            $desc.val('');
        });
    });

})();