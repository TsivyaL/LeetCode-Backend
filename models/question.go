package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
    StatusNotStarted = "Not Started"
    StatusInProgress = "In Progress"
    StatusCompleted  = "Completed"
)


// Question represents a question in the system with inputs and expected outputs
type Question struct {
    ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title            string              `json:"title"`
    Body             string              `json:"body"`
    Inputs           [][]interface{}     `json:"inputs"`  // An array of arrays
    ExpectedOutputs  []interface{}       `json:"expected_outputs"`
    FunctionSignature string             `json:"function_signature"`
    Status           string              `json:"status" bson:"status"`
    CodeInProgress  string             `json:"code_in_progress,omitempty"`
}
