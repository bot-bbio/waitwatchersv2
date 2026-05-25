const go = new Go();
let stationNames = [];
let pollInterval = null;
let currentOrigin = "";
let currentDest = "";

// Keep track of active dropdown indices for keyboard navigation
let activeDropdownIndex = -1;

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    document.getElementById("status").innerText = "Decision Engine Active";
    console.log("WaitWatchersV2 Wasm module initialized.");
    
    // Initialize station names
    stationNames = getStationNames().sort();
    
    // Setup autocomplete search inputs
    setupAutocomplete("origin", "origin-dropdown", "clear-origin");
    setupAutocomplete("destination", "destination-dropdown", "clear-dest");
    
    // Setup swap button
    setupSwapButton();
    
    // Setup window click listener to close dropdowns when clicking outside
    document.addEventListener("click", (e) => {
        closeAllDropdowns(e.target);
    });

    // Enable the calculate button once data is loaded and validation passes
    validateInputs();
}).catch((err) => {
    console.error("Failed to load Wasm module:", err);
    document.getElementById("status").innerText = "Error: Decision Engine Offline";
});

// Setup click and submit listeners
document.addEventListener("DOMContentLoaded", () => {
    const calculateBtn = document.getElementById("calculate-btn");
    
    calculateBtn.addEventListener("click", async () => {
        const originVal = document.getElementById("origin").value.trim();
        const destVal = document.getElementById("destination").value.trim();
        
        if (!originVal || !destVal) {
            return;
        }

        // Save active query
        currentOrigin = originVal;
        currentDest = destVal;

        // Clear any existing polling interval
        if (pollInterval) {
            clearInterval(pollInterval);
        }

        // Perform calculation
        await performCalculation(currentOrigin, currentDest, false);

        // Start polling every 15 seconds
        pollInterval = setInterval(async () => {
            console.log("Auto-polling real-time MTA data...");
            await performCalculation(currentOrigin, currentDest, true);
        }, 15000);
    });
});

// Perform Wait Delta Calculation
async function performCalculation(origin, dest, isPoll = false) {
    const status = document.getElementById("status");
    const resultsDiv = document.getElementById("results");
    const loader = document.getElementById("engine-loader");
    const calcBtn = document.getElementById("calculate-btn");
    const btnSpinner = document.getElementById("btn-spinner");

    if (!isPoll) {
        // Show loader and spinners
        loader.style.display = "flex";
        btnSpinner.style.display = "block";
        calcBtn.disabled = true;
        resultsDiv.innerHTML = "";
        resultsDiv.style.display = "none";
        status.innerText = "Connecting to MTA feed...";
    }

    try {
        const res = await calculateWaitDelta(origin, dest);
        console.log("Calculation Result:", res);
        
        // Hide loader
        loader.style.display = "none";
        btnSpinner.style.display = "none";
        calcBtn.disabled = false;
        
        // Render results
        displayResults(res);
        status.innerText = `Last Updated: ${new Date().toLocaleTimeString()}`;
    } catch (err) {
        console.error("Calculation Error:", err);
        loader.style.display = "none";
        btnSpinner.style.display = "none";
        calcBtn.disabled = false;
        
        status.innerText = "Error: Unable to fetch train data";
        
        // Clear poll if error
        if (pollInterval) {
            clearInterval(pollInterval);
            pollInterval = null;
        }

        resultsDiv.style.display = "block";
        resultsDiv.innerHTML = `
            <div class="result-card" style="border-color: var(--accent-red);">
                <div class="action-prompt" style="color: var(--accent-red);">Calculation Failed</div>
                <p style="text-align: center; color: var(--text-muted); font-size: 0.95rem; margin-top: 5px;">
                    ${err.toString().replace("Error:", "")}
                </p>
            </div>
        `;
    }
}

// Render calculation results in a premium UI format
function displayResults(res) {
    const resultsDiv = document.getElementById("results");
    resultsDiv.style.display = "block";
    
    const deltaSecs = Math.abs(res.waitDelta);
    const mins = Math.floor(deltaSecs / 60);
    const secs = Math.floor(deltaSecs % 60);
    const deltaStr = mins > 0 ? `${mins}m ${secs}s` : `${secs}s`;

    let actionPrompt = "";
    let promptClass = "";
    if (res.waitDelta < 0) {
        actionPrompt = `TAKE THE ${res.options[0].line} TRAIN!`;
        promptClass = "prompt-take";
    } else if (res.waitDelta > 0) {
        actionPrompt = `TAKE THE ${res.options[0].line} TRAIN!`;
        promptClass = "prompt-take";
    } else {
        actionPrompt = "TAKE EITHER TRAIN!";
        promptClass = "prompt-either";
    }

    let resultsHTML = `
        <div class="result-card">
            <div class="action-prompt ${promptClass}">${actionPrompt}</div>
            <div class="wait-delta-text">
                ${deltaStr}
                <span>Time Saved</span>
            </div>
    `;

    // Render the options
    if (res.options && res.options.length > 0) {
        const topCount = Math.min(2, res.options.length);
        for (let i = 0; i < topCount; i++) {
            const opt = res.options[i];
            resultsHTML += `
                <div class="line-info">
                    <div class="line-left">
                        <div class="line-badge line-${opt.line}">${opt.line}</div>
                        <div class="arrival-label">Arrives at Destination</div>
                    </div>
                    <div class="arrival-time">
                        ${formatTime(opt.arrival)}
                        ${i === 0 && res.options.length > 1 ? '<span class="time-tag-fastest">Fastest</span>' : ''}
                    </div>
                </div>
            `;
        }

        // Render collapse toggle if more than 2 options exist
        if (res.options.length > 2) {
            resultsHTML += `
                <button id="toggle-all-btn" class="secondary-btn">View All Options (${res.options.length})</button>
                <div id="all-options" style="display: none; flex-direction: column; gap: 15px; margin-top: 5px; border-top: 1px solid rgba(255,255,255,0.06); padding-top: 15px;">
            `;
            for (let i = 2; i < res.options.length; i++) {
                const opt = res.options[i];
                resultsHTML += `
                    <div class="line-info">
                        <div class="line-left">
                            <div class="line-badge line-${opt.line}">${opt.line}</div>
                            <div class="arrival-label">Arrives at Destination</div>
                        </div>
                        <div class="arrival-time">${formatTime(opt.arrival)}</div>
                    </div>
                `;
            }
            resultsHTML += `</div>`;
        }
    }

    resultsHTML += `</div>`;
    resultsDiv.innerHTML = resultsHTML;

    // Expand/collapse options list
    const toggleBtn = document.getElementById("toggle-all-btn");
    if (toggleBtn) {
        toggleBtn.addEventListener("click", () => {
            const allDiv = document.getElementById("all-options");
            if (allDiv.style.display === "none") {
                allDiv.style.display = "flex";
                toggleBtn.innerText = "Hide Other Options";
            } else {
                allDiv.style.display = "none";
                toggleBtn.innerText = `View All Options (${res.options.length})`;
            }
        });
    }
}

// Utility to parse RFC3339 timestamp to readable 12h clock
function formatTime(rfc3339) {
    const date = new Date(rfc3339);
    return date.toLocaleTimeString([], { hour: 'numeric', minute: '2-digit', second: '2-digit' });
}

// Custom autocomplete component
function setupAutocomplete(inputId, dropdownId, clearBtnId) {
    const input = document.getElementById(inputId);
    const dropdown = document.getElementById(dropdownId);
    const clearBtn = document.getElementById(clearBtnId);

    // Watch for input changes to show filter suggestions
    input.addEventListener("input", () => {
        activeDropdownIndex = -1;
        const val = input.value.trim();
        
        if (val) {
            clearBtn.style.display = "block";
            filterStations(val, dropdown, input);
        } else {
            clearBtn.style.display = "none";
            dropdown.style.display = "none";
            validateInputs();
        }
    });

    // Show recommendations when the field gains focus
    input.addEventListener("focus", () => {
        activeDropdownIndex = -1;
        const val = input.value.trim();
        filterStations(val, dropdown, input);
    });

    // Keyboard controls for suggestion navigation
    input.addEventListener("keydown", (e) => {
        const items = dropdown.getElementsByClassName("autocomplete-item");
        if (dropdown.style.display !== "block" || items.length === 0) return;

        if (e.key === "ArrowDown") {
            e.preventDefault();
            activeDropdownIndex = (activeDropdownIndex + 1) % items.length;
            highlightItem(items);
        } else if (e.key === "ArrowUp") {
            e.preventDefault();
            activeDropdownIndex = (activeDropdownIndex - 1 + items.length) % items.length;
            highlightItem(items);
        } else if (e.key === "Enter") {
            e.preventDefault();
            if (activeDropdownIndex > -1 && items[activeDropdownIndex]) {
                items[activeDropdownIndex].click();
            } else if (items.length > 0) {
                items[0].click(); // Default select first
            }
        } else if (e.key === "Escape") {
            dropdown.style.display = "none";
        }
    });

    // Clear input field on button click
    clearBtn.addEventListener("click", () => {
        input.value = "";
        clearBtn.style.display = "none";
        dropdown.style.display = "none";
        input.focus();
        validateInputs();
        
        // Stop current polling if origin/dest is cleared
        if (pollInterval) {
            clearInterval(pollInterval);
            pollInterval = null;
            document.getElementById("status").innerText = "Decision Engine Active";
            document.getElementById("results").style.display = "none";
        }
    });
}

// Match input to station list and construct suggestion markup
function filterStations(val, dropdown, input) {
    dropdown.innerHTML = "";
    
    // Filter matches
    const matches = stationNames.filter(name => 
        name.toLowerCase().includes(val.toLowerCase())
    );

    if (matches.length === 0) {
        dropdown.style.display = "none";
        return;
    }

    matches.forEach(name => {
        const div = document.createElement("div");
        div.classList.add("autocomplete-item");
        
        // Bold match substring
        const index = name.toLowerCase().indexOf(val.toLowerCase());
        if (index > -1 && val.length > 0) {
            const before = name.substring(0, index);
            const match = name.substring(index, index + val.length);
            const after = name.substring(index + val.length);
            div.innerHTML = `${before}<strong>${match}</strong>${after}`;
        } else {
            div.innerText = name;
        }

        // Add selection listener
        div.addEventListener("click", () => {
            input.value = name;
            dropdown.style.display = "none";
            
            // Show clear button
            const clearBtnId = input.id === "origin" ? "clear-origin" : "clear-dest";
            document.getElementById(clearBtnId).style.display = "block";
            
            validateInputs();
        });

        dropdown.appendChild(div);
    });

    dropdown.style.display = "block";
}

// Highlight autocomplete suggestions when navigating using keys
function highlightItem(items) {
    for (let i = 0; i < items.length; i++) {
        items[i].classList.remove("active");
    }
    if (activeDropdownIndex > -1 && items[activeDropdownIndex]) {
        items[activeDropdownIndex].classList.add("active");
        // Scroll list if item is out of view
        items[activeDropdownIndex].scrollIntoView({ block: "nearest" });
    }
}

// Close suggestion dropdown if click occurred elsewhere
function closeAllDropdowns(elmnt) {
    const dropdowns = document.getElementsByClassName("autocomplete-dropdown");
    const origin = document.getElementById("origin");
    const dest = document.getElementById("destination");

    for (let i = 0; i < dropdowns.length; i++) {
        if (elmnt !== dropdowns[i] && elmnt !== origin && elmnt !== dest) {
            dropdowns[i].style.display = "none";
        }
    }
}

// Swap Station Input Values
function setupSwapButton() {
    const swapBtn = document.getElementById("swap-btn");
    
    swapBtn.addEventListener("click", () => {
        const originInput = document.getElementById("origin");
        const destInput = document.getElementById("destination");
        const clearOrigin = document.getElementById("clear-origin");
        const clearDest = document.getElementById("clear-dest");
        
        // Rotation trigger
        swapBtn.classList.add("rotate-swap");
        setTimeout(() => {
            swapBtn.classList.remove("rotate-swap");
        }, 300);

        const temp = originInput.value;
        originInput.value = destInput.value;
        destInput.value = temp;

        // Toggle clear buttons
        clearOrigin.style.display = originInput.value ? "block" : "none";
        clearDest.style.display = destInput.value ? "block" : "none";

        validateInputs();

        // If results are actively showing, recalculate immediately
        if (document.getElementById("results").style.display === "block" || pollInterval) {
            document.getElementById("calculate-btn").click();
        }
    });
}

// Check validation rules and toggle state of primary button
function validateInputs() {
    const originVal = document.getElementById("origin").value.trim();
    const destVal = document.getElementById("destination").value.trim();
    const calculateBtn = document.getElementById("calculate-btn");

    const originValid = stationNames.includes(originVal);
    const destValid = stationNames.includes(destVal);

    if (originValid && destValid && originVal !== destVal) {
        calculateBtn.disabled = false;
        return true;
    } else {
        calculateBtn.disabled = true;
        return false;
    }
}
