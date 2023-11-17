// Function to fetch column names and populate dropdown options
async function fetchColumns(select) {
  try {
    const response = await axios.get("http://localhost:3000/columns");
    const columns = response.data.columns;

    columns.forEach((column) => {
      const option = document.createElement("option");
      option.value = column;
      option.text = column;
      select.appendChild(option);
    });
  } catch (error) {
    console.error("Error:", error);
  }
}

// Function to add new search parameters
function addParam() {
  const searchParameters = document.getElementById("searchParameters");

  const newSearchParam = document.createElement("div");
  newSearchParam.classList.add("searchParam");

  const columnSelect = document.createElement("select");
  columnSelect.classList.add("column");
  columnSelect.onchange = () => fetchColumns(columnSelect);
  newSearchParam.appendChild(columnSelect);

  const operand = document.createElement("select");
  operand.classList.add("operand");
  ["=", "!=", ">", "<"].forEach((e) => {
    const option = document.createElement("option");
    option.value = e;
    option.text = e;
    operand.appendChild(option);
  });

  newSearchParam.appendChild(operand);

  const valueInput = document.createElement("input");
  valueInput.classList.add("value");
  valueInput.setAttribute("type", "text");
  valueInput.setAttribute("placeholder", "Enter value");
  newSearchParam.appendChild(valueInput);

  const removeButton = document.createElement("button");
  removeButton.textContent = "Remove";
  removeButton.onclick = function () {
    removeParam(this);
  };
  newSearchParam.appendChild(removeButton);

  searchParameters.appendChild(newSearchParam);

  fetchColumns(columnSelect);
}

// Function to remove a search parameter
function removeParam(button) {
  const searchParam = button.parentElement;
  searchParam.remove();
}

// Function to handle the search
function searchLogs() {
  var searchParams = {};
  const criteria = [];
  const searchParamDivs = document.querySelectorAll(".searchParam");
  searchParamDivs.forEach((paramDiv) => {
    const column = paramDiv.querySelector(".column").value;
    const operand = paramDiv.querySelector(".operand").value; // Add operand selection
    const value = paramDiv.querySelector(".value").value;
    const logical = "AND";
    if (column && operand && value) {
      criteria.push({ column, operand, value, logical }); // Include operand in the search parameters
    }
  });

  searchParams = { criteria };
  console.log(searchParams);
  const url = `http://localhost:3000/search`; // Update the search endpoint

  axios
    .post(url, searchParams) // Send a POST request with the searchParams
    .then((response) => {
      // Process and display the search results
      console.log(response.data);
      displayLogs(response.data);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

function displayLogs(logs) {
  const logsDiv = document.getElementById("logs");
  logsDiv.innerHTML = "";

  if (logs.length === 0) {
    logsDiv.innerHTML = "No logs found.";
    return;
  }

  const table = document.createElement("table");
  table.border = "1";

  // Create table header based on the keys in the first log entry
  const thead = document.createElement("thead");
  const headerRow = document.createElement("tr");
  Object.keys(logs[0]).forEach((key) => {
    const th = document.createElement("th");
    th.textContent = key;
    headerRow.appendChild(th);
  });
  thead.appendChild(headerRow);
  table.appendChild(thead);

  // Create table rows with log data
  const tbody = document.createElement("tbody");
  logs.forEach((log) => {
    const row = document.createElement("tr");
    Object.values(log).forEach((value) => {
      const cell = document.createElement("td");
      cell.textContent = value;
      row.appendChild(cell);
    });
    tbody.appendChild(row);
  });
  table.appendChild(tbody);

  logsDiv.appendChild(table);
}

// Add initial search parameter on page load
addParam();
