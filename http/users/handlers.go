package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcctor/marauders/db"
	marauderhttp "github.com/mcctor/marauders/http"
)

const (
	userResultLimit     = 1000
	locationResultLimit = 50
)

func init() {
	marauderhttp.Router.HandleFunc("/users/", allUsers).Methods("GET")
	marauderhttp.Router.HandleFunc("/users/{username}/", specificUser).Methods("GET")
	marauderhttp.Router.HandleFunc("/users/{username}/cloaks/", specificUserCloaks)
	marauderhttp.Router.HandleFunc("/users/{username}/cloaks/{cloak_id}/", specificUserCloak)
	marauderhttp.Router.HandleFunc("/users/{username}/devices/", specificUserDevices)
	marauderhttp.Router.HandleFunc("/users/{username}/devices/{device_id}/", specificUserDevice)
	marauderhttp.Router.HandleFunc("/users/{username}/devices/{device_id}/loc-data/", specificUserDeviceLocData)
	marauderhttp.Router.HandleFunc("/users/{username}/devices/{device_id}/cloaks/", specificUserDeviceCloaks)
}


func allUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "GET" {
		allUsers, err := db.GetUsers(userResultLimit)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(writer).Encode(allUsers)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func specificUser(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "GET" {
		if user, exists := db.GetUser(vars["username"]); exists {
			err := json.NewEncoder(writer).Encode(user)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func specificUserCloaks(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "GET" {
		cloaks, err := db.getCloaksFor(vars["username"])
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(writer).Encode(cloaks)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func specificUserCloak(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "GET" {
		cloak, err := db.GetCloakByID(vars["cloak_id"])
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(writer).Encode(cloak)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func specificUserDevices(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "GET" {
		devices, err := db.getDevicesFor(vars["username"])
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(writer).Encode(devices)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func specificUserDevice(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "GET" {
		deviceID, err := strconv.ParseInt(vars["device_id"], 10, 32)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if device, exists := db.GetDeviceByID(int(deviceID)); exists {
			err = json.NewEncoder(writer).Encode(device)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			writer.WriteHeader(http.StatusNotFound)
		}
	}
}

func specificUserDeviceLocData(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "GET" {
		deviceID, err := strconv.ParseInt(vars["device_id"], 10, 32)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if device, exists := db.GetDeviceByID(int(deviceID)); exists {
			locData, err := device.LocationSnapshots(locationResultLimit)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
			}
			err = json.NewEncoder(writer).Encode(locData)
		} else {
			writer.WriteHeader(http.StatusNotFound)
		}
	}
}

func specificUserDeviceCloaks(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "GET" {
		deviceID, err := strconv.ParseInt(vars["device_id"], 10, 32)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if device, exists := db.GetDeviceByID(int(deviceID)); exists {
			cloaks, err := device.AssociatedCloaks()
			if err != nil {
				log.Println(err)
				writer.WriteHeader(http.StatusInternalServerError)
			}
			err = json.NewEncoder(writer).Encode(cloaks)
		} else {
			writer.WriteHeader(http.StatusNotFound)
		}

	}
}
