function get(url) {
    return new Promise((callback, err_callback) => {
        let token = window.localStorage.getItem('token');
        if (!token) {
            window.location.href = "/login.html";
            return;
        }
        $.ajax({
            type: "GET",
            url: url,
            headers: {
                Authorization: token
            },
            success: (data) => callback(data),
            error: (resp) => {
                if (resp.status == 401) { window.location.href = "/login.html";
                } else { err_callback(resp); }
            }
        })
    })
}

function post(url, data) {
    return new Promise((callback, err_callback) => {
        let token = window.localStorage.getItem('token');
        if (!token) {
            window.location.href = "/login.html";
            return;
        }
        $.ajax({
            type: "POST",
            url: url,
            data: data,
            headers: {
                Authorization: token
            },
            success: (data) => callback(data),
            error: (resp) => {
                if (resp.status == 401) { window.location.href = "/login.html";
                } else { err_callback(resp); }
            }
        })
    })
}

function patch(url, data) {
    return new Promise((callback, err_callback) => {
        let token = window.localStorage.getItem('token');
        if (!token) {
            window.location.href = "/login.html";
            return;
        }
        $.ajax({
            type: "PATCH",
            url: url,
            data: data,
            headers: {
                Authorization: token
            },
            success: (data) => callback(data),
            error: (resp) => {
                if (resp.status == 401) { window.location.href = "/login.html";
                } else { err_callback(resp); }
            }
        })
    })
}

function del(url) {
    return new Promise((callback, err_callback) => {
        let token = window.localStorage.getItem('token');
        if (!token) {
            window.location.href = "/login.html";
            return;
        }
        $.ajax({
            type: "DELETE",
            url: url,
            headers: {
                Authorization: token
            },
            success: (data) => callback(data),
            error: (resp) => {
                if (resp.status == 401) { window.location.href = "/login.html"; } 
                else if (resp.status == 204) { callback(resp); } 
                else { err_callback(resp); }

            }
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