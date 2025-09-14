package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Simple HTTP client to test gRPC HTTP gateway
func main() {
	baseURL := "http://localhost:8081"

	// Test 1: Create a user
	fmt.Println("=== Creating User ===")
	createUser := map[string]string{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(createUser)
	resp, err := http.Post(baseURL+"/grpc/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Create User Response: %s\n", body)

	// Parse response to get user ID
	var createResp map[string]interface{}
	json.Unmarshal(body, &createResp)

	var userID string
	if user, ok := createResp["user"].(map[string]interface{}); ok {
		if id, ok := user["id"].(string); ok {
			userID = id
		}
	}

	if userID == "" {
		log.Fatal("Could not get user ID from response")
	}

	// Test 2: Get the created user
	fmt.Println("\n=== Getting User ===")
	resp, err = http.Get(baseURL + "/grpc/users/" + userID)
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Get User Response: %s\n", body)

	// Test 3: List all users
	fmt.Println("\n=== Listing Users ===")
	resp, err = http.Get(baseURL + "/grpc/users")
	if err != nil {
		log.Fatalf("Error listing users: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("List Users Response: %s\n", body)

	fmt.Println("\n=== gRPC HTTP Gateway Test Complete ===")
}
