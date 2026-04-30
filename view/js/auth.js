// view/js/auth.js
// This file is imported FIRST in every protected page.
// document.cookie returns all cookies for the current page as a string.
// An empty string means the user has not logged in.

if (document.cookie === "") {
  alert("User not logged in!!");
  window.open("index.html", "_self"); // _self = same tab
} else {
  console.log("Cookie found — user is logged in");
}
