package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Question מייצגת שאלה במערכת עם קלטים ופלטים צפויים
type Question struct {
    ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title            string              `json:"title"`
    Body             string              `json:"body"`
    Inputs           [][]interface{}     `json:"inputs"`  // מערך של מערכים
    ExpectedOutputs  []interface{}       `json:"expected_outputs"`
    FunctionSignature string             `json:"function_signature"`
}
