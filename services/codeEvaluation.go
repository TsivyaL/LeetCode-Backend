package services

import (
	"Backend/models"
	//"context"
	//"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
)

// ExecuteAnswer יקבל קוד של תשובה ויבצע אותו בתוך מיכל Docker
func ExecuteAnswer(answer models.Answer) (bool, error) {
	// בודקים את סוג השפה (Python / JS) של הקוד שהמשתמש שלח
	question, err := FetchQuestionByID(answer.QuestionID)
	if err != nil {
		return false, err
	}

	// נבדוק אם השפה היא פייתון
	for i, input := range question.Inputs {
		var result string
		var errRun error

		if strings.Contains(answer.Language, "python") {
			result, errRun = runPythonCode(answer.Code, question.FunctionSignature, input)
		} else if strings.Contains(answer.Language, "js") {
			result, errRun = runJavaScriptCode(answer.Code, question.FunctionSignature, input)
		}

		if errRun != nil {
			return false, errRun
		}

		// בדיקת תוצאה מול הפלט הצפוי עבור כל אינפוט
		isPassed, message := checkSingleTest(result, input, question.ExpectedOutputs[i])
		if !isPassed {
			log.Printf("Test %d failed: %s", i+1, message)
			return false, nil
		}
	}

	log.Printf("All tests passed successfully!")
	return true, nil
}


// פונקציה להשוואת פלט התשובה עם הציפיות
// פונקציה לבדוק כל אינפוט בנפרד מול הפלט הצפוי
func checkSingleTest(result string, inputSet []interface{}, expectedOutput interface{}) (bool, string) {
	// נוודא שהתוצאה מתאימה לפלט הצפוי
	if !strings.Contains(result, fmt.Sprintf("%v", expectedOutput)) {
		return false, fmt.Sprintf("Input set: %v\nExpected output: %v\nReceived output: %v", inputSet, expectedOutput, result)
	}
	return true, ""
}


// הפונקציה להרצת קוד פייתון
func runPythonCode(code string, signature string, inputs []interface{}) (string, error) {
	// נשלב את הקוד עם הקלטים שהמשתמש שלח
	codeWithInput := fmt.Sprintf(`
def solution(%s):
    %s
result = solution(%v)
print(result)
`, signature, code, formatInputs(inputs))
log.Printf(codeWithInput)
	// נריץ את קוד הפייתון ישירות במיכל של דוקר
	cmd := exec.Command("docker", "exec", "python-container", "python3", "-c", codeWithInput)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running Python code: %s", err)
		return "false", err
	}

	log.Printf("Python code output: %s", string(output))
	return string(output), nil
}

// הפונקציה להרצת קוד C#
func runJavaScriptCode(code string, signature string, inputs []interface{}) (string, error) {
    // נשלב את הקוד עם הקלטים
    codeWithInput := fmt.Sprintf(`
function solution(%s) {
    %s;
}
console.log(solution(%v));
`, signature, code, formatInputs(inputs))

    log.Printf(codeWithInput)

    // נריץ את הקוד ישירות במיכל עם Node.js
    cmd := exec.Command("docker", "exec", "nodejs-container", "node", "-e", codeWithInput)

    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("Error running JavaScript code: %s\nOutput: %s", err, string(output))
        return "false", err
    }

    log.Printf("JavaScript code output: %s", string(output))
    return string(output), nil
}
func formatInputs(inputs []interface{}) string {
	var formattedInputs []string
	for _, input := range inputs {
		// נהפוך את כל קלט לערך מחרוזת (כך ניתן להשתמש בו ב-JS)
		formattedInputs = append(formattedInputs, fmt.Sprintf("%v", input))
	}
	return strings.Join(formattedInputs, ", ")
}