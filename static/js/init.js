$(() => {
    $.get("components/header.html", (data) => $("head").html(data))
})