const go = new Go();
let stationNames = [];
let pollInterval = null;

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    document.getElementById("status").innerText = "Decision Engine Ready";
    console.log("WaitWatchersV2 Wasm module initialized.");
    
    // Initialize station names for datalist
    stationNames = getStationNames();
    const datalist = document.getElementById("station-options");
    stationNames.sort().forEach(name => {
        const option = document.createElement("option");
        option.value = name;
        datalist.appendChild(option);
    });
    
    // Enable the calculate button
    document.getElementById("calculate-btn").disabled = false;
}).catch((err) => {
    console.error("Failed to load Wasm module:", err);
    document.getElementById("status").innerText = "Failed to load Decision Engine.";
});

// Logic for button click
document.addEventListener("DOMContentLoaded", () => {
    const btn = document.getElementById("calculate-btn");
    btn.addEventListener("click", async () => {
        const origin = document.getElementById("origin").value;
        const dest = document.getElementById("destination").value;
        
        if (!origin || !dest) {
            document.getElementById("status").innerText = "Please select both stations.";
            return;
        }

        // Clear any existing polling
        if (pollInterval) {
            clearInterval(pollInterval);
        }

        await performCalculation(origin, dest);

        // Start polling every 15 seconds
        pollInterval = setInterval(async () => {
            console.log("Auto-polling MTA data...");
            await performCalculation(origin, dest, true);
        }, 15000);
    });
});

async function performCalculation(origin, dest, isPoll = false) {
    const status = document.getElementById("status");
    const resultsDiv = document.getElementById("results");

    if (!isPoll) {
        status.innerText = "Fetching live MTA data...";
        resultsDiv.innerHTML = "";
        resultsDiv.style.display = "none";
    }

    try {
        const res = await calculateWaitDelta(origin, dest);
        console.log("Calculation Result:", res);
        displayResults(res);
        status.innerText = `Last Updated: ${new Date().toLocaleTimeString()}`;
    } catch (err) {
        console.error(err);
        status.innerText = "Error: " + err;
        if (pollInterval) clearInterval(pollInterval);
    }
}

function displayResults(res) {
    const resultsDiv = document.getElementById("results");
    resultsDiv.style.display = "block";
    
    const deltaSecs = Math.abs(res.waitDelta);
    const mins = Math.floor(deltaSecs / 60);
    const secs = Math.floor(deltaSecs % 60);
    const deltaStr = mins > 0 ? `${mins}m ${secs}s` : `${secs}s`;

    let actionPrompt = "";
    if (res.waitDelta < 0) {
        actionPrompt = `🚀 TAKE THE ${res.options[0].line} TRAIN!`;
    } else if (res.waitDelta > 0) {
        actionPrompt = `🚀 TAKE THE ${res.options[0].line} TRAIN!`;
    } else {
        actionPrompt = "⚖️ TAKE EITHER TRAIN!";
    }

    let resultsHTML = `
        <div class="result-card">
            <div class="action-prompt">${actionPrompt}</div>
            <div class="wait-delta-text">You save ${deltaStr}</div>
    `;

    // Render the top 2 options
    if (res.options && res.options.length > 0) {
        const topCount = Math.min(2, res.options.length);
        for (let i = 0; i < topCount; i++) {
            const opt = res.options[i];
            resultsHTML += `
                <div class="line-info">
                    <div class="line-badge line-${opt.line}">${opt.line}</div>
                    <div class="arrival-time">Arrives at ${formatTime(opt.arrival)} ${i === 0 ? "(Fastest)" : ""}</div>
                </div>
            `;
        }

        // Render "View All" toggle if more than 2 options
        if (res.options.length > 2) {
            resultsHTML += `
                <button id="toggle-all-btn" class="secondary-btn">View All Options (${res.options.length})</button>
                <div id="all-options" style="display: none; margin-top: 20px; border-top: 1px solid var(--mta-light-grey); padding-top: 10px;">
            `;
            for (let i = 2; i < res.options.length; i++) {
                const opt = res.options[i];
                resultsHTML += `
                    <div class="line-info">
                        <div class="line-badge line-${opt.line}">${opt.line}</div>
                        <div class="arrival-time">Arrives at ${formatTime(opt.arrival)}</div>
                    </div>
                `;
            }
            resultsHTML += `</div>`;
        }
    }

    resultsHTML += `</div>`;
    resultsDiv.innerHTML = resultsHTML;

    // Add toggle listener
    const toggleBtn = document.getElementById("toggle-all-btn");
    if (toggleBtn) {
        toggleBtn.addEventListener("click", () => {
            const allDiv = document.getElementById("all-options");
            if (allDiv.style.display === "none") {
                allDiv.style.display = "block";
                toggleBtn.innerText = "Hide Other Options";
            } else {
                allDiv.style.display = "none";
                toggleBtn.innerText = `View All Options (${res.options.length})`;
            }
        });
    }
}

function formatTime(rfc3339) {
    const date = new Date(rfc3339);
    return date.toLocaleTimeString([], { hour: 'numeric', minute: '2-digit' });
}
