package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func renderHome(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "views/index.html")
}

func getUser(res http.ResponseWriter, req *http.Request) {
	var httpError = ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "It's not you it's me.",
	}
	jsonResponse := getUserFromDB()

	if jsonResponse == nil {
		returnErrorResponse(res, req, httpError)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonResponse)
	}
}

func insertUser(res http.ResponseWriter, req *http.Request) {
	var httpError = ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Failed to add user.",
	}
	var userDetails User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userDetails)
	defer req.Body.Close()
	if err != nil {
		returnErrorResponse(res, req, httpError)
	} else {
		httpError.Code = http.StatusBadRequest
		if userDetails.Name == "" {
			httpError.Message = "First Name can't be empty"
			returnErrorResponse(res, req, httpError)
		} else if userDetails.Lname == "" {
			httpError.Message = "Last Name can't be empty"
			returnErrorResponse(res, req, httpError)
		} else if userDetails.Country == "" {
			httpError.Message = "Country can't be empty"
			returnErrorResponse(res, req, httpError)
		} else {
			isInserted := insertUserInDB(userDetails)
			if isInserted {
				getUser(res, req)
			} else {
				returnErrorResponse(res, req, httpError)
			}
		}
	}
}

func deleteUser(res http.ResponseWriter, req *http.Request) {
	var httpError = ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "It's not you it's me.",
	}
	userID := mux.Vars((req))["id"]
	if userID == "" {
		returnErrorResponse(res, req, httpError)
		json.NewEncoder(res).Encode(httpError)
	} else {
		isdeleted := deleteUserFromDB(userID)
		if isdeleted {
			getUser(res, req)
		} else {
			returnErrorResponse(res, req, httpError)
		}
	}
}

func updateUser(res http.ResponseWriter, req *http.Request) {
	var httpError = ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "It's not you it's me.",
	}
	var userDetails User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userDetails)
	defer req.Body.Close()
	if err != nil {
		returnErrorResponse(res, req, httpError)
	} else {
		httpError.Code = http.StatusBadRequest
		if userDetails.Name == "" {
			httpError.Message = "First Name can't be empty"
			returnErrorResponse(res, req, httpError)
		} else if userDetails.ID == 0 {
			httpError.Message = "User Id can't be empty"
			returnErrorResponse(res, req, httpError)
		} else if userDetails.Lname == "" {
			httpError.Message = "Last Name can't be empty"
			returnErrorResponse(res, req, httpError)
		} else if userDetails.Country == "" {
			httpError.Message = "Country Name can't be empty"
			returnErrorResponse(res, req, httpError)
		} else {
			isUpdated := updateUserInDB(userDetails)
			if isUpdated {
				getUser(res, req)
			} else {
				returnErrorResponse(res, req, httpError)
			}
		}
	}
}

func returnErrorResponse(res http.ResponseWriter, _ *http.Request, errorMessage ErrorResponse) {
	httpResponse := &ErrorResponse{
		Code:    errorMessage.Code,
		Message: errorMessage.Message,
	}
	jsonResponse, err := json.Marshal(httpResponse)
	if err != nil {
		panic(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(errorMessage.Code)
	res.Write(jsonResponse)
}
