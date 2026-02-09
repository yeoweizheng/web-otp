$(() => {
    $("#addModal").on("hidden.bs.modal", (e) => { 
        $("#addAccountName").val("");
        $("#addAccountToken").val("");
    })
    $("#addSaveBtn").on("click", async (e) => {
        let accountName = $("#addAccountName").val();
        let token = $("#addAccountToken").val();
        if (!accountName || !token) {
            alert("Account name and token are required", "error");
        } else {
            post("/api/add_account/", JSON.stringify({name: accountName, token: token}))
            .then(async () => {
                $("#addModal").modal("hide");
                alert("Account added", "success");
                await initAccountOTPs();
            }
            ).catch(() => alert("Failed to add account", "error"));
        }
    });
    $("#editModal").on("hidden.bs.modal", (e) => { 
        $("#editAccountToken").val("");
        $("#deleteCheckbox").prop("checked", false);
        $("#editDelBtn").attr("disabled", true);
        $("#showTokenBtn").removeClass("d-none");
        $("#tokenRevealConfirmWrap").addClass("d-none");
        $("#revealPasswordInput").val("");
    })
    $("#editSaveBtn").on("click", async (e) => {
        let accountId = $("#editAccountId").val();
        let accountName = $("#editAccountName").val();
        let token = $("#editAccountToken").val();
        if (!accountName) {
            alert("Account name is required", "error");
        } else {
            let payload = {name: accountName};
            if (token && token.trim()) {
                payload.token = token;
            }
            patch(`/api/update_account/${accountId}/`, JSON.stringify(payload))
            .then(async () => {
                $("#editModal").modal("hide");
                alert("Account updated", "success");
                await initAccountOTPs();
            }
            ).catch(() => alert("Failed to update account", "error"));
        }
    });
    $("#showTokenBtn").on("click", (e) => {
        $("#tokenRevealConfirmWrap").removeClass("d-none");
        $("#revealPasswordInput").trigger("focus");
    });
    $("#confirmShowTokenBtn").on("click", (e) => {
        let accountId = $("#editAccountId").val();
        let password = $("#revealPasswordInput").val();
        if (!password) {
            alert("Password is required", "error");
            return;
        }
        post(`/api/reveal_account_token/${accountId}/`, JSON.stringify({password: password}))
        .then((resp) => {
            $("#editAccountToken").val(resp.token).trigger("input").trigger("focus");
            $("#showTokenBtn").addClass("d-none");
            $("#revealPasswordInput").val("");
            $("#tokenRevealConfirmWrap").addClass("d-none");
            alert("Token loaded", "success");
        })
        .catch(() => {
            alert("Password incorrect or token unavailable", "error");
        });
    });
    $("#cancelShowTokenBtn").on("click", (e) => {
        $("#revealPasswordInput").val("");
        $("#tokenRevealConfirmWrap").addClass("d-none");
    });
    $("#deleteCheckbox").on("change", (e) => {
        if (e.target.checked) {
            $("#editDelBtn").attr("disabled", false);
        } else {
            $("#editDelBtn").attr("disabled", true);
        }
    });
    $("#editDelBtn").on("click", (e) => {
        let accountId = $("#editAccountId").val();
        del(`/api/delete_account/${accountId}/`).then(
            async () => {
                $("#editModal").modal("hide");
                alert("Account deleted", "error");
                await initAccountOTPs();
            }
        ).catch(() => { alert("Failed to delete account", "error")})
    });
})

function openEditModal(id, accountName) {
    $("#editAccountId").val(id);
    $("#editAccountName").val(accountName);
    $("#editAccountToken").val("");
    $("#showTokenBtn").removeClass("d-none");
    $("#revealPasswordInput").val("");
    $("#tokenRevealConfirmWrap").addClass("d-none");
    $("#editModal").modal("show");
}
