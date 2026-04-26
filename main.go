package main

import (
	"encoding/json" // For JSON encoding and decoding - Go structs to JSON and JSON to Go structs.
	"net/http" // For handling HTTP requests and responses - creating a web server and defining endpoints.
)

type Plan struct { // Plan struct defines the structure of a subscription plan with fields for ID, Title, Description, and Price.
	ID          int     `json:"id"` // json Tag: When converting this struct to JSON, use the name id for the ID field.
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var plans = []Plan{ // in-memory list of plans. In a real application, this would likely be stored in a database.
	{ID: 1, Title: "Basic Plan", Description: "Basic features", Price: 9.99},
	{ID: 2, Title: "Premium Plan", Description: "Premium features", Price: 29.99},
}

func getPlans(w http.ResponseWriter, r *http.Request) { // Handler function to return the list of plans as JSON. r is the incoming HTTP request and w is the response writer used to send the response back to the client.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans) // This converts the plans slice into JSON and writes it to the response.
}

func createPlan(w http.ResponseWriter, r *http.Request) { // Handler function to create a new plan.
	var newPlan Plan

	err := json.NewDecoder(r.Body).Decode(&newPlan) // This reads the JSON body from the request and converts it into a Go Plan. The &newPlan means we pass the memory address, so Go can fill the variable with the data from the JSON.
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newPlan.ID = len(plans) + 1
	plans = append(plans, newPlan)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // This sets the HTTP status code to 201 Created, indicating that a new resource has been successfully created.
	json.NewEncoder(w).Encode(newPlan) // This converts the newPlan struct into JSON and writes it to the response, along with a 201 Created status code.
}

func plansHandler(w http.ResponseWriter, r *http.Request) { // This function handles incoming HTTP requests to the /plans endpoint. It checks the HTTP method of the request and calls the appropriate handler function (getPlans for GET requests and createPlan for POST requests). If the method is not allowed, it returns a 405 Method Not Allowed error.
	if r.Method == http.MethodGet {
		getPlans(w, r)
		return
	}

	if r.Method == http.MethodPost {
		createPlan(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/plans", plansHandler)

	http.ListenAndServe(":8080", nil)
}
