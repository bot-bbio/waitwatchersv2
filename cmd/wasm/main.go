package main

import (
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/molus/mach/internal/engine"
	"github.com/molus/mach/internal/mta"
)

func calculateWaitDeltaWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return "Error: origin and destination required"
	}

	originName := args[0].String()
	destName := args[1].String()

	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		go func() {
			originIDs, err := mta.ResolveStation(originName)
			if err != nil {
				reject.Invoke(fmt.Sprintf("Origin error: %v", err))
				return
			}

			destIDs, err := mta.ResolveStation(destName)
			if err != nil {
				reject.Invoke(fmt.Sprintf("Destination error: %v", err))
				return
			}

			urls := []string{
				"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs",
				"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace",
				"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-bdfm",
				"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-g",
				"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-jz",
				"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-nqrw",
				"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-l",
				"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-7",
			}
			client := mta.NewClient(urls...)

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			preds, err := client.GetPredictions(ctx)
			if err != nil {
				reject.Invoke(fmt.Sprintf("MTA Fetch error: %v", err))
				return
			}

			res, err := engine.CalculateWaitDelta(preds, originIDs, destIDs)
			if err != nil {
				reject.Invoke(fmt.Sprintf("Calculation error: %v", err))
				return
			}

			// Convert all options to a format JS can consume.
			optionsJS := make([]interface{}, len(res.Options))
			for i, opt := range res.Options {
				optionsJS[i] = map[string]interface{}{
					"line":    opt.Line,
					"arrival": opt.Arrival.Format(time.RFC3339),
				}
			}

			resultMap := map[string]interface{}{
				"options":   js.ValueOf(optionsJS),
				"waitDelta": res.WaitDelta.Seconds(),
			}

			resolve.Invoke(js.ValueOf(resultMap))
		}()

		return nil
	})

	promiseClass := js.Global().Get("Promise")
	return promiseClass.New(handler)
}

func getStationNamesWrapper(this js.Value, args []js.Value) interface{} {
	names := mta.GetStationNames()
	jsNames := make([]interface{}, len(names))
	for i, name := range names {
		jsNames[i] = name
	}
	return js.ValueOf(jsNames)
}

func main() {
	fmt.Println("WaitWatchersV2 WebAssembly Loaded")
	
	// Expose the functions to the global JS scope
	js.Global().Set("calculateWaitDelta", js.FuncOf(calculateWaitDeltaWrapper))
	js.Global().Set("getStationNames", js.FuncOf(getStationNamesWrapper))

	// Keep the Go program running
	select {}
}
