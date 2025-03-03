
/**
 * Determines which part of the SPA to show depending on the contentID that's passed.
 *
 * @param {string} contentID - The ID of the content-section that's being passed to the function
 */
function showContent(contentID) {
  let content;

  document.querySelectorAll('.content-section').forEach(section => {
      section.style.display = 'none';
  });

  content = document.getElementById(contentID);
  if (content) {
    content.style.display = 'block';

    // If showing API docs, reinitialize Redoc
    if (contentID === "api-docs") {
      Redoc.init("openapi.yml", document.getElementById("api-docs"));
    }
  }
}

/**
 * Adds a listener to toggle the main view to the 'csv-update' page content by default.
 *
 * @param {string} DOMContentLoaded - The type of event listener being added to the page
 */
document.addEventListener("DOMContentLoaded", function () {
  showContent('csv-upload'); // Show CSV Upload by default
});

/**
 * Adds an event listener to override the onClick behavior of the redoc source code copy button.
 *
 * @param {string} click - The type of event listener being added to the page
 */
document.addEventListener("click", /** @param {MouseEvent} event */ function (event) {
  let button = event.target;

  // Find the problematic response samples copy button (really a div) and fix it
  if (button.tagName === "DIV" && button.innerText.trim().startsWith("Copy")) {
    button.querySelectorAll("div").forEach(div => {
      if (!div.innerText.trim() && div.children.length === 0) {
        div.remove(); // Remove the divs that display incorrectly when clicking copy
      } else {
        div.style.backgroundColor = "#fff"; // Give the working ones white backgrounds
      }
    });
  }
});
