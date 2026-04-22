package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/molus/mach/internal/engine"
	"github.com/molus/mach/internal/mta"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "setup":
		fmt.Println("To setup, run: /conductor:setup")
	case "track":
		fmt.Println("To create a track, run: /conductor:newTrack \"description\"")
	case "iterate":
		fmt.Println("To start a Ralph loop (ReAct), run: /ralph:loop \"task\" --completion-promise \"DONE\"")
	case "audit":
		fmt.Println("To scan for security, run: /security:analyze")
	case "checkpoint":
		checkpoint()
	case "revert":
		revert()
	case "jules":
		fmt.Println("To run Jules, run: /jules \"description\"")
	case "calculate":
		calculate()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}
}

func calculate() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: mach calculate <origin> <destination>")
		fmt.Println("Example: mach calculate \"96 St\" \"Chambers St\"")
		return
	}

	originName := os.Args[2]
	destName := os.Args[3]

	originIDs, err := mta.ResolveStation(originName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	destIDs, err := mta.ResolveStation(destName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Origin: %s %v | Destination: %s %v\n", originName, originIDs, destName, destIDs)

	// Fetch live data from all relevant MTA feeds
	urls := getAllFeedURLs()
	client := mta.NewClient(urls...)

	fmt.Println("Fetching live MTA data across all feeds...")
	preds, err := client.GetPredictions(context.Background())
	if err != nil {
		fmt.Printf("Error fetching MTA data: %v\n", err)
		return
	}

	// Calculate wait delta across all lines and station IDs found
	res, err := engine.CalculateWaitDelta(preds, originIDs, destIDs)
	if err != nil {
		fmt.Printf("Decision Error: %v\n", err)
		return
	}

	fmt.Println("\n--- Recommendation ---")
	if len(res.Options) < 2 {
		fmt.Printf("Only one line found: Next %s arrives at %s\n", res.Options[0].Line, res.Options[0].Arrival.Format(time.Kitchen))
		return
	}

	fastest := res.Options[0]
	secondFastest := res.Options[1]

	fmt.Printf("Next %-2s Arrives (Dest): %s\n", fastest.Line, fastest.Arrival.Format(time.Kitchen))
	fmt.Printf("Next %-2s Arrives (Dest): %s\n", secondFastest.Line, secondFastest.Arrival.Format(time.Kitchen))
	fmt.Printf("Wait Delta:             %v\n", res.WaitDelta)

	if res.WaitDelta < 0 {
		fmt.Printf("\n🚀 TAKE the %s train. It arrives %v earlier than the %s!\n", fastest.Line, -res.WaitDelta, secondFastest.Line)
	} else {
		fmt.Printf("\n⚖️ Take either the %s or %s. They arrive at the same time.\n", fastest.Line, secondFastest.Line)
	}

	if len(res.Options) > 2 {
		fmt.Println("\nOther available routes:")
		for i := 2; i < len(res.Options); i++ {
			fmt.Printf("- Line %-2s: %s\n", res.Options[i].Line, res.Options[i].Arrival.Format(time.Kitchen))
		}
	}
}

func getAllFeedURLs() []string {
	return []string{
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-bdfm",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-g",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-jz",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-nqrw",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-l",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-7",
	}
}

func usage() {
	fmt.Println("Multi-Agent Coding Harness (MACH)")
	fmt.Println("Usage: mach <command>")
	fmt.Println("\nCore Commands:")
	fmt.Println("  setup       Initialize a project using Conductor")
	fmt.Println("  track       Create a new track for a feature or bug")
	fmt.Println("  calculate   Compare subway lines between stations")
	fmt.Println("  iterate     Start a Ralph loop (ReAct) for autonomous implementation")
	fmt.Println("  audit       Perform a security scan")
	fmt.Println("\nRegression Prevention:")
	fmt.Println("  checkpoint  Save the current working state (Git commit)")
	fmt.Println("  revert      Roll back to the last checkpoint")
}

func checkpoint() {
	msg := "checkpoint: automated save"
	if len(os.Args) > 2 {
		msg = fmt.Sprintf("checkpoint: %s", os.Args[2])
	}
	run("git", "add", ".")
	run("git", "commit", "-m", msg)
	fmt.Println("✅ Checkpoint created.")
}

func revert() {
	run("git", "reset", "--hard", "HEAD")
	fmt.Println("⏪ Reverted to last checkpoint.")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running %s: %v\n", name, err)
	}
}
