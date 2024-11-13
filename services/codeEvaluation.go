package services

import (
	"Backend/models"
	"fmt"
	"log"
	"strings"
	"github.com/fsouza/go-dockerclient"
)

// ExecuteAnswer receives answer code and runs it inside a Docker container
func ExecuteAnswer(answer models.Answer) (bool, error) {
	// Check the language type (Python / JS) of the code the user sent
	question, err := FetchQuestionByID(answer.QuestionID)
	if err != nil {
		log.Printf("Error fetching question: %v", err)
		return false, err
	}

	// Check if the language is Python
	for i, input := range question.Inputs {
		var result string
		var errRun error

		if strings.Contains(answer.Language, "python") {
			result, errRun = runPythonCode(answer.Code, question.FunctionSignature, input)
		} else if strings.Contains(answer.Language, "js") {
			result, errRun = runJavaScriptCode(answer.Code, question.FunctionSignature, input)
		}

		if errRun != nil {
			log.Printf("Error running code: %v", errRun)
			return false, errRun
		}

		// Check if the result matches the expected output for each input
		isPassed, message := checkSingleTest(result, input, question.ExpectedOutputs[i])
		if !isPassed {
			log.Printf("Test %d failed: %s", i+1, message)
			return false, nil
		}
	}

	log.Printf("All tests passed successfully!")
	return true, nil
}

// Function to compare the output with the expected result
func checkSingleTest(result string, inputSet []interface{}, expectedOutput interface{}) (bool, string) {
	// Ensure the result matches the expected output
	if !strings.Contains(result, fmt.Sprintf("%v", expectedOutput)) {
		return false, fmt.Sprintf("Input set: %v\nExpected output: %v\nReceived output: %v", inputSet, expectedOutput, result)
	}
	return true, ""
}

// Function to run Python code inside Docker
// Function to run Python code inside Docker
// Function to run Python code inside Docker
func runPythonCode(code string, signature string, inputs []interface{}) (string, error) {
	// Combine the code with the inputs the user sent
	codeWithInput := fmt.Sprintf(`
def solution(%s):
    %s
result = solution(%v)
print(result)
`, signature, code, formatInputs(inputs))

	log.Printf("Running Python code: %s", codeWithInput)

	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Printf("Error creating Docker client: %v", err)
		return "", err
	}

	// Create a Docker container with the Python image
	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: "python-test-container",
		Config: &docker.Config{
			Image: "python:3.13-slim",
			Cmd:   []string{"python3", "-c", codeWithInput},
		},
	})
	if err != nil {
		log.Printf("Error creating container: %v", err)
		return "", err
	}

	log.Println("Starting container...")

	err = client.StartContainer(container.ID, nil)
	if err != nil {
		log.Printf("Error starting container: %v", err)
		return "", err
	}

	// Attach to container output
	err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container: container.ID,
		Stream: true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		log.Printf("Error attaching to container: %v", err)
		return "", err
	}

	// Read logs from the container
	logs, err := client.WaitContainer(container.ID)
if err != nil {
    log.Printf("Error waiting for container to finish: %v", err)
    return "", err
}

log.Printf("Container logs: %s", logs)
 //return logs, nil
	// Wait for container to finish and capture the exit status
	_, err = client.WaitContainer(container.ID)
	if err != nil {
		log.Printf("Error waiting for container to finish: %v", err)
		return "", err
	}

	// log.Printf("Final logs from container: %s", logsOutput)
	defer client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true})

	return "", nil
}


// Function to run JavaScript code inside Docker
func runJavaScriptCode(code string, signature string, inputs []interface{}) (string, error) {
	// Combine the code with the inputs
	codeWithInput := fmt.Sprintf(`
function solution(%s) {
    %s;
}
console.log(solution(%v));
`, signature, code, formatInputs(inputs))

	log.Printf("Running JavaScript code: %s", codeWithInput)

	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Printf("Error creating Docker client: %v", err)
		return "", err
	}

	// Create a Docker container with the Node.js image
	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: "nodejs-test-container",
		Config: &docker.Config{
			Image: "tsivyal/my-nodejs-image:latest", // Docker image for Node.js
			Cmd:   []string{"node", "-e", codeWithInput},
		},
	})
	if err != nil {
		log.Printf("Error creating container: %v", err)
		return "", err
	}
	defer client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true})

	// Start the container
	err = client.StartContainer(container.ID, nil)
	if err != nil {
		log.Printf("Error starting container: %v", err)
		return "", err
	}

	// Attach to the container's output
	logs, err := client.WaitContainer(container.ID)
	if err != nil {
		log.Printf("Error waiting for container to finish: %v", err)
		return "", err
	}

	log.Printf("JavaScript code output: %s", logs)
	return fmt.Sprintf("%d", logs), nil
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
