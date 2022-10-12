package main

type Order struct {
	Id         int   `json:"id"`
	Items      []int `json:"items"`
	Priority   int   `json:"priority"`
	MaxWait    int   `json:"max-wait"`
	PickUpTime int   `json:"pick-up-time"`
}
