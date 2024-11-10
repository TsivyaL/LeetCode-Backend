package models

// Question מייצגת שאלה במערכת עם קלטים ופלטים צפויים
type Question struct {
    ID              string        `json:"id" bson:"_id"`                // מזהה השאלה
    Title           string        `json:"title" bson:"title"`          // כותרת השאלה
    Body            string        `json:"body" bson:"body"`            // גוף השאלה
    Inputs          []interface{} `json:"inputs" bson:"inputs"`        // קלטים שהמשתמש ישתמש בהם (יכולים להיות מסוגים שונים)
    ExpectedOutputs []interface{} `json:"expected_outputs" bson:"expected_outputs"` // פלטים צפויים (יכולים להיות מסוגים שונים)
    FunctionSignature string      `json:"function_signature" bson:"function_signature"` // מיני-תיאור של הפונקציה (כמו שם הפונקציה והפרמטרים שלה)
}
