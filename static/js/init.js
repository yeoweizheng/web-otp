$(() => {
    $.get("components/header.html", (data) => $("head").html(data))
    $.get("components/topbar.html", (data) => $("#topbar").html(data))
})