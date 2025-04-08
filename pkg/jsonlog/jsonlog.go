package jsonlog

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LeverOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "Error"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out        io.Writer
	minLevel   Level
	mu         sync.Mutex
	rabbitConn *amqp.Connection
}

func New(out io.Writer, minLevel Level, rabbitMQURL string) (*Logger, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return &Logger{
		out:        out,
		minLevel:   minLevel,
		rabbitConn: conn,
	}, nil
}

func (l *Logger) PrintInfo(message string, properties map[string]any, collectionName string) {

	l.sendToRabbitMQ(LevelInfo, message, properties, collectionName)
}

func (l *Logger) PrintError(err error, properties map[string]any, collectionName string) {
	l.print(LevelError, err.Error(), properties)
	l.sendToRabbitMQ(LevelError, err.Error(), properties, collectionName)
}

func (l *Logger) PrintFatal(err error, properties map[string]any, collectionName string) {
	l.print(LevelFatal, err.Error(), properties)
	l.sendToRabbitMQ(LevelError, err.Error(), properties, collectionName)
	os.Exit(1)
}

func (l *Logger) sendToRabbitMQ(level Level, message string, properties map[string]any, collectionName string) {
	if level < l.minLevel {
		log.Printf("log level below the minimum level: %v", level)
		return
	}

	aux := struct {
		Level      string         `json:"level"`
		Time       string         `json:"time"`
		Message    string         `json:"message"`
		Properties map[string]any `json:"properties,omitempty"`
		Collection string         `json:"collection"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
		Collection: collectionName,
	}

	// Маршалинг данных в JSON
	body, err := json.Marshal(aux)
	if err != nil {
		log.Printf("failed to marshal log message: %v", err)
		return
	}

	// Проверяем соединение с RabbitMQ
	if l.rabbitConn.IsClosed() {
		log.Printf("RabbitMQ connection is closed")
		return
	}

	// Открываем канал для RabbitMQ
	ch, err := l.rabbitConn.Channel()
	if err != nil {
		log.Printf("failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	// Декларируем очередь (если ее нет)
	q, err := ch.QueueDeclare(
		"logs", // Имя очереди
		true,   // Устойчивая очередь
		false,  // Удалить, если не используется
		false,  // Эксклюзивная
		false,  // Блокировка ожидания
		nil,    // Дополнительные аргументы
	)
	if err != nil {
		log.Printf("failed to declare a queue: %v", err)
		return
	}

	// Отправляем сообщение в очередь
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Ключ маршрутизации (queue name)
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("failed to publish a message: %v", err)
		return
	}
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}

func (l *Logger) print(level Level, message string, properties map[string]any) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string         `json:"level"`
		Time       string         `json:"time"`
		Message    string         `json:"message"`
		Properties map[string]any `json:"properties,omitempty"`
		Trace      string         `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}
