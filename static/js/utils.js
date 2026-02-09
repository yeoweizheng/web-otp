let refreshRequest = null;

function rawRequest(method, url, data) {
    return $.ajax({
        type: method,
        url: url,
        data: data,
        contentType: "application/json"
    });
}

function refreshSession() {
    if (!refreshRequest) {
        refreshRequest = rawRequest("POST", "/api/refresh/", JSON.stringify({}));
        refreshRequest.always(() => { refreshRequest = null; });
    }
    return refreshRequest;
}

function request(method, url, data, retried=false) {
    return new Promise((callback, err_callback) => {
        rawRequest(method, url, data).then(
            (resp) => callback(resp),
            (resp) => {
                let isAuthEndpoint = url == "/api/login/" || url == "/api/refresh/";
                if (resp.status == 401 && !retried && !isAuthEndpoint) {
                    refreshSession().then(
                        () => request(method, url, data, true).then(callback).catch(err_callback)
                    ).catch(
                        () => {
                            window.location.href = "/login.html";
                            err_callback(resp);
                        }
                    );
                } else if (resp.status == 401) {
                    window.location.href = "/login.html";
                    err_callback(resp);
                } else {
                    err_callback(resp);
                }
            }
        )
    })
}

function get(url) {
    return request("GET", url);
}

function post(url, data) {
    return request("POST", url, data);
}

function patch(url, data) {
    return request("PATCH", url, data);
}

function del(url) {
    return request("DELETE", url);
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
