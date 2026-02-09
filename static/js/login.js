$("#loginBtn").on("click", (e) => {
    e.preventDefault();
    let username = $("#username").val();
    let password = $("#password").val();
    $.ajax({
        type: "POST",
        url: "/api/login/",
        data: JSON.stringify({username: username, password: password}),
        contentType: "application/json",
        success: (resp) => {
            window.localStorage.setItem("username", resp.username || username);
            window.location.href = "/";
        },
        error: () => {alert("Login failed", "error")}
    })
})
