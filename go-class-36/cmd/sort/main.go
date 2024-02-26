package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gorilla/mux"
)

// STEP 1: must build program as binary first since there is animation -> go build ./cmd/sort
// STEP 2: run `./sort`, which will run the program sort that was build
// (adding `-speed faster` or `-speed fastest` args will run fast versions)
// STEP 3: View animations in browser
// - http://localhost:8081/insert?loop=0 for insertion sort
// - http://localhost:8081/qsort?loop=0 for quick sort
// - http://localhost:8081/qsortf?loop=0 for quick sort fast
// - etc.
// STEP 4: http://localhost:8081/debug/pprof/ for profiling,
// but better to do with script so run `./scripts/profile.sh` (if permission error, run `chmod +x ./scripts/profile.sh` first)
// STEP 5: in a separate terminal run `./scripts/sort.sh`
// STEP 6: when profiling is done, rename output with following command: `mv profile profile-slow`
// (or `mv profile profile-faster` or `mv profile profile-fastest` if running fast versions)
// STEP 7: run `go tool pprof -http=":6060" sort profile-slow ` for visualization
// here (http://localhost:6060/ui/?f=) we can see various graphs such as flame graphs, top, etc.

func main() {
	var speed string

	flag.StringVar(&speed, "speed", "slow", "painting speed")
	flag.Parse()

	switch speed {
	case "faster":
		paintSquare = paintSquareFast
	case "fastest":
		paintSquare = paintSquareFastest
	case "slow":
		// nop
	default:
		log.Fatal("unknown speed:", speed)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	log.Printf("Speed %q", speed)
	log.Printf("Listening on port %s", port)

	router := mux.NewRouter()

	router.Use(logMiddleware)
	router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	router.HandleFunc("/insert", insertHandler).Methods(http.MethodGet)
	router.HandleFunc("/qsort", qsortHigh).Methods(http.MethodGet)
	router.HandleFunc("/qsortm", qsortMiddle).Methods(http.MethodGet)
	router.HandleFunc("/qsort3", qsortMedian).Methods(http.MethodGet)
	router.HandleFunc("/qsorti", qsortInsert).Methods(http.MethodGet)
	router.HandleFunc("/qsortf", qsortFlag).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
