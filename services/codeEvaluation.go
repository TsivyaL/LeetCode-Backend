package services

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"strings"
	"Backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ExecuteAnswer יקבל קוד של תשובה ויבצע אותו בתוך מיכל Docker
// תיקון ב-ExecuteAnswer
func ExecuteAnswer(answer models.Answer) (bool, error) {
	// בודקים את סוג השפה (Python / Go) של הקוד שהמשתמש שלח
	question, err := GetQuestionByID(answer.QuestionID)
	if err != nil {
		return false, err
	}

	if strings.Contains(answer.Language, "python") {  // השתמש בשדה Language
		result, err := runPythonCode(answer.Code)
		if err != nil {
			return false, err
		}
		return strings.Contains(result,question.Answer), nil
	} else if strings.Contains(answer.Language, "go") {  // השתמש בשדה Language
		result, err := runGoCode(answer.Code)
		if err != nil {
			return false, err
		}
		return strings.Contains(result,question.Answer), nil
	}

	return false, errors.New("Unsupported language")
}
func GetQuestionByID(id string) (models.Question, error) {
    var question models.Question
    err := QuestionsCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&question)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return question, errors.New("question not found")
        }
        return question, err
    }
    return question, nil
}

func runPythonCode(code string) (string, error) {
	// נריץ את קוד הפייתון ישירות במיכל של דוקר בלי קובץ מקומי
	cmd := exec.Command("docker", "run", "--rm", "python:latest", "python3", "-c", code)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running Python code: %s", err)
		return "false", err
	}

	log.Printf("Python code output: %s", string(output))
	return string(output) ,nil

	
}


func runGoCode(code string) (string, error) {
	// הפעלת קוד Go ישירות במיכל של דוקר בלי יצירת קובץ פיזי
	cmd := exec.Command("docker", "run", "--rm", "golang:latest", "go", "run", "-e", code)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running Go code: %s", err)
		return "false", err
	}

	log.Printf("Go code output: %s", string(output))
	return string(output) ,nil
}
