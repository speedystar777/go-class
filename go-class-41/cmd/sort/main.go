package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// MUST BE SET by go build -ldflags "-X main.version=999"
// like 0.6.14-0-g26fe727 or 0.6.14-2-g9118702-dirty

// STEP 1: run `make sort` to build program binary
// STEP 2: run program binary: `./sort`
// visit one of the urls (/insert, /qsort, /qsortm, etc.) to see sorting
// STEP 3: `make docker` will create docker image (a template loaded onto the container to run it, like a set of instructions.)
// STEP 4: `docker run -p 8081:8081 sort-anim` will produce a docker container (a self-contained, runnable software application or service)
// -p arg is port and is mapping local-machine-port:port-in-dockerfile
// otherwise it will run on random port; once again, we can visit the listed urls to see program run

var version string // do not remove or modify

func main() {
	port := os.Getenv("PORT")
	router := mux.NewRouter()

	if port == "" {
		port = "8081"
	}

	router.Use(logMiddleware)
	router.HandleFunc("/insert", insertHandler).Methods(http.MethodGet)
	router.HandleFunc("/qsort", qsortHigh).Methods(http.MethodGet)
	router.HandleFunc("/qsortm", qsortMiddle).Methods(http.MethodGet)
	router.HandleFunc("/qsort3", qsortMedian).Methods(http.MethodGet)
	router.HandleFunc("/qsorti", qsortInsert).Methods(http.MethodGet)
	router.HandleFunc("/qsortf", qsortFlag).Methods(http.MethodGet)
	router.HandleFunc("/version", showVersion).Methods(http.MethodGet)

	log.Printf("version %s listening on port %s", version, port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
