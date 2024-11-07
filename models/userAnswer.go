package models

// תיקון ב-model Answer
type Answer struct {
	//ID        string `json:"id" bson:"_id"`
	QuestionID string `json:"question_id" bson:"question_id"`
	//UserID    string `json:"user_id" bson:"user_id"`
	Code       string `json:"code" bson:"code"`        // קוד התשובה שהמשתמש שלח
	Language   string `json:"language" bson:"language"`  // תיקון ל-Language
	IsCorrect  bool   `json:"is_correct" bson:"is_correct"` // אם התשובה נכונה
}
