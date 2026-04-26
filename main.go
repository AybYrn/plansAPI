package main

import (
	"encoding/json" // For JSON encoding and decoding - Go structs to JSON and JSON to Go structs.
	"net/http" // For handling HTTP requests and responses - creating a web server and defining endpoints.
	"strconv" // For converting strings to integers and vice versa - used for handling plan IDs in the URL.
	"strings" // For string manipulation - used for parsing the plan ID from the URL path.
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

func getIDFromPath(path string) (int, error) {
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]

	return strconv.Atoi(idStr)
}

func getPlans(w http.ResponseWriter, r *http.Request) { // Handler function to return the list of plans as JSON. r is the incoming HTTP request and w is the response writer used to send the response back to the client.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans) // This converts the plans slice into JSON and writes it to the response.
}

func getPlanByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	for _, plan := range plans { // // use only value, not index, since we don't need the index of the plan in the slice. This loop iterates through the plans slice and checks if any plan has an ID that matches the one extracted from the URL.
		if plan.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(plan)
			return
		}
	}

	http.Error(w, "Plan not found", http.StatusNotFound)
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

func deletePlanByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	for i, plan := range plans { // use index + value since we need the index to remove the plan from the slice. This loop iterates through the plans slice and checks if any plan has an ID that matches the one extracted from the URL.
		if plan.ID == id {
			plans = append(plans[:i], plans[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Plan not found", http.StatusNotFound)
}

func putPlanByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	var updatedPlan Plan
	err = json.NewDecoder(r.Body).Decode(&updatedPlan)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i, plan := range plans {
		if plan.ID == id {
			updatedPlan.ID = id
			plans[i] = updatedPlan
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedPlan)
			return
		}
	}

	http.Error(w, "Plan not found", http.StatusNotFound)
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

func planByIDHandler(w http.ResponseWriter, r *http.Request) { // This function handles requests to the /plans/ endpoint, which is used to get or delete a plan by its ID. It checks the HTTP method and calls the appropriate handler function (getPlanByID for GET requests and deletePlanByID for DELETE requests). If the method is not allowed, it returns a 405 Method Not Allowed error.
	if r.Method == http.MethodGet {
		getPlanByID(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		deletePlanByID(w, r)
		return
	}

	if r.Method == http.MethodPut {
		putPlanByID(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/plans", plansHandler) // sets up routes, This registers handler functions for the /plans endpoint. When a request is made to /plans, the plansHandler function will be called to handle the request.
	http.HandleFunc("/plans/", planByIDHandler) // This registers a handler for the /plans/ endpoint, which is used to get or delete a plan by its ID. The planByIDHandler function will handle requests to this endpoint.
	
	http.ListenAndServe(":8080", nil) // starts server using those routes, This starts the HTTP server on port 8080. The second argument is nil, which means it will use the default HTTP handler (which we have set up with http.HandleFunc).
}
