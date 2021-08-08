# eprocessor

eprocessor is a simple Go-based tool to process a specific and well formatted csv-based file.

* Click to watch the live [demo video](https://youtu.be/vcIizhXkPwg)


## Table of contents
* [Description](#description)
* [Technologies](#technologies)
* [Setup](#setup)
* [Usage](#usage)
* [Outputs](#outputs)
* [Testing](#testing)
* [Upcomings](#upcomings)
* [Contribution](#contribution)
* [License](#license)


## Description

Once launched the program will perform at least these below operations in same order :

1. Download the structured data file from https://s3.amazonaws.com/ecompany/data.csv.
2. Remove the field named 'Memo' from all records.
3. Add a field named "import_date" and populate it appropriately.
4. For any record that has an empty value, set the value of the field to the value "missing".
5. Remove any duplicate records.
6. Submit the records as JSON objects named 'PaymentRecord' to a REST API url with an API key in the 'X-API-KEY' header.

The REST API URL and API KEY are configurable at launching time via positional arguments.  Also the program has been
improved to allow the data source URL to be configurable at lauching time. When adding -save option, provided arguments will
be saved as environnement variables on your system for futher usage without mentionning then again. If successfully set then
*"EPROCESSOR_API_URL"* is for the API URL and *"EPROCESSOR_API_KEY"* is for the API KEY and *"EPROCESSOR_SOURCE_URL"* for the SOURCE URL.
See [Usage](#Usage) section for some practical examples with local dummy backend server and sample data test.

The repository contains a folder named bonus. Inside you will find a dummy api service and a sample data for testing locally the tool.
Once launched, this server expects to receive request for data downloading at *http://localhost:8080/data.csv* and payment records post call
at *http://127.0.0.1:8080/records* with a custom header *(X-API-KEY)* set with "very-long-complex-key" as value. For each payment record received
it will just display that on the console for confirmation purpose.

Finally, at each launch of the eprocessor tool, a dedicated working folder will be created with the name matching the pattern loggingATcurrentdateDOTcurrenttime.
This folder will be used by the program to store the two generated files and the downloaded data file. The first log file generated will be details*.*log
It will contain the program logs such as errors and infos level details. The second log file will be created with the name statistics*.*log
under the name of statistics . log and it will contain all records sent with SUCCESS or FAILURE as prefic according to the API POST call response.

* Click to watch the live [demo video](https://youtu.be/vcIizhXkPwg)


## Technologies

This project is developed with:
* Golang version: 1.13
* Native libraries only


## Setup

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

* You can quickly watch the demo from this youtube link [demo video](https://youtu.be/vcIizhXkPwg)

```Usage:
    
    eprocessor [-source  <download-link-of-the-data>] [-api  <url-of-the-api-service>] [-key  <value-of-the-api-key>] [-save]

Subcommands:
    version    Display the current version of this tool.
    help       Display the help - how to use this tool.


Options:
    -api      Specify the API URL where the payment records will be posted.
    -key      Specify the key to use into the custom HTTP header 'X-API-KEY'.
    -source   Specify the full URL (inc. filename) for download the data.
    -save     If present then provided arguments would be saved as env variables for later use.
    

Arguments:
    url-of-the-api-service     route of the rest api service.
    value-of-the-api-key       value of the X-API-KEY header.
    download-link-of-the-data  url from where to fetch the data.

You have to provide at least the two mandatory arguments values [-api and -key]. In case
you want to launch the tool without any arguments make sure the required parameters are
set as environnement variables ["EPROCESSOR_API_URL" and "EPROCESSOR_API_KEY"] on your system.
In case the source url is not provided or not set as environnement variable ["EPROCESSOR_SOURCE_URL"],
the default link will be used (check the documentation). To have the the parameters set as environnement
variables for the first time, just add -save flag when launching the program. See below third example.


Examples:
	$ eprocessor
    $ eprocessor -api https://ecompany.com/v1/paymentsrecords -key complex-api-key
    $ eprocessor -source https://ecompany.com/data.csv -api https://ecompany.com/v1/paymentsrecords -key complex-api-key
    $ eprocessor -source https://ecompany.com/data.csv -api https://ecompany.com/v1/paymentsrecords -key complex-api-key -save
	
```


## Outputs

* Running the dummy backend server from the source or executable will display:

```

2021/08/02 15:44:59 backend-service up & running at http://127.0.0.1:8080/data.csv to serve sample file.
2021/08/02 15:44:59 backend-service up & running at http://127.0.0.1:8080/records to handle post calls.

```	

* For local testing. Make sure backend service is Up. Then run the eprocessor tool from the command line

```
# Running the eprocessor from the source code with mention of source url and api url and api key options
~$ go run eprocessor.go -source http://localhost:8080/data.csv -api http://127.0.0.1:8080/records -key my-key

# Running the eprocessor from the source code with mention of source url and api url and api key and save flag options
~$ go run eprocessor.go -source http://localhost:8080/data.csv -api http://127.0.0.1:8080/records -key my-key -save

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
=== RUN   TestRemoveMemoField
--- PASS: TestRemoveMemoField (0.00s)
=== RUN   TestReplaceEmptyValues
--- PASS: TestReplaceEmptyValues (0.00s)
=== RUN   TestRemoveDuplicateRecords
--- PASS: TestRemoveDuplicateRecords (0.00s)
=== RUN   TestToJson
--- PASS: TestToJson (0.00s)
=== RUN   ExamplePause_exit
--- PASS: ExamplePause_exit (0.00s)
=== RUN   ExamplePause_continue
--- PASS: ExamplePause_continue (0.00s)
=== RUN   ExampleBanner
--- PASS: ExampleBanner (0.07s)
PASS
ok      github.com/jeamon/eprocessor    0.188s
```


## Upcomings

* add string flag -secret. when specified will be the AES key to encrypt the passed arguments formatted into json then saved into a local file.
* add subcommand load. when specified will require two flags -config (to indicate the local config file) and -secret to precise the key to decrypt.
* when flags -save and -secret are provided - program will only save parameters passed into a local encrypted file and ignore saving into ENV variables.  


## Contribution

Pull requests are welcome. However, I would be glad to be contacted for discussion before.


## License

proprietary - please [contact me](https://blog.cloudmentor-scale.com/contact) before any action.