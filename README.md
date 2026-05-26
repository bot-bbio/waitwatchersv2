# WaitWatchers

WaitWatchers is a high-performance utility designed for New York City subway riders to solve the decision-making paralysis when choosing between local and express trains. Using real-time MTA GTFS-RT data, it calculates the "Wait Delta"—the precise difference in arrival times at your destination between boarding the next arriving train or waiting for a faster express train—to help you make the fastest travel decision in seconds.

The application operates as both a command-line interface (CLI) and an interactive, client-side WebAssembly (WASM) dashboard.

## Key Features

- **Wait Delta Engine**: Computes exact arrival times at your destination for all available route options.
- **MTA Live Integration**: Direct connection to official GTFS-RT feeds for the 1, 2, 3, 4, 5, 6, 7, S, A, C, E, B, D, F, M, N, Q, R, W, L, J, and Z lines.
- **Logical Transit Hubs**: Groups connected stations (such as Times Sq-42 St and 42 St-Port Authority Bus Terminal) to allow comparison across different platforms and lines.
- **Terminal Station Support**: Robust calculations that fall back to train departures when arrivals are not published (at terminal stations).
- **Glanceable Web UI**: A glassmorphic dark-mode web application featuring custom autocomplete inputs with served line indicators, loading states, and side-by-side comparison tables.

## Tech Stack

- **Backend / Engine**: Go (Golang)
- **Frontend / Client**: WebAssembly (Go WASM), HTML5, and Vanilla CSS
- **Deployment**: Static hosting (GitHub Pages) with a dynamic CORS-bypassing proxy

## Usage

### Command-Line Interface (CLI)

To compile and run the CLI tool, use:

```bash
go build -o mach.exe main.go
./mach.exe calculate "<origin station>" "<destination station>"
```

Example command:
```bash
./mach.exe calculate "42 St-Times Sq / Port Authority" "168 St-Washington Hts"
```

Example output:
```text
Origin: 42 St-Times Sq / Port Authority [127 R16 725 A27] | Destination: 168 St-Washington Hts [112 A09]
Fetching live MTA data across all feeds...

--- Recommendation ---
Next A  Arrives (Dest): 6:15PM
Next 1  Arrives (Dest): 6:17PM
Wait Delta:             -1m23s

TAKE the A train. It arrives 1m23s earlier than the 1!

Other available routes:
- Line C : 6:22PM
```

### Web Application

To run the application locally:

1. Ensure Go is installed on your system.
2. Build the WebAssembly client and start the local development server:
   ```powershell
   powershell -ExecutionPolicy Bypass -File .\run.ps1
   ```
3. Open `http://localhost:8080` in your web browser.

For production, the static web client is deployed to GitHub Pages and uses a public proxy to fetch live MTA feeds securely from the browser.

---
*Developed using the Multi-Agent Coding Harness (MACH) and Conductor Protocol.*
