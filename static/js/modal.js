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
        $("#deleteCheckbox").prop("checked", false);
    })
    $("#editSaveBtn").on("click", async (e) => {
        let accountId = $("#editAccountId").val();
        let accountName = $("#editAccountName").val();
        let token = $("#editAccountToken").val();
        if (!accountName || !token) {
            alert("Account name and token are required", "error");
        } else {
            patch(`/api/update_account/${accountId}/`, JSON.stringify({name: accountName, token: token}))
            .then(async () => {
                $("#editModal").modal("hide");
                alert("Account updated", "success");
                await initAccountOTPs();
            }
            ).catch(() => alert("Failed to update account", "error"));
        }
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

function openEditModal(id, accountName, token) {
    $("#editAccountId").val(id);
    $("#editAccountName").val(accountName);
    $("#editAccountToken").val(token);
    $("#editModal").modal("show");
}