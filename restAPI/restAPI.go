package restAPI

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Creates and maintains localhost that listens for GET requests and outputs the specified JSON file stored within the output folder.
func Get() {
	fmt.Println("REST API now operating...")
	fmt.Println("Use GET at: http://localhost:8080/")
	fmt.Println("Example: http://localhost:8080/threads")

	mux := http.NewServeMux()

	// Multiplexer matches requests to this server and can then intake a request
	mux.HandleFunc("GET /{dataObject}", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the dataObject from the request
		dataObject := r.PathValue("dataObject")

		// Sets the header to JSON
		w.Header().Set("Content-Type", "application/json")

		// Retrieve file from output folder
		file, err := os.ReadFile("output/output" + dataObject + ".json")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Output the file
		fmt.Fprintf(w, "%s", file)

	})

	// Listen and serve the server, passing in the multiplexer
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err.Error())
	}
}
