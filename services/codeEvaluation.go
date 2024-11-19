package services

import (
	"Backend/models"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"github.com/google/uuid"
)

// ExecuteAnswer receives answer code and runs it inside a Kubernetes Pod (Job)
func ExecuteAnswer(answer models.Answer) (bool, error) {
	// Fetch the question details based on the provided QuestionID
	question, err := FetchQuestionByID(answer.QuestionID)
	if err != nil {
		log.Printf("Error fetching question: %v", err)
		return false, fmt.Errorf("error fetching question: %v", err)
	}

	// Iterate through the inputs and run the code based on the specified language (Python / JS)
	for i, input := range question.Inputs {
		var result string
		var errRun error

		// Run Python or JavaScript code depending on the language
		if strings.Contains(answer.Language, "python") {
			result, errRun = runCodeInKubernetes("python", answer.Code, question.FunctionSignature, input)
		} else if strings.Contains(answer.Language, "js") {
			result, errRun = runCodeInKubernetes("js", answer.Code, question.FunctionSignature, input)
		} else {
			return false, fmt.Errorf("unsupported language: %s", answer.Language)
		}

		// If an error occurred while running the code, log it and return
		if errRun != nil {
			log.Printf("Error running code for Question ID %s: %v", answer.QuestionID, errRun)
			return false, errRun
		}

		// Check if the result matches the expected output for each input
		isPassed, message := checkSingleTest(result, input, question.ExpectedOutputs[i])
		if !isPassed {
			log.Printf("Test %d failed: %s", i+1, message)
			return false, fmt.Errorf("test %d failed: %s", i+1, message)
		}
	}

	log.Printf("All tests passed successfully!")
	return true, nil
}

// Function to compare the output with the expected result
func checkSingleTest(result string, inputSet []interface{}, expectedOutput interface{}) (bool, string) {
	// Ensure the result matches the expected output
	resultStr := fmt.Sprintf("%v", result)
	expectedStr := fmt.Sprintf("%v", expectedOutput)

	if resultStr != expectedStr {
		return false, fmt.Sprintf("Input set: %v\nExpected output: %v\nReceived output: %v", inputSet, expectedStr, resultStr)
	}
	return true, ""
}

// Unified function to run code inside a Kubernetes Pod (Job) based on language
func runCodeInKubernetes(language, code, signature string, inputs []interface{}) (string, error) {
	var codeWithInput string
	// Build code with input based on language
	if language == "python" {
		codeWithInput = fmt.Sprintf(`
%s
result = solution(%v)
print(result)
`, code, formatInputs(inputs))
	} else if language == "js" {
		codeWithInput = fmt.Sprintf(`
%s;
console.log(solution(%v));
`, code, formatInputs(inputs))
	} else {
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	log.Printf("Running %s code: %s", language, codeWithInput)

	// Create the job YAML file dynamically based on the code
	yamlContent := createK8sJobYAML(language, codeWithInput)

	// Write the YAML to a file
	yamlFilePath := "./tmps/job.yaml"
	if err := os.WriteFile(yamlFilePath, []byte(yamlContent), 0644); err != nil {
		log.Printf("Error writing YAML file: %v", err)
		return "", fmt.Errorf("error writing YAML file: %v", err)
	}
	log.Printf("YAML file written to: %s", yamlFilePath)

	// Apply the YAML file to Kubernetes using kubectl
	if err := runKubectlApply(yamlFilePath); err != nil {
		return "", fmt.Errorf("error applying YAML: %v", err)
	}

	// Wait for the job to finish and retrieve the logs
	logs, err := waitForJobAndGetLogs("function-test-job")
	if err != nil {
		// Log the error and check the pod status to determine failure reason
		log.Printf("Error retrieving logs: %v", err)
		podStatus, podErr := getPodStatus("function-test-job")
		if podErr != nil {
			log.Printf("Error checking pod status: %v", podErr)
		} else {
			log.Printf("Pod status: %s", podStatus)
		}
		return "", fmt.Errorf("error retrieving logs: %v", err)
	}

	// Delete the job after execution
	defer func() {
		if deleteErr := runKubectlDelete("function-test-job"); deleteErr != nil {
			log.Printf("Error deleting job: %v", deleteErr)
		}
	}()

	return logs, nil
}

// getPodStatus checks the status of the Kubernetes Pod to diagnose failure
func getPodStatus(jobName string) (string, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-l", fmt.Sprintf("job-name=%s", jobName), "-o", "jsonpath='{.items[0].status.phase}'")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error getting pod status: %v", err)
		return "", err
	}

	return out.String(), nil
}

// createK8sJobYAML creates a dynamic YAML file for the Kubernetes job based on the language and code
func createK8sJobYAML(language, code string) string {
	// Build the command array according to the language
	var command []string

	// Set command according to the language
	if language == "python" {
		lines := strings.Split(code, "\n")
		command = append([]string{"python", "-c"}, lines...)
	} else if language == "js" {
		lines := strings.Split(code, "\n")
		command = append([]string{"node", "-e"}, lines...)
	} else {
		log.Printf("Unsupported language: %s", language)
		return ""
	}

	// Convert the command array to a string format for YAML
	commandStr := fmt.Sprintf(`"%s"`, command[0])
	for _, part := range command[1:] {
		commandStr += fmt.Sprintf(`, "%s"`, part)
	}
	jobName := generateUniqueJobName()
	// Build the YAML content for the job
	yamlContent := fmt.Sprintf(`
apiVersion: batch/v1
kind: Job
metadata:
  name: function-test-job-%s
spec:
  template:
    spec:
      containers:
      - name: code-executor
        image: %s
        command: [%s]
      restartPolicy: OnFailure
`,jobName, getImageForLanguage(language), commandStr)

	// Print the YAML content to the console
	fmt.Println("Generated YAML content:")
	fmt.Println(yamlContent)

	return yamlContent
}

// runKubectlApply applies the YAML file using kubectl
func runKubectlApply(yamlFilePath string) error {
	cmd := exec.Command("kubectl", "apply", "-f", yamlFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		log.Printf("Error applying YAML: %v", err)
		return fmt.Errorf("error applying YAML: %v", err)
	}
	return nil
}

// waitForJobAndGetLogs waits for the job to complete and retrieves the logs
func waitForJobAndGetLogs(jobName string) (string, error) {
	// Wait for the Job to complete
	cmd := exec.Command("kubectl", "wait", fmt.Sprintf("job/%s", jobName), "--for=condition=complete", "--timeout=20s")
	err := cmd.Run()
	if err != nil {
		log.Printf("Error waiting for job to complete: %v", err)
		return "", err
	}

	// Get the logs from the job's pod
	cmd = exec.Command("kubectl", "logs", "-l", fmt.Sprintf("job-name=%s", jobName))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error getting logs: %v", err)
		return "", err
	}

	return out.String(), nil
}

// runKubectlDelete deletes the job after execution
func runKubectlDelete(jobName string) error {
	cmd := exec.Command("kubectl", "delete", "job", jobName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error deleting job: %v", err)
		return fmt.Errorf("error deleting job: %v", err)
	}
	return nil
}


// Get the appropriate Docker image based on the language
func getImageForLanguage(language string) string {
	switch language {
	case "python":
		return "python:3.13-slim"
	case "js":
		return "node:18-slim"
	default:
		log.Printf("Unsupported language: %s", language)
		return ""
	}
}

// Helper function to format inputs
func formatInputs(inputs []interface{}) string {
	var formattedInputs []string
	for _, input := range inputs {
		// Convert each input to a string value (so it can be used in JS)
		formattedInputs = append(formattedInputs, fmt.Sprintf("%v", input))
	}
	return strings.Join(formattedInputs, ", ")
}
func generateUniqueJobName() string {
	return fmt.Sprintf("%s", uuid.New().String())
}
