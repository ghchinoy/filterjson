package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const port = 12001

func main() {
	fmt.Println("Starting filterjson service")

	r := mux.NewRouter()
	r.HandleFunc("/", timedHandler(filterHandler)).
		Methods("POST") //.
		//Queries("remove", "filter")

	fmt.Printf("Listening on port %v ...\n", port)

	//http.Handle("/", r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), handler))
}

func timedHandler(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer timethis(time.Now(), fmt.Sprintf("%+v", r.URL.Query()))
		h.ServeHTTP(w, r)
	})
}

func timethis(start time.Time, name string) {
	log.Printf("%s took %s", name, time.Since(start))
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func filterHandler(w http.ResponseWriter, r *http.Request) {

	// body
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "unable to read body", http.StatusNotAcceptable)
		return
	}

	// filters
	removes := r.URL.Query()["remove"]
	filters := r.URL.Query()["filter"]

	var testpayload interface{}
	err = json.Unmarshal(data, &testpayload)
	if err != nil {
		log.Println(err)
		http.Error(w, "ugh", http.StatusInternalServerError)
		return
	}

	var payload interface{}
	switch p := testpayload.(type) {
	case map[string]interface{}: // JSON Object
		if len(filters) > 0 {
			log.Println("filtering")
			payload = returnOnlyKeysInJSON(p, filters)
			p = payload.(map[string]interface{})
		}
		if len(removes) > 0 {
			log.Println("removing")
			payload = removeKeysFromJSON(p, removes)
		}
	case []interface{}: // JSON Array
		if len(filters) > 0 {
			payload = returnOnlyKeysInJSONArray(p, filters)
			p = payload.([]interface{})
		}
		if len(removes) > 0 {
			payload = removeKeysFromJSONArray(p, removes)
		}
	default: // unknown
		log.Printf("idk wtf: %+v", p)
	}

	w.Header().Set("filters", fmt.Sprintf("%s", filters))
	w.Header().Set("removes", fmt.Sprintf("%s", removes))
	w.Header().Set("unfiltered-content-length", fmt.Sprintf("%v", len(data)))

	p, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "unable to byte removed payload", http.StatusInternalServerError)
		return
	}

	w.Write(p)
}

// removeKeysFromJSONArray Removes the list of keys in "removes" from the array
// of JSON objects
func removeKeysFromJSONArray(payload []interface{}, removes []string) []interface{} {
	var response []interface{}
	for _, v := range payload {
		r := removeKeysFromJSON(v.(map[string]interface{}), removes)
		response = append(response, r)
	}
	return response
}

// removeKeysFromJSON Removes the keys in "removes" from the JSON object given
func removeKeysFromJSON(payload map[string]interface{}, removes []string) map[string]interface{} {
	for _, r := range removes {
		delete(payload, r)
	}
	return payload
}

// returnOnlyKeysInJSON returns only the keys listed in "filters" in the given JSON object
func returnOnlyKeysInJSON(payload map[string]interface{}, filters []string) map[string]interface{} {
	kept := make(map[string]interface{})
	for _, f := range filters {
		if _, ok := payload[f]; ok == true {
			kept[f] = payload[f]
		}
	}
	return kept
}

// returnOnlyKeysInJSONArray returns an array of JSON objects with only the keys specified
// in "filters"
func returnOnlyKeysInJSONArray(payload []interface{}, filters []string) []interface{} {
	var response []interface{}
	for _, v := range payload {
		// this will error if the type conversion doesn't work
		// valid json is array of arrays, not everything is an array of objects,
		// and an array of arrays will panic here
		kept := returnOnlyKeysInJSON(v.(map[string]interface{}), filters)
		response = append(response, kept)
	}
	return response
}
