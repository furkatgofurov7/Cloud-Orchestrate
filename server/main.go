package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/ec2", GetEC2Instances).Methods("GET")
	router.HandleFunc("/ec2/{id}/{command}", GetEC2InstanceCommand).Methods("GET")
	router.HandleFunc("/monitor/{id}", MonitorEC2Metrics).Methods("GET")

	router.HandleFunc("/s3", GetS3Buckets).Methods("GET")
	// router.HandleFunc("/s3/create", CreateS3Bucket).Methods("POST")
	// router.HandleFunc("/s3/{id}/upload", UploadS3Object).Methods("POST")
	// router.HandleFunc("/s3/{id}/delete/{key}", DeleteS3Object).Methods("DELETE")

	http.ListenAndServe(":3000", router)
}

func GetEC2Instances(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := listInstances()

	jsonResponse, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", jsonResponse)
}

func GetS3Buckets(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := listBuckets()

	jsonResponse, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", jsonResponse)
}

func GetEC2InstanceCommand(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	instanceID := vars["id"]
	command := vars["command"]

	response := commandInstance(command, instanceID)
	fmt.Fprintf(w, "%s", response)
}

func MonitorEC2Metrics(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	instanceID := vars["id"]
	v := r.URL.Query()
	metricName := v.Get("metric")
	namespace := v.Get("namespace")
	unit := v.Get("unit")

	response := getMetrics(metricName, instanceID, namespace, unit)
	jsonResponse, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", jsonResponse)
}

// func GetS3Buckets(w http.ResponseWriter, r *http.Request)   {}
// func CreateS3Bucket(w http.ResponseWriter, r *http.Request) {}
// func UploadS3Object(w http.ResponseWriter, r *http.Request) {}
// func DeleteS3Object(w http.ResponseWriter, r *http.Request) {}
