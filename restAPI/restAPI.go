package restAPI

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// List of GET endpoints:

/*
http://localhost:8080/{dataObject}
http://localhost:8080/threads
http://localhost:8080/matches
http://localhost:8080/rankings

http://localhost:8080/Ranking/{region}
http://localhost:8080/Ranking/Asia-Pacific
http://localhost:8080/Ranking/Europe
http://localhost:8080/Ranking/North-America
http://localhost:8080/Ranking/Brazil
http://localhost:8080/Ranking/Korea
http://localhost:8080/Ranking/Japan
etc...
*/

// Creates and maintains localhost that listens for GET requests and outputs the specified JSON file stored within the output folder.
func Get() {
	fmt.Println("REST API now operating...")
	fmt.Println("Use GET at: http://localhost:8080/")
	fmt.Println("Endpoint: http://localhost:8080/{dataObject}")
	fmt.Println("Example: http://localhost:8080/threads")
	fmt.Println("Endpoint: http://localhost:8080/Ranking/{region}")
	fmt.Println("Example: http://localhost:8080/Ranking/Asia-Pacific")

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

	mux.HandleFunc("GET /Ranking/{region}", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the dataObject from the request
		region := r.PathValue("region")

		// Sets the header to JSON
		w.Header().Set("Content-Type", "application/json")

		// Retrieve file from output folder
		file, err := os.ReadFile("output/ranking/output" + region + "Rankings.json")
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
