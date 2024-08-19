package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"

	//"encoding/hex"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func FileListByDirectory(pathName string) []fs.DirEntry {
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
	return fileList
}

func FileListByCurrentDirectory() []fs.DirEntry {
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
	return fileList
}

func FileListByDirectoryDisplay() {
	fileList := FileListByCurrentDirectory()
	fmt.Println("These files where found: ")
	for _, file := range fileList {
		fmt.Println("  " + file.Name())
	}
}

func FileListByDirectorySearch(fileName string) bool {
	fileList := FileListByCurrentDirectory()
	fileNameMatched := false
	for _, file := range fileList {
		if fileName == file.Name() {
			fileNameMatched = true
		}
	}
	return fileNameMatched
}

func SignatureEncrypt(message string, key string) string {

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))
	signature := h.Sum(nil)
	base64Signature := base64.StdEncoding.EncodeToString(signature)
	//signature := hex.EncodeToString(h.Sum(nil))
	//signature := h.Sum(nil)

	key = "" // Clear the key, tiny bit more secure.
	fmt.Printf("\nOur generated signature: %v", base64Signature)
	return base64Signature

}

func main() {
	FileListByCurrentDirectory()

	// Display the files in the current credential directory.
	apiCredentialFileName := "api-credentials-"
	currentFileName := apiCredentialFileName + "unleashed.txt"
	apiCredentialsFound := FileListByDirectorySearch(currentFileName)

	if apiCredentialsFound == true {
		fmt.Println("We found the API credential file: ", currentFileName)
	} else {
		fmt.Println("No found API credential file called: ", currentFileName)
	}
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// Open the api credentials text file.
	currentFile, err := os.Open(path + "\\" + currentFileName)
	if err != nil {
		fmt.Printf("ERROR: Error opening file.")
		log.Fatal(err)
	}
	defer currentFile.Close()

	// Display text file contents and set the API credentils
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
	// Close the api credentials text file.

	// Display the credentials we grabbed from the file.
	fmt.Printf("\nWe have found credentials. %v, %v\n", unleashedAPIID, unleashedAPIKey)

	// Test API connection to Unleashed. "https://go.unleashedsoftware.com/v2"
	unleashedWebsiteURL := "https://api.unleashedsoftware.com"
	unleashedWebsiteResource := "/Customers?"
	unleashedWebsiteParameters := "CustomerCode=100001" //!!!!!
	unleashedRequestURL := unleashedWebsiteURL + unleashedWebsiteResource + unleashedWebsiteParameters
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// Create the new HTTP GET request.
	request, err := http.NewRequest("GET", unleashedRequestURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	signatureMessage := unleashedWebsiteParameters
	// Adding custom headers for Unleashed API
	// for Content-Type and Accept we can use application/xml or application/json
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("api-auth-id", unleashedAPIID)
	request.Header.Set("api-auth-signature", SignatureEncrypt(signatureMessage, unleashedAPIKey))
	request.Header.Set("client-type", "Sandbox-Billson's Beverages Pty Ltd (Administrators Appointed)/james_tynan_integration_testing")

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	defer response.Body.Close()
	// Check the response status code returned from the request.
	unleashedStatusCodeMap := map[int]string{
		200: "OK: Operation was successful.",
		400: "Bad Request: The request was not in a correct form, or the posted data failed a validation test. Check the error returned to see what was wrong with the request.",
		403: "Forbidden: Method authentication failed.",
		404: "Not Found: Endpoint does not exist (eg using /SalesOrder instead of /SalesOrders).",
		405: "Not Allowed: The method used is not allowed (eg PUT, DELETE) or is missing a required parameter (eg POST requires an /{id} parameter).",
		500: "Internal Server Error: The object passed to the API could not be parsed.",
	}
	responseReturnCode := unleashedStatusCodeMap[response.StatusCode]
	fmt.Printf("\nAPI Reponse Status Code:: %v", responseReturnCode)

	// Read the response body from the request.
	responseBodyContent, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("\nFailed to read the response body: %v", err)
	}
	fmt.Println("Response Status:", response.Status)
	fmt.Println("Response Headers:", response.Header)
	//fmt.Println("Raw response body:", string(responseBodyContent))

	// Unmarshal the JSON response into a go structure
	var responseBodyresult map[string]interface{}
	if err := json.Unmarshal(responseBodyContent, &responseBodyresult); err != nil {
		log.Fatalf("\nFailed to unmarshal the response: %v", err)
	}

	// Display the response from the API request.
	for key, value := range responseBodyresult {
		fmt.Printf("\n %s:%v\n", key, value)
	}

	// request Customer extract json file from Unleashed.
	// display Customer extract from json file.

	// connect to PostgreSQL database.
	// enter new row into customer table.
	// close connection to PostgreSQL.

}
