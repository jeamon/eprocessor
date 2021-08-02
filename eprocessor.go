package main

// This program is a command-line tool built by e-COMPANY for automated payment records data processing.
// The program named `eprocessor` performs into its current version, at least these below actions :
//
// 1. Download the structured data file from https://s3.amazonaws.com/ecompany/data.csv.
// 2. Remove the field named 'Memo' from all records.
// 3. Add a field named "import_date" and populate it appropriately.
// 4. For any empty value, set the value of the field to the value "missing".
// 5. Remove any duplicate records.
// 6. Submit each record as JSON object named 'PaymentRecord' to a REST API with a key in 'X-API-KEY' header.
//
// The Download URL and API URL and API Key are configurable via arguments at launch time. Check help details.
// The API service return valid HTTP status codes with errors. Check the backend service documentation.
//
// For local testing purpose - a dummy backend server was provided into under the name dummy-backend-server.go
// It could be run from source by `go run dummy-backend-server.go` and serve a sample data and handle POST calls.
//
// To contact the author for any feedback or cooperation, use this link https://blog.cloudmentor-scale.com/contact
//
// Version  : 1.0
// Author   : Jerome AMON
// Created  : 01 August 2021

import (
	"bytes"
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// A Record is a final structure of each record after the data file has been proccessed.
type Record struct {
	Date       string `json:"date"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	Zipcode    string `json:"zipcode"`
	Telephone  string `json:"telephone"`
	Mobile     string `json:"mobile"`
	Amount     string `json:"amount"`
	Processor  string `json:"processor"`
	ImportDate string `json:"importdate"`
}

// A PaymentRecord is a structure to be used to build json string before posting to API.
type PaymentRecord struct {
	PaymentRecord Record `json:"PaymentRecord"`
}

// waiting time before program exit at failure.
const waitingTime = 3

// maximum number of pool workers to POST records to API
const maxworkers = 10

// maximum waiting time to establish http connection.
const timeout = 15

// this stores the url to download the data file.
var sourceURL string

// this stores the api uri where to post records.
var apiURL string

// this stores the key to fill into X-API-KEY header.
var apiKEY string

// map console cleaning function based on OS type.
var clear map[string]func()

// custom logger for program INFO only details.
var logInfos *log.Logger

// custom logger for program ERROR only details.
var logError *log.Logger

// custom logger for saving successful sent payment record.
var logSuccessRecords *log.Logger

// custom logger for saving failed to send payment record.
var logFailureRecords *log.Logger

// init is an initializtion function that performs log files creation
// and their associate logger handlers.
func init() {

	// enforce the usage of all available cores on the computer
	runtime.GOMAXPROCS(runtime.NumCPU())

	// initialize the map of functions
	clear = make(map[string]func())
	// add function tp clear linux-based console
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	// add function to clear windows-based console
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// clearConsole is a function that clears the console
// it exits the program if the OS is not supported.
func clearConsole() {
	if clearFunc, ok := clear[runtime.GOOS]; ok {
		clearFunc()
	} else {
		fmt.Println(" [*] Program aborted // failed to clear the console // platform unsupported")
		time.Sleep(waitingTime * time.Second)
		os.Exit(0)
	}
}

// Pause is a function that helps wait until the user press any key.
func Pause(action string) {
	fmt.Printf("\n\t\t{:} Press [Enter] key to %s", action)
	fmt.Scanln()
}

// Banner is a function to display the program title.
func Banner() {
	// first clean the console.
	clearConsole()
	// message to display as program title.
	bannerMsg := " E-COMPANY TOOL // CSV FILE PROCESSOR v1.0 "
	lgMsg := len(bannerMsg)
	// full lengh of the frame - added 20 more characters per each side.
	lgFrame := lgMsg + 40
	// building a centered banner inside the frame.
	fmt.Println("\n" + strings.Repeat("/", lgFrame))
	fmt.Println(strings.Repeat("@", (lgFrame-lgMsg)/2) + bannerMsg + strings.Repeat("@", (lgFrame-lgMsg)/2))
	fmt.Print(strings.Repeat("/", lgFrame), "\n\n\n")
}

// ExtractFilename is a function to retrieve the filename (expected to be the last part) from a url.
func ExtractFilename(srcURL string) string {
	fileUrl, err := url.ParseRequestURI(srcURL)
	// stop the program if failed to parse.
	if err != nil {
		fmt.Print("[ FAILURE ]\n\n\t[-] please check log file for more detailed reason. // ")
		logError.Fatalf("program execution aborted - Errmsg: %v", err)
	}

	path := fileUrl.Path
	// build a slice of each parts from the path.
	parts := strings.Split(path, "/")
	// filename as the last element of the slice.
	filename := parts[len(parts)-1]

	if len(filename) == 0 {
		fmt.Print("[ FAILURE ]\n\n\t[-] please check log file for more detailed reason. // ")
		logError.Fatal("program execution aborted - Errmsg: the link provided does not seems to load a file.")
	}

	return filename
}

// downloadFile is a function that fetches the source data file from the given url
// and save the content into the current directory for further usage by processFile.
func downloadFile() (string, string) {
	fmt.Print("\n\t[+] downloading the formatted file from the url ... ")

	logInfos.Println("extracting the filename from the url.")
	filename := ExtractFilename(sourceURL)
	logInfos.Println("extraction successfully completed.")

	logInfos.Print("downloading the content from the url.")

	// set the http connection timeout.
	client := http.Client{Timeout: timeout * time.Second}

	// get the full file content
	resp, err := client.Get(sourceURL)
	if err != nil {
		fmt.Print("[ FAILURE ]\n\n\t[-] please check log file for more detailed reason. // ")
		logError.Fatalf("failed to download the content - Errmsg: %v", err)
	}
	defer resp.Body.Close()
	logInfos.Println("downloading successfully completed.")

	logInfos.Println("creating destination file for saving.")
	// create an empty destination file.
	dest, err := os.Create(filename)

	if err != nil {
		fmt.Print("[ FAILURE ]\n\n\t[-] please check log file for more detailed reason. // ")
		logError.Fatalf("failed to create destination file - Errmsg: %v", err)
	}
	defer dest.Close()
	logInfos.Println("creation of file successfully completed.")

	logInfos.Println("saving downloaded content to the disk.")
	// flush the content to the file.
	_, err = io.Copy(dest, resp.Body)

	if err != nil {
		fmt.Print("[ FAILURE ]\n\n\t[-] please check log file for more detailed reason. // ")
		logError.Fatalf("failed to save content - Errmsg: %v", err)
	}
	logInfos.Printf("saving of file %s successfully completed.", filename)
	fmt.Println("[ SUCCESS ]")

	// return the name of the file with current date into UTC+0.
	return filename, time.Now().UTC().Format("01/02/2006")
}

// processFile is a function that loads the csv file from disk and performs in order these actions
// 1/ remove "Memo" field. 2/ add "import_date" as new field and fill with current date
// 3/ replace any emply value by "missing". 4/ remove duplicate records. 5/ POST each payment record.
func processFile(filename, importDate string) {

	fmt.Print("\n\t[+] opening csv file from disk for processing ... ")

	logInfos.Println("opening csv file from disk for processing.")
	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Print("[ FAILURE ]\n\n\t[-] please check log file for more detailed reason. // ")
		logError.Fatalf("failed to load the file - Errmsg: %v", err)
	}
	defer csvFile.Close()

	logInfos.Println("opening csv file successfully completed.")
	fmt.Println("[ SUCCESS ]")

	fmt.Print("\n\t[+] loading csv all records for processing ... ")
	logInfos.Println("loading csv all records for processing.")
	reader := csv.NewReader(csvFile)
	// read all entries into slice of slice of string
	allRecords, err := reader.ReadAll()
	if err != nil {
		fmt.Print("[ FAILURE ]\n\n\t[-] please check log file for more detailed reason. // ")
		logError.Fatalf("failed to load the all csv records - Errmsg: %v", err)

	}
	logInfos.Println("loading of records successfully completed.")
	fmt.Println("[ SUCCESS ]")

	// no need to continue if the file does not have any records.
	if len(allRecords) <= 1 {
		logInfos.Println("the downloaded data file seems does not have records entries.")
		fmt.Print("\n\t[+] leaving the program since the there is no records for processing.")
		return
	}

	fmt.Print("\n\t[+] removing of \"Memo\" field from all records ... ")
	logInfos.Println("removing of \"Memo\" field from all records.")
	// retrieve the headers names - first row
	headers := allRecords[0]

	// compute the index of Memo into the slice headers
	memoIndex := -1
	for i, header := range headers {
		if header == "Memo" {
			memoIndex = i
			break
		}
	}

	// remove all Memo value using re-slicing in case we found the field before
	if memoIndex != -1 {
		for i, record := range allRecords {
			// remove the memoIndex and append importation date to the record
			allRecords[i] = append(append(record[:memoIndex], record[memoIndex+1:]...), importDate)
		}
	}
	logInfos.Println("removal of Memo field successfully completed.")
	fmt.Println("[ SUCCESS ]")

	// remove headers record - first element from allRecords slice
	allRecords = allRecords[1:]

	// replace any empty value by "missing" for all records
	fmt.Print("\n\t[+] replacing all empty values by \"missing\" ... ")
	logInfos.Println("replacement of all empty values by \"missing\" started.")
	for _, record := range allRecords {
		for i, v := range record {
			if len(strings.TrimSpace(v)) == 0 {
				record[i] = "missing"
			}
		}
	}
	logInfos.Println("replacement of empty values successfully completed.")
	fmt.Println("[ SUCCESS ]")

	// this following section consists of removing any duplicate records
	// build a Record struct from each record element then add it to
	// the map as key with empty struct as value for memory saving.
	fmt.Print("\n\t[+] removing of any duplicate records ... ")
	logInfos.Println("removal of any duplicate records started.")
	mapOfRecords := make(map[Record]struct{})
	for _, record := range allRecords {
		r := Record{
			Date:       record[0],
			Name:       record[1],
			Address:    record[2],
			Address2:   record[3],
			City:       record[4],
			State:      record[5],
			Zipcode:    record[6],
			Telephone:  record[7],
			Mobile:     record[8],
			Amount:     record[9],
			Processor:  record[10],
			ImportDate: record[11],
		}

		mapOfRecords[r] = struct{}{}
	}

	// compute the initial number of records
	initNumOfRecords := len(allRecords)
	// compute the non-duplicated number of records
	currentNumOfRecords := len(mapOfRecords)

	logInfos.Printf("removal of %d duplicated records successfully completed.\n", (initNumOfRecords - currentNumOfRecords))
	fmt.Println("[ SUCCESS ]")

	// silently clear all records from the slice for memory optimization.
	allRecords = nil

	// compute number of goroutines with a maximum of maxworkers.
	numOfWorkers := int(len(mapOfRecords)/maxworkers) + 1
	if numOfWorkers > maxworkers {
		numOfWorkers = maxworkers
	}

	// posting each record to the API Endpoint as PaymentRecord.
	jobs := make(chan []byte, numOfWorkers)
	// channel to hold each worker success. True when post call succeeds.
	results := make(chan bool)
	// channel to notify end of aggretationResult goroutine.
	done := make(chan bool)
	// number of successfully sent payment records.
	successNum := 0
	// number of failed to send payment records.
	failureNum := 0

	// goroutines to add each json record on the jobs channel for workers.
	go addRecordsAsJobs(jobs, mapOfRecords)
	logInfos.Println("goroutine to jsonify and add records to jobs channel started.")

	// goroutines to monitor results of all workers.
	go aggregateResults(done, results, &successNum, &failureNum, len(mapOfRecords))
	logInfos.Println("goroutine to monitor and compute success rate started.")

	fmt.Printf("\n\t[+] submission of all %d records to rest api backend ... [ STARTED ]\n\t\n", currentNumOfRecords)
	logInfos.Printf("submission of all %d records to rest api backend started.\n", currentNumOfRecords)

	// creating a pool of pre-computed numOfWorkers workers and start them.
	var wg sync.WaitGroup
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go postWorker(&wg, jobs, results)
	}
	// wait for all workers to finish.
	wg.Wait()

	// notify results channel that no more date will come in.
	close(results)
	// block here until we read true from the aggregate goroutine
	<-done

	logInfos.Println("submission of all records successfully completed.")
	fmt.Println()

	// this value could be different from the total records number after the processing
	// in case some payment records failed to be added to the jobs channel at json Marshalling.
	sent := successNum + failureNum

	// success rate is accurate only wi
	successRate := (float64(successNum) / float64(sent)) * 100

	fmt.Printf("\n\t[+] Initial Records: %d / After processed: %d / sent: %d / success: %d / fails: %d / success rate: %.2f%%\n", initNumOfRecords, currentNumOfRecords, sent, successNum, failureNum, successRate)
	// log as INFO the stats into the logging file
	logInfos.Printf("Initial Records: %d / After proccessed: %d / sent: %d / success: %d / fails: %d / success rate: %.2f%%\n", initNumOfRecords, currentNumOfRecords, sent, successNum, failureNum, successRate)
}

// addRecordsAsJobs is function that will be used into a goroutine fashion to
// pick each record from the map and build its associated payment record then
// then marshall it into json and finally add it to the jobs channel for workers.
func addRecordsAsJobs(jobs chan<- []byte, mapOfRecords map[Record]struct{}) {
	for r, _ := range mapOfRecords {
		data, err := json.Marshal(PaymentRecord{PaymentRecord: r})
		if err != nil {
			// unexpected to happen for each record - progression will not reach 100.00% but sucess rate will be accurate
			// track by generating failure id and manually try to build and associated json payment record into stats log.
			sid := generateID()
			logError.Printf("failure to allocate jobs [sid: %s] - Errmsg: %v\n", sid, err)
			// trying to jsonify the record itself
			if d, e := json.Marshal(r); e == nil {
				logFailureRecords.Printf("[sid: %s] {\"PaymentRecord\":%s}", sid, string(d))
			} else {
				// if failed to jsonify the record itslef then manually build
				// the json string like formatted record and log it into stats file.
				s := r.RecordToJson()
				logFailureRecords.Printf("[sid: %s] {\"PaymentRecord\":%s}", sid, s)
			}
			// don't add to jobs channel and move to next record
			continue
		}

		jobs <- data
	}
	close(jobs)
}

// RecordToJson is a function that converts a Record object into json string.
func (r *Record) RecordToJson() string {
	return fmt.Sprintf("{\"date\":%q,\"name\":%q,\"address\":%q,\"address2\":%q,\"city\":%q,\"state\":%q,\"zipcode\":%q,\"telephone\":%q,\"mobile\":%q,\"amount\":%q,\"processor\":%q,\"importdate\":%q}", r.Date, r.Name, r.Address, r.Address2, r.City, r.State, r.Zipcode, r.Telephone, r.Mobile, r.Amount, r.Processor, r.ImportDate)
}

// aggregateResults watchs the results channel and increment the number of success when hits true and
// increment the number of fails when hits false. At the same time, displays real-time progression.
func aggregateResults(done chan<- bool, results <-chan bool, success *int, fails *int, numOfRecords int) {

	total := 0
	// monitor the results channel
	for r := range results {
		// increment the number of post submitted
		total += 1

		if r == true {
			// increment the success numbers
			(*success) = (*success) + 1
		}

		if r == false {
			// increment the failure numbers
			(*fails) = (*fails) + 1
		}
		// enable this below next line to mimic delay into submission progression display
		time.Sleep(time.Duration(100) * time.Millisecond)

		fmt.Printf("\t[+] please wait ... all records submission progression : %2.2f%% [%d/%d]\r", ((float64(total) / float64(numOfRecords)) * 100), total, numOfRecords)
	}

	// send True to the channel once results channel closed
	done <- true
}

// postWorker is a function that will be used as worker in charge of posting payment record
// to the API service and add to the results channel either true or false if success or failure.
func postWorker(wg *sync.WaitGroup, jobs <-chan []byte, results chan<- bool) {
	// loop over the channel of jobs and initiate separate API POST call.
	for job := range jobs {
		// based on status add true or false
		if ok := postPaymentRecord(job); ok {
			results <- true
		} else {
			results <- false
		}
	}
	wg.Done()
}

// generateID uses rand from crypto module to generate random ID into hexadecimal mode this value
// will be used as api call id (cid) and jsonify failure id (sid) for each payment record.
func generateID() string {

	// randomly fill the 8 capacity slice of bytes
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// should not happen but if there - use the current nanosecond time into hexa
		return fmt.Sprintf("%x", time.Now().UTC().UnixNano())
	}
	return fmt.Sprintf("%x", b)
}

// postPaymentRecord is a function to post a payment record to API service.
func postPaymentRecord(jsonBytes []byte) bool {
	// generate an ID for this specific API call. will be used into stats logging.
	cid := generateID()

	// build the http request
	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBytes))
	request.Header.Set("X-API-KEY", apiKEY)
	request.Header.Set("Content-Type", "application/json")

	// set the http connection timeout.
	client := &http.Client{Timeout: timeout * time.Second}
	response, err := client.Do(request)
	if err != nil {
		logError.Printf("failure to submit record - [cid: %s] - Errmsg: %v", cid, err)
		logFailureRecords.Printf("[cid :%s] %s", cid, string(jsonBytes))
		return false
	}
	defer response.Body.Close()

	// check HTTP response header for quick success.
	if response.Status == "200 OK" || response.Status == "201 Created" {
		logInfos.Printf("success to submit record [cid: %s]", cid)
		logSuccessRecords.Printf("[cid: %s] %s", cid, string(jsonBytes))
		return true
	}

	// probable failure on backend side. so we will double check by decoding the response json body.
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return false
	}

	// check response data - status for accurate validation of the failure
	// based on API documentation add or remove successful status codes.
	if result["status"].(float64) == 200 || result["status"].(float64) == 202 {
		log.Printf("success to submit record - [cid: %s]", cid)
		// log the payment record into the stats file with SUCCCESS prefix.
		logSuccessRecords.Printf("%v", string(jsonBytes))
		return true
	} else {
		logError.Printf("failure to create record - [cid: %s] - Errmsg: %s", cid, result["error"].(string))
		// log the payment record into the stats file with FAILURE prefix.
		logFailureRecords.Printf("[cid: %s] %s", cid, string(jsonBytes))
	}

	return false
}

// setupLoggers is a function that create logs files and initialize loggers
func setupLoggers() {

	// get current launch time and build log file name
	startTime := time.Now()
	logTime := fmt.Sprintf("%d%02d%02d.%02d%02d%02d", startTime.Year(), startTime.Month(), startTime.Day(), startTime.Hour(), startTime.Minute(), startTime.Second())

	// create the file to log execution details
	programInfosFile, err := os.OpenFile("logging@"+logTime+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf(" [-] Program aborted. failed to create log file for program execution infos - Errmsg: %v", err)
		time.Sleep(waitingTime * time.Second)
		os.Exit(1)
	}

	// create the file to log post status of all payment records submitted
	recordsStatsFile, err := os.OpenFile("statistics@"+logTime+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf(" [-] Program aborted. failed to create log file to track records submission statistics - Errmsg: %v", err)
		time.Sleep(waitingTime * time.Second)
		os.Exit(1)
	}

	// setup all loggers parameters with microsecnds at timestamp
	logInfos = log.New(programInfosFile, "[ INFOS ] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logError = log.New(programInfosFile, "[ ERROR ] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logSuccessRecords = log.New(recordsStatsFile, "[ SUCCESS ] ", 0)
	logFailureRecords = log.New(recordsStatsFile, "[ FAILURE ] ", 0)
}

func main() {

	// set the default download url - to be used if not provided.
	sourceURL = "https://s3.amazonaws.com/ecompany/data.csv"

	// will be triggered to display help instructions.
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	// configure all flags with globally declared variables
	flag.StringVar(&sourceURL, "source", sourceURL, "Download data file - specify url from where to fetch")
	flag.StringVar(&apiURL, "api", "", "Post payment records - specify the api url where to send")
	flag.StringVar(&apiKEY, "key", "", "Post payment records - specify the api key to be used")

	// check for subcommands : version or help
	if len(os.Args) == 2 {
		if os.Args[1] == "version" || os.Args[1] == "--version" || os.Args[1] == "-v" {
			fmt.Fprintf(os.Stderr, "\n%s\n", version)
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "\n%s\n", usage)
			os.Exit(0)
		}
	}

	// parse the arguments only number of args matches
	// expected number (included the program name itself).
	if len(os.Args) == 5 || len(os.Args) == 7 {
		flag.Parse()
	} else {
		flag.Usage()
		os.Exit(0)
	}

	// -api and -key are mandatory options. stop the program if not provided.
	if apiURL == "" || apiKEY == "" {
		flag.Usage()
		os.Exit(0)
	}

	// display the banner
	Banner()
	// setup all loggers
	setupLoggers()
	// download and save file locally
	filename, importDate := downloadFile()
	// process the downloaded csv file
	processFile(filename, importDate)

	Pause("exit")
}

const version = "current version 1.0 By jeamon@e-company.com"

const usage = `Usage:
    
    eprocessor [-data  <download-link-of-the-data>] [-api  <url-of-the-api-service>] [-key  <value-of-the-api-key>]

Subcommands:
    version    Display the current version of this tool.
    help       Display the help - how to use this tool.


Options:
    -api      Specify the API URL where the payment records will be posted.
    -key      Specify the key to use into the custom HTTP header 'X-API-KEY'.
    -source   Specify the full URL (inc. filename) for download the data.
    

Arguments:
    url-of-the-api-service     route of the rest api service.
    value-of-the-api-key       value of the X-API-KEY header.
    download-link-of-the-data  url from where to fetch the data.

You have to provide at least the two mandatory arguments values [-api and -key]. Upcoming 
version will integrate the capability to launch the tool without any arguments and later be 
prompted to provide at least the two options values (or load them from environnement variables).
In case the source data url link is not provided, it will use the default link. check the 
documentation to view it. Below the two ways to run the current version of this csv processing tool.

Examples:
    $ eprocessor -api https://ecompany.com/v1/paymentsrecords -key complex-api-key
    $ eprocessor -source https://ecompany.com/data.csv -api https://ecompany.com/v1/paymentsrecords -key complex-api-key`
