package services

import (
	"Backend/models"
	"bytes"
	"context"
	"fmt"
	"log"
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
	return true, nil
}

// validateSyntax checks if the provided code has syntax errors
func validateSyntax(language, code string) error {
	var cmd *exec.Cmd
	log.Printf("Code received: %s", code)

	// הכנת הפקודה בהתאם לשפת הקוד
	if language == "python" {
		// הפעלת python עם דגל -c שמבצע את הקוד שמועבר לו
		cmd = exec.Command("python3", "-c", code) // שים לב, python3 במקום python
	} else if language == "js" {
		// הפעלת node עם דגל -e שמבצע את הקוד שמועבר לו
		cmd = exec.Command("node", "-e", code)
	} else {
		return fmt.Errorf("unsupported language: %s", language)
	}

	// הכנת משתנה לתפיסת שגיאות stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// הפעלת הפקודה ורישום פלט השגיאה אם יש
	if err := cmd.Run(); err != nil {
		log.Printf("Command failed: %v", err)
		log.Printf("stderr: %s", stderr.String())  // הדפסת שגיאות שקשורות לפלט stderr
		return fmt.Errorf("syntax validation failed: %s", stderr.String())
	}

	log.Printf("Syntax validation passed for %s code.", language)
	return nil
}


// checkSingleTest compares the output with the expected result
func checkSingleTest(result string, inputSet []interface{}, expectedOutput interface{}) (bool, string) {
	// Ensure the result matches the expected output
	if !strings.Contains(result, fmt.Sprintf("%v", expectedOutput)) {
		return false, fmt.Sprintf("Input set: %v\nExpected output: %v\nReceived output: %v", inputSet, expectedOutput, result)
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



// Create the Job specification
func createKubernetesJobSpec(language, code string) *batchv1.Job {
	var command []string
	var args []string

	// Adjust the command based on the language
	if language == "python" {
		// Split the code into lines, if necessary
		lines := strings.Split(code, "\n")
		// python -c expects code as a string, so pass the whole code as args
		args = append(args, strings.Join(lines, " "))  // Join lines into one string for the argument
		command = []string{"python", "-c"}
	} else if language == "js" {
		// For JavaScript, we need to send the code as a single argument
		args = append(args, code)
		command = []string{"node", "-e"}
	} else {
		log.Printf("Unsupported language: %s", language)
		return nil
	}

	// Define the Kubernetes Job spec
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("code-execution-%s", time.Now().Format("20060102150405")),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "code-executor",
							Image:   getImageForLanguage(language),
							Command: command,  // The main command to run (python or node)
							Args:    args,      // The code as arguments for the command
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever, // Prevent container from restarting on failure
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
	// Wait for the Job to complete
	log.Printf("Waiting for job to complete...")
	time.Sleep(10 * time.Second) // wait for pod to start and finish execution

	// Use the job's label selector to find the corresponding pod
	labelSelector := fmt.Sprintf("job-name=%s", job.Name)
	podsClient := clientset.CoreV1().Pods("default")
	podList, err := podsClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("Error getting pods: %v", err)
		return "", fmt.Errorf("error getting pods: %v", err)
	}

	if len(podList.Items) == 0 {
		return "", fmt.Errorf("no pods found for job %s", job.Name)
	}

	// Assuming the first pod is the one we're interested in
	podName := podList.Items[0].Name

	// Get logs from the pod using the correct PodLogOptions
	podLogs, err := podsClient.GetLogs(podName, &corev1.PodLogOptions{}).Stream(context.TODO())
	if err != nil {
		log.Printf("Error getting pod logs: %v", err)
		return "", fmt.Errorf("error getting pod logs: %v", err)
	}
	defer podLogs.Close()

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(podLogs)
	if err != nil {
		log.Printf("Error reading logs from pod: %v", err)
		return "", fmt.Errorf("error reading logs from pod: %v", err)
	}

	return buffer.String(), nil
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
