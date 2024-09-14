$("#loginBtn").on("click", (e) => {
    e.preventDefault();
    let username = $("#username").val();
    let password = $("#password").val();
    $.post("/api/login/", JSON.stringify({username: username, password: password}), (resp) => {
        window.localStorage.setItem("token", resp.token);
        window.location.href = "/";
    }, "json").fail(() => {alert("Login failed", "error")})
})