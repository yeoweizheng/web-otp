function get(url, callback, error) {
    $.ajax({
        type: "GET",
        url: url,
        headers: {
            Authorization: `JWT ${window.localStorage.getItem('token')}`
        },
        success: callback,
        error: error
    })
}

function post(url, data, callback, error) {
    $.ajax({
        type: "POST",
        url: url,
        data: data,
        headers: {
            Authorization: `JWT ${window.localStorage.getItem('token')}`
        },
        success: callback,
        error: error
    })
}