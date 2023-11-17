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
var isFirst = true;
// Function to add new search parameters
function addParam() {
  const searchParameters = document.getElementById("searchParameters");
  const newSearchParam = document.createElement("div");
  newSearchParam.classList.add("searchParam");
  if (!isFirst) {
    const logicalSelect = document.createElement("select");
    logicalSelect.classList.add("logical");
    const and = document.createElement("option");
    and.value = "AND";
    and.text = "AND";
    logicalSelect.appendChild(and);
    const or = document.createElement("option");
    or.value = "OR";
    or.text = "OR";
    logicalSelect.appendChild(or);
    newSearchParam.appendChild(logicalSelect);
  }
  const columnSelect = document.createElement("select");
  columnSelect.classList.add("column");
  const defaultOption = document.createElement("option");
  defaultOption.value = "Select an option";
  defaultOption.selected = true;
  defaultOption.disabled = true;
  defaultOption.text = defaultOption.value;
  columnSelect.appendChild(defaultOption);
  newSearchParam.appendChild(columnSelect);

  const operand = document.createElement("select");
  operand.classList.add("operand");
  ["=", "!=", ">", "<", "=~"].forEach((e) => {
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
  if (!isFirst) newSearchParam.appendChild(removeButton);

  isFirst = false;
  searchParameters.appendChild(newSearchParam);
  fetchColumns(columnSelect);

  columnSelect.onchange = () => {
    if (columnSelect.value == "Select an option") {
      operand.disabled = true;
      valueInput.disabled = true;
      document.getElementById("addSearchBtn").disabled = true;
    } else {
      valueInput.disabled = false;
      operand.disabled = false;
      document.getElementById("addSearchBtn").disabled = false;
    }
  };
  columnSelect.dispatchEvent(new Event("change"));
}

// Function to remove a search parameter
function removeParam(button) {
  console.log(button);
  const searchParam = button.parentElement;
  searchParam.remove();
}
function constructSearchParamForTime(startTime, endTime) {
  const timeCriteria = [];

  if (startTime && endTime) {
    const startParam = {
      column: "timestamp",
      operand: ">=",
      value: startTime,
    };
    const endParam = {
      column: "timestamp",
      operand: "<=",
      value: endTime,
    };
    timeCriteria.push(startParam, "AND", endParam);
  } else if (startTime) {
    const startParam = {
      column: "timestamp",
      operand: ">=",
      value: startTime,
    };
    timeCriteria.push(startParam);
  } else if (endTime) {
    const endParam = {
      column: "timestamp",
      operand: "<=",
      value: endTime,
    };
    timeCriteria.push(endParam);
  }

  return timeCriteria;
}

// Function to handle the search
function searchLogs() {
  var searchParams = {};
  const criteria = [];
  const startTime = document.getElementById("startTime").value;
  const endTime = document.getElementById("endTime").value;
  const timeCriteria = constructSearchParamForTime(
    startTime || null,
    endTime || null
  );

  criteria.push(...timeCriteria);

  const searchParamDivs = document.querySelectorAll(".searchParam");
  searchParamDivs.forEach((paramDiv) => {
    const column = paramDiv.querySelector(".column").value;
    const operand = paramDiv.querySelector(".operand").value; // Add operand selection
    const value = paramDiv.querySelector(".value").value;
    const logical = paramDiv.querySelector(".logical") || "AND";

    if (column && operand && value) {
      if (logical) {
        const logicalValue = logical.value || "AND";
        criteria.push(logicalValue);
      }
      criteria.push({ column, operand, value }); // Include operand in the search parameters
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
  var count = logs.count;
  logs = logs.result;
  console.log("logs", logs);
  const logsDiv = document.getElementById("logs");
  logsDiv.innerHTML = "";
  logsDiv.innerText = count + " Logs found";
  if (count == 0) {
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

function updateTimeRange() {
  const timeRange = document.getElementById("timeRange");
  const startTime = document.getElementById("startTime");
  const endTime = document.getElementById("endTime");

  timeRange.addEventListener("change", function () {
    const selectedValue = this.value;
    const now = new Date().toISOString().slice(0, 16); // Current date and time in YYYY-MM-DDTHH:mm format

    switch (selectedValue) {
      case ".08":
        endTime.value = now;
        startTime.value = new Date(Date.now() - 0.08 * 60 * 60 * 1000)
          .toISOString()
          .slice(0, 16);
        break;
      case "1":
        endTime.value = now;
        startTime.value = new Date(Date.now() - 60 * 60 * 1000)
          .toISOString()
          .slice(0, 16); // 1 hour ago
        break;
      case "5":
        endTime.value = now;
        startTime.value = new Date(Date.now() - 5 * 60 * 60 * 1000)
          .toISOString()
          .slice(0, 16); // 5 hours ago
        break;
      case "Forever":
        endTime.value = "";
        startTime.value = "";
      // Add cases for other predefined ranges
      default:
        break;
    }
  });
  timeRange.dispatchEvent(new Event("change"));
}

function handleCustomTimeSelection() {
  const startTime = document.getElementById("startTime");
  const endTime = document.getElementById("endTime");
  const timeRange = document.getElementById("timeRange");

  startTime.addEventListener("input", function () {
    timeRange.value = "custom";
  });

  endTime.addEventListener("input", function () {
    timeRange.value = "custom";
  });
}

// Add initial search parameter on page load
addParam();
// Initialize the functions
updateTimeRange();
handleCustomTimeSelection();
