# WaitWatchersV2

WaitWatchersV2 is a high-performance utility designed for NYC subway riders to solve the "fog of war" when choosing between local and express trains. Using real-time MTA data, it calculates the "Wait Delta"—the precise time difference between taking the next local train or waiting for the express—to help you make the fastest travel decision in seconds.

## 🚀 Key Features

- **Wait Delta Engine:** Real-time calculation of arrival times at your destination for both local and express options.
- **MTA Live Integration:** Direct connection to official GTFS-RT feeds for the 1, 2, 3, 4, 5, 6, 7, and S lines.
- **Decision Clarity:** Clear recommendations (e.g., "WAIT for the express" or "TAKE the local now") with saves calculated in minutes and seconds.
- **High Reliability:** Engineered for accuracy within 2 minutes of real-world arrivals.

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Data Source:** MTA Real-time GTFS-RT APIs
- **Architecture:** Modular engine with a focused CLI interface

## 💻 Usage

To get a travel recommendation, use the `calculate` command:

```bash
# Format: mach calculate "<origin station>" "<destination station>"
./mach calculate "96 St" "72 St"
```

### Example Output:
```text
Origin: 96 St (120) | Destination: 72 St (123)
Fetching live MTA data...

--- Recommendation ---
Next Local Arrives (Dest):   3:25PM
Next Express Arrives (Dest): 3:25PM
Wait Delta:                  22s

🐢 TAKE the local train now. It's faster by 22s.
```

## 🎯 UX Principles

- **Speed (Glanceability):** Designed for decisions in under 2 seconds.
- **Accuracy & Trust:** Real-time data freshness is prioritized to ensure you never miss a train.
- **Utility-Focused:** No fluff. Just the data you need to keep moving.

---
*Developed using the Multi-Agent Coding Harness (MACH) and Conductor Protocol.*
