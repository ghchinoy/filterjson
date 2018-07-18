package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = 12001

func main() {
	fmt.Println("Starting filterjson service")

	r := mux.NewRouter()
	r.HandleFunc("/", filterHandler).Methods("POST")

	http.Handle("/", r)
	fmt.Printf("Listening on port %v ...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)

}

func filterHandler(w http.ResponseWriter, r *http.Request) {

	// body
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "unable to read body", http.StatusNotAcceptable)
		return
	}

	// filter
	removes := r.URL.Query()["remove"]

	var testpayload interface{}
	err = json.Unmarshal(data, &testpayload)
	if err != nil {
		log.Println(err)
		http.Error(w, "ugh", http.StatusInternalServerError)
		return
	}

	var payload interface{}
	switch p := testpayload.(type) {
	case map[string]interface{}:
		log.Println("top level JSON object")
		payload = removeKeysFromJSON(p, removes)
	case []interface{}:
		log.Printf("an array: %+v", p)
		payload = removeKeysFromJSONArray(p, removes)
	default:
		log.Printf("idk wtf: %+v", p)
	}

	/*
			payload := make(map[string]interface{})
			err = json.Unmarshal(data, &payload)
			if err != nil {
				log.Println(err)
				http.Error(w, "unable to payload", http.StatusInternalServerError)
				return
		  }
	*/

	/*
			var keys []string
			for k := range payload {
				keys = append(keys, k)
		  }
		  	w.Header().Set("original-keys", fmt.Sprintf("%s", keys))
	*/

	w.Header().Set("removes", fmt.Sprintf("%s", removes))
	w.Header().Set("unfiltered-content-length", fmt.Sprintf("%v", len(data)))

	p, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "unable to byte removed payload", http.StatusInternalServerError)
		return
	}

	w.Write(p)
}

func removeKeysFromJSONArray(payload []interface{}, removes []string) []interface{} {
	var response []interface{}
	for _, v := range payload {
		r := removeKeysFromJSON(v.(map[string]interface{}), removes)
		response = append(response, r)
	}
	return response
}

func removeKeysFromJSON(payload map[string]interface{}, removes []string) map[string]interface{} {
	for _, r := range removes {
		delete(payload, r)
	}
	return payload
}
