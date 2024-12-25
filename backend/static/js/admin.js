// Handle delete user form submission
document.getElementById("deleteForm").addEventListener("submit", async (e) => {
    e.preventDefault();

    const email = document.getElementById("deleteEmail").value;

    try {
        const response = await fetch('/admin/deleteUser', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email: email }),
        });

        if (response.ok) {
            alert('User deleted successfully!');
            location.reload(); // Обновление страницы
        } else {
            const error = await response.text();
            alert('Error deleting user: ' + error);
        }
    } catch (err) {
        console.error(err);
        alert('Failed to delete user.');
    }
});

// Handle update user form submission
document.getElementById("updateForm").addEventListener("submit", async (e) => {
    e.preventDefault();

    const email = document.getElementById("updateEmail").value;

    const newLogin = document.getElementById("newLogin").value;

    try {
        const response = await fetch('/admin/updateUserLogin', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, newLogin }),
        });

        if (response.ok) {
            alert('User login updated successfully!');
            location.reload(); // Обновление страницы
        } else {
            const error = await response.text();
            alert('Error updating login: ' + error);
        }
    } catch (err) {
        console.error(err);
        alert('Failed to update login.');
    }
});

// Handle find user form submission
document.getElementById("findForm").addEventListener("submit", async (e) => {
    e.preventDefault();

    const email = document.getElementById("findEmail").value;
    const resId = document.getElementById("resultId")
    const resEmail = document.getElementById("resultEmail")
    const resLogin = document.getElementById("resultLogin")


        const response = await fetch('/admin/findUser', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email: email }),
        });

        if (response.ok) {
            const data = await response.json(); // Await the JSON response
            resId.textContent = "ID: " + data["id"]
            resEmail.textContent = "email: " +  data["email"]
            resLogin.textContent = "login: " +  data["login"]
            // console.log(data["login"])
        } else {
            const error = await response.json();
            alert('Error finding login: ' + error["error"]);
        }
    
});

document.getElementById("sendInfo").addEventListener("submit", async (e) => {
    e.preventDefault();
    const msg = document.getElementById("message").value;

    const response = await fetch('/admin/sendMessage', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ msg: msg }),
    });

    if (response.ok) {
        alert("Mail sent")
    } else {
        const error = await response.json();
        alert('Error finding login: ' + error["error"]);
    }
})