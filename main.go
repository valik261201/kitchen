package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var mut sync.Mutex

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func postOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order Order

	_ = json.NewDecoder(r.Body).Decode(&order)
	mut.Lock()
	// add orders to the end of the queue
	orderList.Enqueue(order)

	json.NewEncoder(w).Encode(&order)

	ord, err := PrettyStruct(order)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\nKitchen recieved order:\n", ord)

	go performPostRequest(order)
	mut.Unlock()
}

func performPostRequest(order Order) {
	if !orderList.isEmpty() {
		const myUrl = "http://localhost:3030/distribution"

		// return the first order form the queue
		order := orderList.Dequeue()
		mut.Lock()
		var requestBody, _ = json.Marshal(order)

		fmt.Println("\nOrder with id ", order.Id, " is being cooked.")
		time.Sleep(time.Second + 3)

		fmt.Printf("\nOrder %v was sent to the dining-hall\n", string(requestBody))
		response, err := http.Post(myUrl, "application/json", bytes.NewBuffer(requestBody))
		mut.Unlock()
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		time.Sleep(time.Second + 1)
	}
}

func main() {

	router := mux.NewRouter()

	//URL path and the function to handle
	router.HandleFunc("/order", postOrder).Methods("POST")

	http.ListenAndServe(":8080", router)
}
