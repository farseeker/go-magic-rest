package main

// (c) 2018 E-Net Solutions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	http.HandleFunc("/", handleMagic)
	http.ListenAndServe(":8800", nil)
}

func handleMagic(w http.ResponseWriter, r *http.Request) {
	//Build the URL to proxy the request back to
	originalProtocol := r.Header["X-Forwarded-Proto"][0]
	originalHostname := r.Header["X-Forwarded-Host"][0]
	originalAppName := r.Header["X-Forwarded-Appname"][0]
	scriptsPath := r.Header["X-Scripts-Path"][0]

	urlParts := stripEmptyTokens(strings.Split(r.URL.Path, "/"))
	if len(urlParts) != 1 {
		fmt.Fprintf(w, "Invalid number of URL parts, should be exactly one\n")
		fmt.Fprintf(w, "%+v", urlParts)
		return
	}

	arguments := ""
	programName := urlParts[0]

	// These values need to be the first ones for the Magic broker to recognise them
	formValues := url.Values{
		"APPNAME":   []string{originalAppName},
		"PRGNAME":   []string{programName},
		"ARGUMENTS": []string{arguments},
	}

	// Now add in our original query items
	for k, v := range r.URL.Query() {
		formValues[k] = v
	}

	//Formulate our URL path to mgrqispi.dll
	requestPath := fmt.Sprintf("%s://%s/%s", originalProtocol, originalHostname, scriptsPath)
	requestURL, err := url.Parse(requestPath)
	if err != nil {
		fmt.Fprintf(w, "Error in Parse\n")
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Printf("%s\n", requestURL.String())

	// Now that we have our URL we can parse the JSON into Magic parameters

	if r.ContentLength > 0 {
		decoder := json.NewDecoder(r.Body)
		var t map[string]string
		err = decoder.Decode(&t)

		if err != nil {
			fmt.Fprintf(w, "Error in decoding\n")
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Printf("\n%+v\n", t)

		for k, v := range t {
			formValues[k] = []string{v}
		}
	}

	fmt.Printf("\n%+v\n", formValues)

	// Submit our request to the actual Magic engine
	request, err := http.NewRequest("POST", requestURL.String(), strings.NewReader(formValues.Encode()))
	if err != nil {
		fmt.Fprintf(w, "Error in NewRequest\n")
		fmt.Fprintf(w, err.Error())
		return
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(w, "Error in Do\n")
		fmt.Fprintf(w, err.Error())
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(w, "Error in ReadAll\n")
		fmt.Fprintf(w, err.Error())
		return
	}

	w.Write(body)

	return
}

// https://stackoverflow.com/a/46798238/69683
func stripEmptyTokens(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
