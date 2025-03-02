package main

import (
	"backend-demo/api/healthcheck"
	"backend-demo/models"
	"backend-demo/pkg/config"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	db      *sql.DB
	rdb     *redis.Client
	mutex   sync.Mutex
	running bool
)

var ctx = context.Background()

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	appConfig := config.Read()

	defer zap.L().Sync()
	zap.L().Info("app starting...")

	initDB(appConfig)
	initRedis()

	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	healthCheckHandler := healthcheck.NewHealthCheckHandler()
	app.Get("/health", func(c *fiber.Ctx) error {
		req := &healthcheck.HealthCheckRequest{}
		resp, err := healthCheckHandler.Handle(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	app.Get("/start", func(c *fiber.Ctx) error {
		zap.L().Info("Starting message processing...")
		mutex.Lock()
		if !running {
			running = true
			go processMessages()
		}
		mutex.Unlock()
		return c.SendString("started")
	})
	app.Get("/stop", func(c *fiber.Ctx) error {
		zap.L().Info("Stopping message processing...")
		mutex.Lock()
		running = false
		mutex.Unlock()
		return c.SendString("stopped")
	})

	app.Get("/messages", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, recipient, content, sent, created_at FROM messages WHERE sent = TRUE")
		if err != nil {
			zap.L().Error("Error querying messages", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		var messages []models.Message
		for rows.Next() {
			var msg models.Message
			rows.Scan(&msg.ID, &msg.Recipient, &msg.Content, &msg.Sent, &msg.CreatedAt)
			messages = append(messages, msg)
		}
		return c.JSON(messages)
	})

	go func() {
		if err := app.Listen(fmt.Sprintf("0.0.0.0:%s", appConfig.Port)); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	zap.L().Info("Server started on port", zap.String("port", appConfig.Port))

	gracefulShutdown(app)
}

func gracefulShutdown(app *fiber.App) {
	// Create channel for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	zap.L().Info("Shutting down server...")

	// Shutdown with 5 second timeout
	if err := app.ShutdownWithTimeout(3 * time.Second); err != nil {
		zap.L().Error("Error during server shutdown", zap.Error(err))
	}

	zap.L().Info("Server gracefully stopped")
}

func initDB(config *config.AppConfig) {

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.DatabaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, config.DatabaseName, config.DatabaseSSLMode))
	if err != nil {
		zap.L().Fatal(err.Error())
	}
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func processMessages() {
	for {
		zap.L().Info("Processing messages...")
		mutex.Lock()
		if !running {
			mutex.Unlock()
			return
		}
		mutex.Unlock()
		rows, err := db.Query("SELECT id, recipient, content FROM messages WHERE sent = FALSE LIMIT 2")
		if err != nil {
			zap.L().Error("Error querying messages", zap.Error(err))
			return
		}
		var messages []models.Message
		for rows.Next() {
			var msg models.Message
			rows.Scan(&msg.ID, &msg.Recipient, &msg.Content)
			messages = append(messages, msg)
		}
		for _, msg := range messages {
			zap.L().Info("Processing message", zap.Int("id", msg.ID))
			err := sendMessage(msg)
			if err != nil {
				zap.L().Error("Error sending message", zap.Error(err))
				return
			}
		}
		zap.L().Info("Messages processed")
		time.Sleep(2 * time.Minute)
	}
}

func sendMessage(msg models.Message) error {
	payload, _ := json.Marshal(map[string]string{
		"to":      msg.Recipient,
		"content": msg.Content,
	})
	resp, err := http.Post("https://webhook.site/f0c1ed93-8a60-41d8-a661-7b9ce22774d7", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var res map[string]string
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		zap.L().Error("Error decoding response", zap.Error(err))
		return err
	}
	zap.L().Info("Message sent", zap.Int("id", msg.ID), zap.Any("response", res))
	if msgID, exists := res["message_id"]; exists {
		rdb.Set(ctx, fmt.Sprintf("msg:%d", msg.ID), msgID, 0)
	}
	_, err = db.Exec("UPDATE messages SET sent = TRUE WHERE id = $1", msg.ID)
	return err
}
