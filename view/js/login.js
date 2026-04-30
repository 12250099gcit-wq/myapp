// view/js/login.js

function login() {
  var _data = {
    email: document.getElementById("email").value,
    password: document.getElementById("pw").value,
  };

  fetch("/login", {
    method: "POST",
    body: JSON.stringify(_data),
    headers: { "Content-type": "application/json; charset=UTF-8" },
  })
    .then((response) => {
      if (response.ok) {
        // Login successful — redirect to the student page
        window.open("student.html", "_self");
      } else {
        // Throw to reach the .catch() handler
        throw new Error(response.status);
      }
    })
    .catch((e) => {
      if (e.message == "401") {
        alert("Unauthorized. Credentials do not match!");
      }
    });
}

function logout() {
  fetch("/logout")
    .then((response) => {
      if (response.ok) {
        window.open("index.html", "_self");
      } else {
        throw new Error(response.statusText);
      }
    })
    .catch((e) => alert(e));
}
