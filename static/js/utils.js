function get(url) {
    return new Promise((callback, err_callback) => {
        $.ajax({
            type: "GET",
            url: url,
            headers: {
                Authorization: window.localStorage.getItem('token')
            },
            success: (data) => callback(data),
            error: (error) => err_callback(error)
        })
    })
}

function post(url, data) {
    return new Promise((callback, err_callback) => {
        $.ajax({
            type: "POST",
            url: url,
            data: data,
            headers: {
                Authorization: window.localStorage.getItem('token')
            },
            success: (data) => callback(data),
            error: (error) => err_callback(error)
        })
    })
}

var alertTimeoutIds = [];
function alert(message, type, duration=2000) {
    let elem = $("#alert")
    if (type == "success") {
        elem.removeClass("alert-danger");
        elem.addClass("alert-success");
    } else if (type == "error") {
        elem.removeClass("alert-success");
        elem.addClass("alert-danger");
    }
    elem.html(message);
    elem.fadeIn(200);
    for (id of alertTimeoutIds) { clearTimeout(id); }
    alertTimeoutIds.push(setTimeout(() => {elem.fadeOut(200)}, duration));
}