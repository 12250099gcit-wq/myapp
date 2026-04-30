// view/js/enroll.js

var selectedRow = null;

// ── Page load: populate dropdowns + display all enrollments ──────
window.onload = function () {
  fetch("/students")
    .then((r) => r.text())
    .then((data) => getStudents(data));

  fetch("/courses")
    .then((r) => r.text())
    .then((data) => getCourses(data));

  fetch("/enrolls")
    .then((r) => r.text())
    .then((data) => getAllEnroll(data));
};

// ── Build the Student ID dropdown ─────────────────────────────────
function getStudents(data) {
  const students = [];
  const allStudents = JSON.parse(data);
  allStudents.forEach((stud) => students.push(stud.stdid));

  var select = document.getElementById("sid");
  for (var i = 0; i < students.length; i++) {
    var option = document.createElement("option");
    option.textContent = students[i];
    option.value = students[i];
    select.appendChild(option);
  }
}

// ── Build the Course ID dropdown ──────────────────────────────────
function getCourses(data) {
  const allCourses = JSON.parse(data);
  var option = "";
  allCourses.forEach((course) => {
    option += '<option value="' + course.cid + '">' + course.cid + "</option>";
  });
  document.getElementById("cid").innerHTML = option;
}

// ── Helper: add one row to the enroll table ───────────────────────
function showTable(enrolled) {
  var table = document.getElementById("myTable");
  var row = table.insertRow(table.length);
  var td = [];
  for (var i = 0; i < table.rows[0].cells.length; i++) {
    td[i] = row.insertCell(i);
  }
  td[0].innerHTML = enrolled.stdid;
  td[1].innerHTML = enrolled.cid;
  // enrolled.date is ISO 8601; split("T")[0] shows only the date part e.g. "2024-05-01"
  td[2].innerHTML = enrolled.date.split("T")[0];
  td[3].innerHTML =
    '<input type="button" onclick="deleteEnroll(this)" value="Delete">';
}

// ── Display one enrollment (after POST) ───────────────────────────
function getEnrolled(data) {
  const enrolled = JSON.parse(data);
  showTable(enrolled);
}

// ── Display all enrollments (on page load) ────────────────────────
function getAllEnroll(data) {
  const allenroll = JSON.parse(data);
  allenroll.forEach((enroll) => showTable(enroll));
}

function resetFields() {
  document.getElementById("sid").value = "";
  document.getElementById("cid").value = "";
}

// ── POST: enroll a student ────────────────────────────────────────
function addEnroll() {
  var _data = {
    stdid: parseInt(document.getElementById("sid").value),
    cid: document.getElementById("cid").value,
  };

  var sid = _data.stdid;
  var cid = _data.cid;

  if (isNaN(sid) || cid === "") {
    alert("Select valid data");
    return;
  }

  fetch("/enroll", {
    method: "POST",
    body: JSON.stringify(_data),
    headers: { "Content-type": "application/json; charset=UTF-8" },
  })
    .then((response) => {
      if (response.ok) {
        fetch("/enroll/" + sid + "/" + cid)
          .then((r) => r.text())
          .then((data) => getEnrolled(data));
      } else {
        throw new Error(response.statusText);
      }
    })
    .catch((e) => {
      if (e.message === "Forbidden") {
        alert(e + ". Duplicate entry!");
      }
    });

  resetFields();
}

// ── DELETE: unenroll a student ────────────────────────────────────
const deleteEnroll = async (r) => {
  if (confirm("Are you sure you want to DELETE this?")) {
    selectedRow = r.parentElement.parentElement;
    var sid = selectedRow.cells[0].innerHTML;
    var cid = selectedRow.cells[1].innerHTML;

    fetch("/enroll/" + sid + "/" + cid, {
      method: "DELETE",
      headers: { "Content-type": "application/json; charset=UTF-8" },
    }).then((response) => {
      if (response.ok) {
        var rowIndex = selectedRow.rowIndex;
        if (rowIndex > 0) {
          document.getElementById("myTable").deleteRow(rowIndex);
        }
      }
    });
  }
};
