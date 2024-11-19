package metrics

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Metric представляет метрику в памяти и в файле
type Metric struct {
	Name   string `json:"name"`
	Value  int64  `json:"value"`
	Status string `json:"status"`
}

var (
	metricsFile = filepath.Join(".", "metrics.json")
	metricsData = make(map[string]*Metric)

	// Определение счетчиков для Prometheus
	taskCreatedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tasks_created_total",
			Help: "Total number of tasks created",
		},
		[]string{"status"},
	)
	taskUpdatedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tasks_updated_total",
			Help: "Total number of tasks updated",
		},
		[]string{"status"},
	)
	taskDeletedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tasks_deleted_total",
			Help: "Total number of tasks deleted",
		},
		[]string{"status"},
	)
)

func init() {

	prometheus.MustRegister(taskCreatedCounter)
	prometheus.MustRegister(taskUpdatedCounter)
	prometheus.MustRegister(taskDeletedCounter)

	if err := LoadMetrics(); err != nil {
		log.Printf("Ошибка при загрузке метрик: %v", err)
	}
}

// LoadMetrics загружает метрики из JSON-файла или создает файл с начальными значениями, если его нет
func LoadMetrics() error {
	if _, err := os.Stat(metricsFile); os.IsNotExist(err) {
		log.Println("Файл метрик не найден. Создание нового файла с начальными значениями.")

		metricsData["tasks_created_total_created"] = &Metric{Name: "tasks_created_total", Value: 0, Status: "created"}
		metricsData["tasks_updated_total_updated"] = &Metric{Name: "tasks_updated_total", Value: 0, Status: "updated"}
		metricsData["tasks_deleted_total_deleted"] = &Metric{Name: "tasks_deleted_total", Value: 0, Status: "deleted"}

		if err := SaveMetrics(); err != nil {
			log.Printf("Ошибка при сохранении метрик: %v", err)
			return err
		}
		log.Println("Файл метрик успешно создан и заполнен начальными значениями.")
		return nil
	}

	file, err := os.Open(metricsFile)
	if err != nil {
		log.Printf("Ошибка при открытии файла метрик: %v", err)
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&metricsData)
	if err != nil {
		log.Printf("Ошибка при декодировании метрик из файла: %v", err)
		return err
	}

	log.Println("Файл метрик успешно загружен.")

	for _, metric := range metricsData {
		switch metric.Name {
		case "tasks_created_total":
			taskCreatedCounter.WithLabelValues(metric.Status).Add(float64(metric.Value))
		case "tasks_updated_total":
			taskUpdatedCounter.WithLabelValues(metric.Status).Add(float64(metric.Value))
		case "tasks_deleted_total":
			taskDeletedCounter.WithLabelValues(metric.Status).Add(float64(metric.Value))
		}
	}

	log.Println("Метрики успешно восстановлены в Prometheus из файла.")
	return nil
}

// SaveMetrics сохраняет метрики в JSON-файл
func SaveMetrics() error {
	log.Println("Попытка создания и записи в файл метрик...")

	file, err := os.Create(metricsFile)
	if err != nil {
		log.Printf("Ошибка при создании файла метрик: %v", err)
		return err
	}
	defer file.Close()

	log.Println("Файл метрик успешно создан, начинаем запись данных...")

	err = json.NewEncoder(file).Encode(metricsData)
	if err != nil {
		log.Printf("Ошибка при кодировании метрик в JSON: %v", err)
		return err
	}
	log.Println("Метрики успешно сохранены в файл.")
	return nil
}

// IncrementTaskCreated увеличивает счетчик созданных задач для Prometheus и сохраняет в JSON
func IncrementTaskCreated(status string) {
	metricKey := "tasks_created_total_" + status
	if metric, exists := metricsData[metricKey]; exists {
		metric.Value++
	} else {
		metricsData[metricKey] = &Metric{Name: "tasks_created_total", Value: 1, Status: status}
	}
	taskCreatedCounter.WithLabelValues(status).Inc()
	if err := SaveMetrics(); err != nil {
		log.Printf("Ошибка при сохранении метрик после IncrementTaskCreated: %v", err)
	}
}

// StartMetricsServer запускает HTTP-сервер для экспорта метрик Prometheus
func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("Ошибка при запуске сервера метрик: %v", err)
		}
	}()
}
