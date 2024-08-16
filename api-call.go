package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	// this is our current directory.
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// look at the files currently in the credential directory.
	fmt.Println("Looking for files in current directory: ", path+"\\")
	fileList, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	// display the files in the current credential directory.
	apiCredentialFileName := "api-credentials-"
	currentFileName := apiCredentialFileName + "unleashed.txt"
	apiCredentialsFound := false
	fmt.Println("These files where found: ")
	for _, file := range fileList {
		fmt.Println("  " + file.Name())
		if file.Name() == currentFileName {
			apiCredentialsFound = true
		}
	}
	if apiCredentialsFound == true {
		fmt.Println("We found the API credential file: ", currentFileName)
	} else {
		fmt.Println("No found API credential file called: ", currentFileName)
	}

	// open the api credentials text file.
	currentFile, err := os.Open(path + "\\" + currentFileName)
	if err != nil {
		fmt.Printf("ERROR: Error opening file.")
		log.Fatal(err)
	}
	defer currentFile.Close()

	// print the contents of the text file.
	/*
		scanner := bufio.NewScanner(openFile)
		for scanner.Scan() {
			_ = scanner.Text()
			fmt.Println("", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error: Error opening file:  %v", err)
			return
		}
	*/

	// load the api credentials in.
	scanner := bufio.NewScanner(currentFile)
	var unleashedAPIID string
	var unleashedAPIKey string
	fmt.Printf("We have found these credentials:\n")
	for scanner.Scan() {
		_ = scanner.Text()
		currentLine := scanner.Text()
		currentLinePrefix, currentLineSuffix, currentLineDividerFound := strings.Cut(currentLine, ":")
		if currentLineDividerFound == true {
			fmt.Printf("%v: %v\n", currentLinePrefix, currentLineSuffix)
			if currentLinePrefix == "ID" {
				unleashedAPIID = currentLineSuffix
			} else if currentLinePrefix == "Key" {
				unleashedAPIKey = currentLineSuffix
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: Error opening file:  %v", err)
		return
	}
	currentFile.Close()
	// close the api credentials text file.

	// display the credentials.
	// performed in the credential extraction process.

	// test API connection to Unleashed.
	unleashedWebsite := "https://go.unleashedsoftware.com/v2"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// Create request
	request, err := http.NewRequest("GET", unleashedWebsite, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer response.Body.Close()

	unleashedStatusCodeMap := map[int]string{
		200: "OK: Operation was successful.",
		400: "Bad Request: The request was not in a correct form, or the posted data failed a validation test. Check the error returned to see what was wrong with the request.",
		403: "Forbidden: Method authentication failed.",
		404: "Not Found: Endpoint does not exist (eg using /SalesOrder instead of /SalesOrders).",
		405: "Not Allowed: The method used is not allowed (eg PUT, DELETE) or is missing a required parameter (eg POST requires an /{id} parameter).",
		500: "Internal Server Error: The object passed to the API could not be parsed.",
	}
	responseReturnCode := unleashedStatusCodeMap[200]
	responseReturnCode = unleashedStatusCodeMap[response.StatusCode]
	fmt.Printf("\nAPI Reponse Status Code:: %v", responseReturnCode)
	// display reponse from Unleashed.

	// request Customer extract json file from Unleashed.
	// display Customer extract from json file.

	// connect to PostgreSQL database.
	// enter new row into customer table.
	// close connection to PostgreSQL.
	fmt.Printf("\n%v, %v\n", unleashedAPIID, unleashedAPIKey)
}
