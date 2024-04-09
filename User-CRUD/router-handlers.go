package main

import (
	"encoding/json"
	"net/http"
)

func renderHome(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "views/index.html")
}

func getUsers(res http.ResponseWriter, req *http.Request) {
	var httpError = ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "It's not you it's me.",
	}
	jsonResponse := getUserFromDB()

	if jsonResponse == nil {
		return ErrorResponse(httpError)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonResponse)
	}
}

func insertUser(res http.ResponseWriter, req *http.Request) {
	var httpErrorr = ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Failed to add user.",
	}
	var userDetails User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userDetails)
	defer req.Body.Close()
	if err != nil {
		return ErrorResponse(res, req, httpErrorr)
	} else {
		httpErrorr.Code = http.StatusBadRequest
		if userDetails.Name == "" {
			httpErrorr.Message = "First Name can't be empty"
			return ErrorResponse(res, req, httpErrorr)
		} else if userDetails.Lname == "" {
			httpErrorr.Message = "Last Name can't be empty"
			return ErrorResponse(res, req, httpErrorr)
		} else if userDetails.Country == "" {
			httpErrorr.Message = "Country can't be empty"
			return ErrorResponse(res, req, httpErrorr)
		} else {
			isInserted := insertUserInDB(userDetails)
			if isInserted {
				getUsers(res, req)
			} else {
				return ErrorResponse(res, req, httpErrorr)
			}
		}
	}
}
