package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

// RootObject is the rootobject the returned data
type RootObject struct {
	Results []Result `json:"results"`
	Count   int64    `json:"count"`
	CurPage int64    `json:"cur_page"`
	PerPage int64    `json:"per_page"`
}

// Result is each individual consultant
type Result struct {
	FormID                    string   `json:"form_id"`
	Surname                   string   `json:"step1_surname"`
	GivenName                 string   `json:"step1_given_name"`
	ConsultantID              string   `json:"form_consultant_id"`
	MembershipStatus          string   `json:"form_membership_status"`
	DateStatusChanged         string   `json:"form_date_status_changed"`
	Country                   string   `json:"country"`
	State                     string   `json:"state"`
	City                      string   `json:"city"`
	Postal                    string   `json:"postal"`
	Address                   string   `json:"address"`
	StatusShortName           string   `json:"status_short_name"`
	AgentsSurname             string   `json:"agents_surname"`
	AgentsName                string   `json:"agents_name"`
	Companies                 []string `json:"companies"`
	MembershipStatusText      string   `json:"form_membership_status_text"`
	MembershipStatusShortText string   `json:"form_membership_status_short_text"`
	ReasonText                string   `json:"form_reason_text"`
}

//GetHeaders returns the field names as a string slice
func (r *Result) GetHeaders() []string {
	var headers []string
	val := reflect.Indirect(reflect.ValueOf(r))

	for i := 0; i < val.NumField(); i++ {

		x := val.Type().Field(i).Name
		headers = append(headers, x)
	}
	return headers

}

func main() {

	var data RootObject
	var list []Result

	for c := 'A'; c <= 'Z'; c++ {

		extjsonData := getJSONData(string(c), string(0))

		json.Unmarshal(extjsonData, &data)
		pagesBeforeRounding := float64(data.Count / 25)
		rounded := int(math.Ceil(pagesBeforeRounding))
		start := 0
		for pages := 0; pages < rounded; pages++ {

			jsonData := getJSONData(string(c), string(start))

			json.Unmarshal(jsonData, &data)
			for _, result := range data.Results {
				list = append(list, result)
			}

			start += 25

		}

	}

	writeCSVFile(list)

}

func getJSONData(letter string, start string) []byte {

	resp, err := http.PostForm("https://secure.iccrc-crcic.ca/search/do/lang/en",
		url.Values{"form_fields": {""}, "query": {""}, "letter": {letter}, "start": {start}})

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	jsonData, err := ioutil.ReadAll(resp.Body)

	return jsonData

}

func writeCSVFile(list []Result) {
	//Get Current Date and format according
	t := time.Now().Format("2006_01_02")

	// Create a csv file and name it with the date included

	p := filepath.FromSlash("./")

	f, err := os.Create(p + "ListOfConsultants" + t + ".csv")

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	//These next four lines write the headers to the first line of the CSV
	//The GetStructHeaders gives me the headers as a string slice which I then write as a single line
	var a Result
	headers := a.GetHeaders()
	w := csv.NewWriter(f)
	w.Write(headers)

	//Now that we've written headers, iterate through the array of records and write each line
	for _, obj := range list {
		var record []string
		record = append(record, obj.FormID)
		record = append(record, obj.Surname)
		record = append(record, obj.GivenName)
		record = append(record, obj.ConsultantID)
		record = append(record, obj.MembershipStatus)
		record = append(record, obj.DateStatusChanged)
		record = append(record, obj.Country)
		record = append(record, obj.State)
		record = append(record, obj.City)
		record = append(record, obj.Postal)
		record = append(record, obj.Address)
		record = append(record, obj.StatusShortName)
		record = append(record, obj.AgentsSurname)
		record = append(record, obj.AgentsName)
		record = append(record, "null") //Need to figure this out later
		record = append(record, obj.MembershipStatusText)
		record = append(record, obj.MembershipStatusShortText)
		record = append(record, obj.ReasonText)
		w.Write(record)
	}

	w.Flush()
	println("The CSV file is ready - check the directory you ran this program in")
}

// GetStructHeaders gets the field names for Result struct and returns them as a string array. Not using this right now
func GetStructHeaders(a Result) []string {

	var headers []string
	val := reflect.Indirect(reflect.ValueOf(a))

	for i := 0; i < val.NumField(); i++ {

		x := val.Type().Field(i).Name
		headers = append(headers, x)
	}
	return headers

}
