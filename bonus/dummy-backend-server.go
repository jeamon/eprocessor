/*
 This is a dummy RestFul API & Basic File Server Backend which could be used as backend test for the csv-processor program
 It runs on localhost at port 8080. Fill free to change it into the code and align the new vurls on the client side as well.
 It expects to receive this request http://127.0.0.1:8080/data.csv in order to send back the content of sample_data_test.csv 
 Regarding the API, it expects to receive json data which matches PaymentRecord structure and will 200 or 202 status code.
 In case it failed to receive proper json data or failed to process the payload, it will reply to client with a json message 
 which follows ApiResponse structure while setting appropriate http status error code. 
*/

package main

import (

	"io"
	"os"
	"log"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

//== Pour envoyer la liste des questions en Json
type ApiResponse struct {
	Status int `json:"status"`
	Error string `json:"error"`
}

// format for each record after processed
type Record struct {
	Date string `json:"date"`
	Name string `json:"name"`
	Address string `json:"address"`
	Address2 string `json:"address2"`
	City string `json:"city"`
	State string `json:"state"`
	Zipcode string `json:"zipcode"`
	Telephone string `json:"telephone"`
	Mobile string `json:"mobile"`
	Amount string `json:"amount"`
	Processor string `json:"processor"`
	ImportDate string `json:"importdate"`
}

// struct to build json data
type PaymentRecord struct {
	PaymentRecord Record `json:"PaymentRecord"`
}

// sample key for verification
const API_KEY = "my-key"

// createPaymentRecord is a function that handles /records POST requests and emulate the record creation by printing on console. 
func createPaymentRecord(w http.ResponseWriter, r *http.Request) {

	// by default return only json data
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	// handle only http post method
	if r.Method == "POST" {

		if key := r.Header.Get("X-API-KEY"); key != API_KEY {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(ApiResponse{Status:401, Error:"unauthorized access. bad key provided."})
			return
		}

		// read the payload with into a safety manner by limiting
		reqBody, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(ApiResponse{Status:500, Error:"failed to parse request payload."})
			return
		}

    	var paymentRecord PaymentRecord 
    	if err := json.Unmarshal(reqBody, &paymentRecord); err != nil {
    		w.WriteHeader(422) // object non traitable
    		json.NewEncoder(w).Encode(ApiResponse{Status:422, Error:"data submitted is not expected format."})
			return
    	}

    	// mimic creation of the payment record by displaying on screen
    	log.Printf("successfully received new record - %v\n", paymentRecord)

		// could also use w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ApiResponse{Status:200, Error:""})
		return
	}

	// received http request method other than POST
	w.WriteHeader(400)	
	json.NewEncoder(w).Encode(ApiResponse{Status:400, Error:"unsupported request method. check documentation."})
	return
}

// sendDataFile is a function that handles /data.csv requests and send back the content of file sample_data_test.csv 
func sendDataFile(w http.ResponseWriter, r *http.Request) {

	filename := "sample_data_test.csv"
	// Open the file for reading. Assuming the existence of sample data.
	f, err := os.Open(filename)
    if err != nil {
    	// stop the dummy server
    	log.Fatal(err)
    }
    defer f.Close()

	// Read file content into memory
    fileBytes, err := ioutil.ReadAll(f)
    if err != nil {
    	// stop the dummy server
        log.Fatal(err)
    }

    // add some basics custom headers
    r.Header.Add("Content-Type", "txt/csv")
    r.Header.Add("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
    r.Header.Add("Content-Description", "File Transfer")
    r.Header.Add("Content-Disposition", "attachment; filename=data.csv")
    r.Header.Add("Expires", "0")
    r.Header.Add("Content-Length", strconv.Itoa(len(fileBytes)))
    http.ServeFile(w, r, filename)
}


func main() {
	// setup the route to create to create record with its handler
	http.HandleFunc("/data.csv", sendDataFile)
	// setup the route to create to create record with its handler
	http.HandleFunc("/records", createPaymentRecord)
	// spin up the dummy server on localhost at port 8080
	log.Println("backend-service up & running at http://127.0.0.1:8080/data.csv to serve sample file.")
	log.Println("backend-service up & running at http://127.0.0.1:8080/records to handle post calls.")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}