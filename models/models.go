package models

import "time"

// User представляет собой модель пользователя
type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time `json:"created_at"`
}

// Session представляет модель сессии пользователя
type Session struct {
	SessionID      int       `json:"session_id"`
	UserID         int       `json:"user_id"`
	ExpirationTime time.Time `json:"expiration_time"`
	IPAddress      string    `json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	CreatedAt      time.Time `json:"created_at"`
}

// Shipment представляет модель отправки файла
type Shipment struct {
	ShipmentID   int       `json:"shipment_id"`
	UserID       int       `json:"user_id"`
	ProjectName  string    `json:"project_name"`
	ModelType    string    `json:"model_type"`
	Algorithm    string    `json:"algorithm"`
	TargetColumn string    `json:"target_column"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
}

// File представляет модель файла
type File struct {
	FileID     int       `json:"file_id"`
	ShipmentID int       `json:"shipment_id"`
	FilePath   string    `json:"file_path"`
	Timestamp  time.Time `json:"timestamp"`
}

type ModelMetrics struct {
	MetricID    int
	FileID      int
	MetricName  string
	MetricValue float64
}

func MetricsToDict(metrics []ModelMetrics) map[string]float64 {
	dict := make(map[string]float64)
	for _, metric := range metrics {
		dict[metric.MetricName] = metric.MetricValue
	}
	return dict
}

func ParseMetricsToModelMetrics(fileID int, metrics map[string]float64) []ModelMetrics {
	var modelMetricsList []ModelMetrics

	for name, value := range metrics {
		modelMetric := ModelMetrics{
			FileID:      fileID,
			MetricName:  name,
			MetricValue: value,
		}
		modelMetricsList = append(modelMetricsList, modelMetric)
	}

	return modelMetricsList
}
