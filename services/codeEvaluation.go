package services

import (
	"Backend/models"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

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
	} else if strings.Contains(answer.Language, "c#") {  // השתמש בשדה Language
		result, err := runCSharpCode(answer.Code)
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

func runCSharpCode(code string) (string, error) {
	// נריץ את קוד ה-C# ישירות במיכל של דוקר
	cmd := exec.Command("docker", "run", "--rm", "mcr.microsoft.com/dotnet/sdk:7.0", "bash", "-c", fmt.Sprintf("echo '%s' > /tmp/answer.cs && dotnet new console -o /tmp/csharp-app && cd /tmp/csharp-app && dotnet run", code))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running C# code: %s", err)
		return "false", err
	}

	log.Printf("C# code output: %s", string(output))
	return string(output), nil
}
// func runGoCode(code string) (string, error) {
// 	// הפעלת קוד Go ישירות במיכל של דוקר בלי יצירת קובץ פיזי
// 	cmd := exec.Command("docker", "run", "--rm", "golang:latest", "go", "run", "-e", code)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		log.Printf("Error running Go code: %s", err)
// 		return "false", err
// 	}

// 	log.Printf("Go code output: %s", string(output))
// 	return string(output) ,nil
// }
