package services

import (
	//"context"
	"errors"
	//"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"Backend/models"
)

// ExecuteAnswer יקבל קוד של תשובה ויבצע אותו בתוך מיכל Docker
// תיקון ב-ExecuteAnswer
func ExecuteAnswer(answer models.Answer) (bool, error) {
	// בודקים את סוג השפה (Python / Go) של הקוד שהמשתמש שלח
	if strings.Contains(answer.Language, "python") {  // השתמש בשדה Language
		result, err := runPythonCode(answer.Code)
		if err != nil {
			return false, err
		}
		return result, nil
	} else if strings.Contains(answer.Language, "go") {  // השתמש בשדה Language
		result, err := runGoCode(answer.Code)
		if err != nil {
			return false, err
		}
		return result, nil
	}

	return false, errors.New("Unsupported language")
}


func runPythonCode(code string) (bool, error) {
	// ניצור קובץ קוד פייתון בתוך ה-volume
	file, err := os.Create("/tmp/answer.py")  // ייווצר בתוך ה-volume המשותף במיכל
	if err != nil {
		log.Printf("Error creating Python file: %s", err)  // לוג שגיאה
		return false, err
	}
	defer file.Close()

	_, err = file.WriteString(code)
	if err != nil {
		log.Printf("Error writing Python code to file: %s", err)  // לוג שגיאה
		return false, err
	}

	// נריץ את קוד הפייתון בתוך מיכל דוקר
	cmd := exec.Command("docker", "run", "--rm", "-v", "answer-volume:/tmp", "python:latest", "python3", "/tmp/answer.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running Python code: %s", err)  // לוג שגיאה
		return false, err
	}

	log.Printf("Python code output: %s", string(output))  // לוג פלט
	if strings.Contains(string(output), "Correct") {
		return true, nil
	}

	return false, nil
}


func runGoCode(code string) (bool, error) {
	// ניצור קובץ קוד Go בזמן ריצה
	file, err := os.Create("/tmp/answer.go")
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = file.WriteString(code)
	if err != nil {
		return false, err
	}

	// נריץ את קוד ה-Go בתוך מיכל דוקר
	cmd := exec.Command("docker", "run", "--rm", "-v", "/tmp:/tmp", "golang:latest", "go", "run", "/tmp/answer.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running Go code: %s", err)
		return false, err
	}

	// נניח שנבדוק אם יש פלט מדויק
	if strings.Contains(string(output), "Correct") {
		return true, nil
	}

	return false, nil
}
