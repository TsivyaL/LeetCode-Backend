package models

//import "time"

type Question struct {
	ID      string    `json:"id" bson:"_id"`
	Title   string    `json:"title" bson:"title"`
	Body    string    `json:"body" bson:"body"`
	Answer  string    `json:"answer" bson:"answer"`
	
	//Created time.Time `json:"created" bson:"created"`
	
}