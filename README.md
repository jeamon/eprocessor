# eprocessor

eprocessor is a simple Go-based tool to process a specific and well formatted csv-based file.


## General info

Once launched the program will perform at least these below operations in same order :

1. Download the structured data file from https://s3.amazonaws.com/ecompany/data.csv.
2. Remove the field named 'Memo' from all records.
3. Add a field named "import_date" and populate it appropriately.
4. For any record that has an empty value, set the value of the field to the value "missing".
5. Remove any duplicate records.
6. Submit the records as JSON objects named 'PaymentRecord' to a REST API url with an API key in the 'X-API-KEY' header.

The REST API URL and API KEY are configurable at lauching time via positional arguments.  Also the program has been
improved to allow the data source URL to be configurable at lauching time. see [Usage](#Usage) section for some examples.

The repository contains a folder named bonus. Inside you will find a dummy api service and a sample data for testing locally the tool.
Once launched, this server expects to receive request for data downloading at http://localhost:8080/data.csv and payment records post call
at http://127.0.0.1:8080/records with a custom header (X-API-KEY) set with "very-long-complex-key" as value. For each payment record received
it will just display that on the console for confirmation purpose.

Finally, at each launch of the eprocessor tool, a logging file will be generated in the format of logging@currentdate.currenttime.log
It will contains the program logs such as errors and infos level details. Also, a second log file will be created with same datetime
under the name of statistics@currentdate.currenttime.log and will contains all records sent with SUCCESS or FAILURE as prefic according
to the API POST call response.


## Technologies

This project is developed with:
* Golang version: 1.13
* Native libraries only


## Installation

On Windows, Linux macOS, and FreeBSD you will be able to download the pre-built binaries once available.
If your system has [Go 1.13+](https://golang.org/dl/), you can pull the codebase and build from the source.

```
# build the eprocessor program
git clone https://github.com/jeamon/eprocessor && cd eprocessor
go build eprocessor.go

# build the dummy backend server
cd bonus
go build dummy-backend-server.go
```


## Usage

```Usage:
    
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
    $ eprocessor -source https://ecompany.com/data.csv -api https://ecompany.com/v1/paymentsrecords -key complex-api-key
```


## Outputs

* Running the dummy backend server from the source or executable will display:

```

2021/08/02 15:44:59 backend-service up & running at http://127.0.0.1:8080/data.csv to serve sample file.
2021/08/02 15:44:59 backend-service up & running at http://127.0.0.1:8080/records to handle post calls.

```	

* For local testing. Make sure backend service is Up. Then run the eprocessor tool from the command line

```
# Running the eprocessor from the source code
~$ go run eprocessor.go -source http://localhost:8080/data.csv -api http://127.0.0.1:8080/records -key my-key

# Running the eprocessor on windows from the executable file (for linux - set the permission on it before)
~$ eprocessor.exe -source http://localhost:8080/data.csv -api http://127.0.0.1:8080/records -key my-key

```


* Below is the output from the eprocessor regarding the above example.


```

///////////////////////////////////////////////////////////////////////////////////
@@@@@@@@@@@@@@@@@@@@ E-COMPANY TOOL // CSV FILE PROCESSOR v1.0 @@@@@@@@@@@@@@@@@@@@
///////////////////////////////////////////////////////////////////////////////////



        [+] downloading the formatted file from the url ... [ SUCCESS ]

        [+] opening csv file from disk for processing ... [ SUCCESS ]

        [+] loading csv all records for processing ... [ SUCCESS ]

        [+] removing of "Memo" field from all records ... [ SUCCESS ]

        [+] replacing all empty values by "missing" ... [ SUCCESS ]

        [+] removing of any duplicate records ... [ SUCCESS ]

        [+] submission of all 273 records to rest api backend ... [ STARTED ]

        [+] please wait ... all records submission progression : 100.00% [273/273]

        [+] Initial Records: 800 / After processed: 273 / sent: 273 / success: 273 / fails: 0 / success rate: 100.00%

                {:} Press [Enter] key to exit


```				
	

## Testing

```
$ go test -v
=== RUN   TestExtractFilename
--- PASS: TestExtractFilename (0.00s)
=== RUN   TestToJson
--- PASS: TestToJson (0.00s)
=== RUN   ExamplePause_exit
--- PASS: ExamplePause_exit (0.00s)
=== RUN   ExamplePause_continue
--- PASS: ExamplePause_continue (0.00s)
=== RUN   ExampleBanner
--- PASS: ExampleBanner (0.08s)
PASS
ok      github.com/jeamon/eprocessor    0.219s
```


## Contribution

Pull requests are welcome. However, I would be glad to be contacted first for discussion before.


## License

proprietary - please [contact me](https://blog.cloudmentor-scale.com/contact) before any action.