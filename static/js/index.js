var currTableData = [];
var nextTableData = [];
var timestamp = null;
var nextUpdateTimestamp = null;

$(async () => {
    $("#searchInput").on("input", () => { updateTable(); });
    $("#logoutLink").on("click", () => {
        window.localStorage.clear();
        window.location.href = "/login.html";
    });
    await initAccountOTPs();
})

async function initAccountOTPs() {
    timestamp = new Date().valueOf();
    updateProgress();
    nextUpdateTimestamp = timestamp - (timestamp % 30000) + 30000;
    currTableData = await get("/api/account_otps/");
    updateTable();
    nextTableData = await get(`/api/account_otps/?timestamp=${Math.round(nextUpdateTimestamp/1000)}`);
    setInterval(refreshAccountOTPs, 100);
}

async function refreshAccountOTPs() {
    timestamp = new Date().valueOf();
    updateProgress();
    if (timestamp >= nextUpdateTimestamp) {
        currTableData = nextTableData;
        nextUpdateTimestamp = timestamp - (timestamp % 30000) + 30000;
        updateTable();
        nextTableData = await get(`/api/account_otps/?timestamp=${Math.round(nextUpdateTimestamp/1000)}`);
    }
}

function updateTable() {
    let html = "";
    if (!currTableData) {
        $("#tableData").html("");
        updateEvents();
        return;
    }
    currTableData.sort((a, b) => a.name.localeCompare(b.name))
    let searchText = $("#searchInput").val().toLowerCase()
    for (let row of currTableData) {
        if (!row.name.toLowerCase().includes(searchText)) continue;
        html += `
        <tr data-id=${row.id}>
            <td class="fs-6">${row.name}</td>
            <td class=""><button type="button" class="btn bgn-lg btn-link fs-6 px-2 otp-btn">${row.otp}</button></td>
            <td>
                <button type="button" class="btn btn-lg btn-link px-2 edit-btn" data-mdb-ripple-init><i class="fas fa-pen-to-square text-dark"></i></button>
                <button type="button" class="btn btn-lg btn-link px-2 delete-btn" data-mdb-ripple-init><i class="fas fa-trash-can text-danger"></i></button>
            </td>
        `;
    }
    $("#tableData").html(html);
    updateEvents();
}

function updateProgress() {
    let progress = (timestamp % 30000) / 30000 * 100;
    $("#progressBar").attr('aria-valuenow', progress).css('width', `${progress}%`);
}

function updateEvents() {
    $(".otp-btn").on("click", (e) => {
        let id = e.target.closest("tr").dataset.id;
        let otp;
        for (let row of currTableData) {
            if (row.id == id) {
                otp = row.otp;
                break
            }
        }
        navigator.clipboard.writeText(otp).then(() => {alert("Copied to clipboard", "success")})
    });
    $(".edit-btn").on("click", (e) => {
        let id = e.target.closest("tr").dataset.id;
        for (let row of currTableData) {
            if (row.id == id) {
                openEditModal(id, row.name, row.token);
                break
            }
        }
    })
    $(".delete-btn").on("click", (e) => {
        let id = e.target.closest("tr").dataset.id;
        del(`/api/delete_account/${id}/`).then(
            async () => {
                alert("Account deleted", "error");
                await initAccountOTPs();
            }
        ).catch(() => { alert("Failed to delete account", "error")})
    })
}
