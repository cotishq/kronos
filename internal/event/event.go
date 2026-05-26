package event

import "time"

type Event struct {
	ID 			string   	`json:"id"`
	Type		string   	`json:"type"`
	Payload 	string  	`json:"payload"`
	Timestamp 	time.Time	`json:"timestamp"`
}