$(function () {
    $.ajaxSetup({
        error: function (x, status, error) {
            if (x.status == 406) {
                $.ajax({
                    url: 'http://localhost:8080/refresh',
                    type: 'POST',
                    beforeSend: function(xhr){
                       xhr.withCredentials = true;
                    },
                    success: (message, status, jqxhr) => {
                        if(status === 200) {
                            console.log("refreshed")
                        }
                    }
                });
                console.log("refreshed")
            }
        }
    });
});