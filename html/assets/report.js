
// Add an on-load listener for generating the validation report.
document.addEventListener('DOMContentLoaded', function() {
  const jsonDiv = document.getElementById('json');
  const jsonString = jsonDiv.textContent || jsonDiv.innerText || '';
  let json;

  if (!jsonString) {
    document.getElementById('report').innerText('No JSON report data found')
  }

  // Create a JSON object with the report data
  json = JSON.parse(jsonString);

  try {
    // Create a validation report and display it on the webpage
    document.getElementById('report').appendChild(createReport(json));
  } catch (error) {
    document.getElementById('report').innerText = 'JSON Parsing Error: ' + error.message;
  }

  // Add a listener on our download button just to confirm the work is done
  setUpReportDownload(json.time).then(available => console.log("PDF download available"));
});

// Function to set up the save report as a PDF functionality.
async function setUpReportDownload(time) {
  document.getElementById('pdf-dl').addEventListener('click', async function() {
    const reportElement = document.getElementById('report');

    // Get the report element so we can convert it into a PDF and serve it
    if (reportElement) {
      const timestamp = time.replace("T", "_").replaceAll(":", "-").split(".")[0];
      const pdfFileName = 'validation_report_' + timestamp + '.pdf';
      const reportTitle = document.getElementById('report-title')
      const bgColor = window.getComputedStyle(reportTitle).backgroundColor;
      const opt = {
        margin:       [0, 0.5, 0.5, 0.5],
        filename:     pdfFileName,
        image:        { type: 'jpeg', quality: 0.98 },
        html2canvas:  { scale: 2 },
        jsPDF:        { unit: 'in', format: 'letter', orientation: 'landscape' }
      };

      try {
        // Set the report's title background to match the PDF background
        reportTitle.style.backgroundColor = "#FFFFFF";
        await html2pdf().set(opt).from(reportElement).save();
      } catch (error) {
        document.getElementById('report').innerText = 'PDF Generation Error: ' + error.message;
      } finally {
        // Set the report's title background back to the website's color
        reportTitle.style.backgroundColor = bgColor;
      }
    } else {
      document.getElementById('report').innerText = 'Report Generation Error: No report found';
    }
  });
}

// Function to format the timestamp
function formatDateTime(timestamp) {
  return new Date(timestamp).toISOString().slice(0, 19).replace('T', ' ');
}

// Function to create an HTML report from JSON data.
function createReport(data) {
  const div = document.createElement('div');
  const h3 = document.createElement("h3");
  const details = document.createElement("div")
  const table = document.createElement('table');
  const thead = document.createElement('thead');
  const tbody = document.createElement('tbody');
  const headers = ['Header', 'Row', 'Value', 'Message'];
  const headerRow = document.createElement('tr');

  // Make the validation report look pretty
  table.classList.add('table', 'is-bordered', 'is-hoverable', 'is-fullwidth');
  details.classList.add('has-text-right', 'is-size-7')

  headers.forEach(header => {
    const th = document.createElement('th');
    th.textContent = header;
    headerRow.appendChild(th);
  });

  thead.appendChild(headerRow);

  // noinspection JSUnresolvedVariable
  data.warnings.forEach(warning => {
    // Populate table rows with our validation results
    const row = document.createElement('tr');
    row.innerHTML = `
          <td>${warning.header}</td>
          <td>${warning.row + 1}<!-- Row index is 1-based --></td>
          <td>${warning.value}</td>
          <td class="warning">${warning.message}</td>
        `;
    tbody.appendChild(row);
  });

  table.appendChild(thead);
  table.appendChild(tbody);

  // Add some additional markup to make the report pretty
  h3.id = 'report-title';
  h3.classList.add('title', 'is-3');
  h3.innerText = "Validation Report";

  // noinspection JSUnresolvedVariable
  details.innerText = "Profile: " + data.profile + " [" + formatDateTime(data.time) + "]";

  div.appendChild(h3);
  div.appendChild(table);
  div.appendChild(details);

  return div;
}
