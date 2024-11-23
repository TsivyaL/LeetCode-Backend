package services

import (
	"Backend/models"
	"bytes"
	"context"
	"fmt"
	// "io/ioutil"
	"log"
	// "os"
	"os/exec"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ExecuteAnswer receives answer code and runs it inside a Kubernetes Pod (Job)
func ExecuteAnswer(answer models.Answer) (bool, error) {
	// Fetch the question details based on the provided QuestionID
	question, err := FetchQuestionByID(answer.QuestionID)
	if err != nil {
		log.Printf("Error fetching question: %v", err)
		return false, fmt.Errorf("error fetching question: %v", err)
	}
	// Update question status to "In Progress" before running code
	err = updateQuestionStatusToInProgress(answer.QuestionID, answer.Code)
	if err != nil {
		log.Printf("Error updating question status to 'In Progress': %v", err)
		return false, fmt.Errorf("error updating question status to 'In Progress': %v", err)
	}

	// Perform syntax validation before proceeding
	if err := validateSyntax(answer.Language, answer.Code); err != nil {
		log.Printf("Syntax error: %v", err)
		return false, fmt.Errorf("syntax error: %v", err)
	}

	// Create Kubernetes client from kubeconfig
	clientset, err := createKubernetesClient()
	if err != nil {
		log.Printf("Error creating Kubernetes client: %v", err)
		return false, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	// Iterate through the inputs and run the code based on the specified language (Python / JS)
	for i, input := range question.Inputs {
		var result string
		var errRun error

		// Run Python or JavaScript code depending on the language
		if strings.Contains(answer.Language, "python") {
			result, errRun = runCodeInKubernetes(clientset, "python", answer.Code, question.FunctionSignature, input)
		} else if strings.Contains(answer.Language, "js") {
			result, errRun = runCodeInKubernetes(clientset, "js", answer.Code, question.FunctionSignature, input)
		} else {
			return false, fmt.Errorf("unsupported language: %s", answer.Language)
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
			return false, fmt.Errorf("test %d failed: %s", i+1, message)
		}
	}

	log.Printf("All tests passed successfully!")

	// Update the status to "Completed" when all tests pass
	err = updateQuestionStatusToCompleted(answer.QuestionID, answer.Code)
	if err != nil {
		log.Printf("Error updating question status: %v", err)
		return false, fmt.Errorf("error updating question status: %v", err)
	}

	log.Printf("Question status updated to completed!")
	return true, nil
}
func updateQuestionStatusToInProgress(questionID string, code string) error {
	// Assuming you have a function to update the question status in the database
	// You would fetch the question and set its status to "In Progress"
	err := UpdateQuestionStatus(questionID, "In Progress", code)
	if err != nil {
		return fmt.Errorf("failed to update question status to 'In Progress': %v", err)
	}
	return nil
}

// updateQuestionStatusToCompleted updates the status of the question to "Completed"
func updateQuestionStatusToCompleted(questionID string, code string) error {
	// Assuming you have a function to update the question status in the database
	// You would fetch the question and set its status to "Completed" or "Solved"
	err := UpdateQuestionStatus(questionID, "Completed",code)
	if err != nil {
		return fmt.Errorf("failed to update question status: %v", err)
	}
	return nil
}

// validateSyntax checks if the provided code has syntax errors
func validateSyntax(language, code string) error {
	var cmd *exec.Cmd
	log.Printf("Code received: %s", code)

	// Prepare the command according to the programming language
	if language == "python" {
		// Run python with the -c flag to execute the provided code
		cmd = exec.Command("python3", "-c", code) // Note: use python3 instead of python
	} else if language == "js" {
		// Run node with the -e flag to execute the provided code
		cmd = exec.Command("node", "-e", code)
	} else {
		return fmt.Errorf("unsupported language: %s", language)
	}

	// Prepare a variable to capture stderr errors
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command and log any stderr output
	if err := cmd.Run(); err != nil {
		log.Printf("Command failed: %v", err)
		log.Printf("stderr: %s", stderr.String())  // Print errors related to stderr output
		return fmt.Errorf("syntax validation failed: %s", stderr.String())
	}

	log.Printf("Syntax validation passed for %s code.", language)
	return nil
}

// checkSingleTest compares the output with the expected result
func checkSingleTest(result string, inputSet []interface{}, expectedOutput interface{}) (bool, string) {
	var resultValue string
	if strings.Contains(result, "\n") {
		resultValue = strings.TrimSpace(result) // מסיר רווחים בסוף הפלט
	} else {
		resultValue = result
	}

	expectedValue := fmt.Sprintf("%v", expectedOutput)
	if resultValue != expectedValue {
		return false, fmt.Sprintf("Input set: %v\nExpected output: %v\nReceived output: %v", inputSet, expectedValue, resultValue)
	}

	return true, ""
}

// runCodeInKubernetes executes code in a Kubernetes Job
func runCodeInKubernetes(clientset *kubernetes.Clientset, language, code, signature string, inputs []interface{}) (string, error) {
	var codeWithInput string
	if language == "python" {
		codeWithInput = fmt.Sprintf(
			`%s
result = solution(%v)
print(result)
`, code, formatInputs(inputs))
	} else if language == "js" {
		codeWithInput = fmt.Sprintf(
			`%s;
console.log(solution(%v));
`, code, formatInputs(inputs))
	} else {
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	log.Printf("Running %s code: %s", language, codeWithInput)

	// Define the Job specification in Kubernetes
	job := createKubernetesJobSpec(language, codeWithInput)

	// Create the Job in Kubernetes
	jobClient := clientset.BatchV1().Jobs("default")
	job, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating Kubernetes Job: %v", err)
		return "", fmt.Errorf("error creating Kubernetes Job: %v", err)
	}

	// Wait for the job to complete and retrieve the result (logs)
	logs, err := waitForJobAndGetLogs(clientset, job)
	if err != nil {
		log.Printf("Error retrieving logs: %v", err)
		return "", fmt.Errorf("error retrieving logs: %v", err)
	}

	return logs, nil
}




// createKubernetesJobSpec creates the Job spec, adding checks for temporary file creation and errors.
func createKubernetesJobSpec(language, code string) *batchv1.Job {
	
	var command []string
	var args []string

	if language == "python" {
		// Create a temporary Python file
		// tempFilePath := "/tmp/python_code_" + time.Now().Format("20060102150405") + ".py"
		// err := ioutil.WriteFile(tempFilePath, []byte(code), 0644)
		// if err != nil {
		// 	log.Printf("Error creating temporary file for Python code: %v", err)
		// 	return nil
		// }
		// defer os.Remove(tempFilePath) // Ensure cleanup of the temporary file

		command = []string{"python3","-c", code}
	} else if language == "js" {
		args = append(args, code)
		command = []string{"node", "-e"}
	} else {
		log.Printf("Unsupported language: %s", language)
		return nil
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("code-execution-%s", time.Now().Format("20060102150405")),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "code-executor",
							Image:           getImageForLanguage(language),
							Command:         command,
							Args:            args,
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever, // Prevent retries on failure
				},
			},
		},
	}
	return job
}

// Get the appropriate Docker image based on the language
func getImageForLanguage(language string) string {
	if language == "python" {
		return "python:3.13-slim"
	} else if language == "js" {
		return "node:18-slim"
	}
	return ""
}

// Wait for the Kubernetes Job to finish and retrieve the logs
func waitForJobAndGetLogs(clientset *kubernetes.Clientset, job *batchv1.Job) (string, error) {
	log.Printf("Waiting for job %s to complete...", job.Name)

	// Define maximum timeout for waiting
	timeout := time.Minute * 5
	interval := time.Second * 5
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		// Find the pod associated with the job
		labelSelector := fmt.Sprintf("job-name=%s", job.Name)
		podsClient := clientset.CoreV1().Pods("default")
		podList, err := podsClient.List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			log.Printf("Error getting pods: %v", err)
			return "", fmt.Errorf("error getting pods: %v", err)
		}

		// If no pods are found, retry after the interval
		if len(podList.Items) == 0 {
			log.Printf("No pods found for job %s, retrying...", job.Name)
			time.Sleep(interval)
			continue
		}

		// Assume the first pod in the list is the one we are interested in
		pod := podList.Items[0]
		podName := pod.Name
		log.Printf("Found pod for job %s: %s", job.Name, podName)

		// Check if the pod has completed successfully
		if pod.Status.Phase == corev1.PodSucceeded {
			log.Printf("Pod %s completed successfully!", podName)
			// Retrieve and return pod logs
			return getPodLogs(clientset, podName)
		} else if pod.Status.Phase == corev1.PodFailed {
			log.Printf("Pod %s failed with status: %s", podName, pod.Status.Phase)
			// Retrieve logs even if the pod failed to help debug issues
			logs, logErr := getPodLogs(clientset, podName)
			if logErr != nil {
				log.Printf("Error retrieving pod logs: %v", logErr)
			}
			return "", fmt.Errorf("Logs:\n%s", logs)
		}

		// Log the current pod status and retry after the interval
		log.Printf("Pod %s status: %s, retrying in %v...", podName, pod.Status.Phase, interval)
		time.Sleep(interval)
	}

	// If the timeout is reached, return an error
	errMsg := fmt.Sprintf("Timeout waiting for job %s to complete", job.Name)
	log.Printf(errMsg)
	return "", fmt.Errorf(errMsg)
}


// Helper function to fetch logs from a pod
// Helper function to fetch logs from a pod and return only the first few lines
func getPodLogs(clientset *kubernetes.Clientset, podName string) (string, error) {
	podLogOpts := corev1.PodLogOptions{}
	req := clientset.CoreV1().Pods("default").GetLogs(podName, &podLogOpts)
	logStream, err := req.Stream(context.TODO())
	if err != nil {
		return "", fmt.Errorf("error opening log stream: %v", err)
	}
	defer logStream.Close()

	var logsBuffer bytes.Buffer
	_, err = logsBuffer.ReadFrom(logStream)
	if err != nil {
		return "", fmt.Errorf("error reading logs: %v", err)
	}

	logs := logsBuffer.String()

	// Split logs into lines and take the first 4 lines only
	lines := strings.Split(logs, "\n")
	if len(lines) > 4 {
		lines = lines[:4]
	}

	return strings.Join(lines, "\n"), nil
}


// Helper function to format inputs
func formatInputs(inputs []interface{}) string {
	var formattedInputs []string

	// מעבר על כל האיברים במערך
	for _, input := range inputs {
		switch v := input.(type) {
		case string:
			formattedInputs = append(formattedInputs, fmt.Sprintf("\"%s\"", v))
		case []interface{}:
			formattedInputs = append(formattedInputs, formatInputs(v))
		default:
			formattedInputs = append(formattedInputs, fmt.Sprintf("%v", v))
		}
	}

	return strings.Join(formattedInputs, ", ")
}

// Create Kubernetes client using kubeconfig
func createKubernetesClient() (*kubernetes.Clientset, error) {
	kubeconfig := "/root/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("Error building kubeconfig: %v", err)
		log.Println("Kubeconfig path:", kubeconfig)
		return nil, fmt.Errorf("error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error creating Kubernetes client: %v", err)
		return nil, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	return clientset, nil
}