package main

//User is Interface for user details
type User struct {
	ID      int
	Name    string
	Lname   string
	Country string
}

//ErrorResponse is interface for sending error message
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
