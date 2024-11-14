package services

import (
	"Backend/models"
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

// ExecuteAnswer receives answer code and runs it inside a Docker container
func ExecuteAnswer(answer models.Answer) (bool, error) {
	// Fetch the question details based on the provided QuestionID
	question, err := FetchQuestionByID(answer.QuestionID)
	if err != nil {
		log.Printf("Error fetching question: %v", err)
		return false, err
	}

	// Iterate through the inputs and run the code based on the specified language (Python / JS)
	for i, input := range question.Inputs {
		var result string
		var errRun error

		// Run Python or JavaScript code depending on the language
		if strings.Contains(answer.Language, "python") {
			result, errRun = runCodeInContainer("python", answer.Code, question.FunctionSignature, input)
		} else if strings.Contains(answer.Language, "js") {
			result, errRun = runCodeInContainer("js", answer.Code, question.FunctionSignature, input)
		}

		// If an error occurred while running the code, log it and return
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

// Unified function to run code inside a Docker container based on language
func runCodeInContainer(language, code, signature string, inputs []interface{}) (string, error) {
    var codeWithInput string
    // Build code with input based on language
    if language == "python" {
        codeWithInput = fmt.Sprintf(`
def solution(%s):
    %s
result = solution(%v)
print(result)
`, signature, code, formatInputs(inputs))
    } else if language == "js" {
        codeWithInput = fmt.Sprintf(`
function solution(%s) {
    %s;
}
console.log(solution(%v));
`, signature, code, formatInputs(inputs))
    } else {
        return "", fmt.Errorf("unsupported language: %s", language)
    }

    log.Printf("Running %s code: %s", language, codeWithInput)

    client, err := docker.NewClientFromEnv()
    if err != nil {
        log.Printf("Error creating Docker client: %v", err)
        return "", err
    }

    // Set up container image and command based on the language
    var image string
    var cmd []string
    if language == "python" {
        image = "python:3.13-slim"
        cmd = []string{"python3", "-c", codeWithInput}  // שינוי ל- -c
    } else if language == "js" {
        image = "node:18-slim"
        cmd = []string{"node", "-e", codeWithInput}     // שמירה על -e עבור JavaScript
    }

    // Create the container with the appropriate language image
    container, err := client.CreateContainer(docker.CreateContainerOptions{
        Name: "test-container",
        Config: &docker.Config{
            Image: image,
            Cmd:   cmd,
        },
    })
    if err != nil {
        log.Printf("Error creating container: %v", err)
        return "", err
    }

    // Ensure container is removed, even if there's an error later
    defer client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true})

    log.Println("Starting container...")

    err = client.StartContainer(container.ID, nil)
    if err != nil {
        log.Printf("Error starting container: %v", err)
        return "", err
    }

    // Attach to container output and get logs
    var buffer bytes.Buffer
    err = client.AttachToContainer(docker.AttachToContainerOptions{
        Container:    container.ID,
        OutputStream: &buffer,
        Stream:       true,
        Stdout:       true,
        Stderr:       true,
        Logs:         true,
    })
    if err != nil {
        log.Printf("Error attaching to container: %v", err)
        return "", err
    }

    // Wait for container to finish executing
    exitCode, err := client.WaitContainer(container.ID)
    if err != nil {
        log.Printf("Error waiting for container to finish: %v", err)
        return "", err
    }

    if exitCode != 0 {
        log.Printf("Container exited with non-zero code: %d", exitCode)
        return "", fmt.Errorf("container execution failed with code %d", exitCode)
    }

    // Capture logs from buffer
    logsOutput := buffer.String()
    log.Printf("Container logs: %s", logsOutput)

    // Return the logs output
    return logsOutput, nil
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
