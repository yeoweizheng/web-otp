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
})

function openEditModal(id, accountName, token) {
    $("#editAccountId").val(id);
    $("#editAccountName").val(accountName);
    $("#editAccountToken").val(token);
    $("#editModal").modal("show");
}