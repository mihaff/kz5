package python

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type PyModel struct{}

func (p *PyModel) RunModel(modelType, algorithm, targetColumn, inputFilePath, outputFilePath string) (map[string]float64, error) {
	var pythonScript string
	switch modelType {
	case "reg":
		pythonScript = "model_reg.py"
	case "class":
		pythonScript = "model_class.py"
	default:
		return nil, fmt.Errorf("unsupported model type: %s", modelType)
	}

	algorithmMapping := map[string]string{
		"Логистическая регрессия": "logistic_regression",
		"Случайный лес":           "random_forest",
		"Линейная регрессия":      "linear_regression",
		"Метод опорных векторов":  "support_vector_machine",
	}

	algorithm, ok := algorithmMapping[algorithm]
	if !ok {
		return nil, fmt.Errorf("unsupported algorithm type: %s", algorithm)
	}

	log.Printf("Executing: python %s %s %s %s %s", pythonScript, algorithm, targetColumn, inputFilePath, outputFilePath)
	cmd := exec.Command("python", pythonScript, algorithm, targetColumn, inputFilePath, outputFilePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println("Stderr:", stderr.String())
		return nil, fmt.Errorf("failed to run Python model: %v", err)
	}
	lines := strings.Split(stdout.String(), "\n")

	var nonEmptyLines []string
	for _, line := range lines {
		if line != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	result := strings.Join(nonEmptyLines, "")

	log.Println("Stdout:", result)
	metrics, err := ParseMetrics(lines)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metrics: %v", err)
	}
	return metrics, nil
}

func ParseMetrics(outputStrings []string) (map[string]float64, error) {
	metrics := make(map[string]float64)

	for _, line := range outputStrings[:len(outputStrings)-1] {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid output format: %s", line)
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		var floatValue float64
		_, err := fmt.Sscanf(value, "%f", &floatValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value as float: %s", value)
		}

		metrics[name] = floatValue
	}

	return metrics, nil
}
