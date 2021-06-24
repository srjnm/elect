$.ajaxSetup({
    statusCode: {
        406: function () {
            $.ajax({
                url: 'http://localhost:8080/refresh',
                type: 'POST',
                beforeSend: function(xhr){
                    xhr.withCredentials = true;
                },
                success: (data) => {
                    console.log(data)
                }
            });
        }
    }
});