package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	
	_ "automated-message-sender/docs"
	"automated-message-sender/models"
)

var (
	isSending bool = false
	db        *gorm.DB
	rdb       *redis.Client
	ctx       = context.Background()
	c         *cron.Cron
)

// @title Automated Message Sender API
// @version 1.0
// @description This API is for sending and managing messages in automated way.
// @host localhost:8080
// @BasePath /
func main() {
	var err error
    dsn := "host=db user=user password=password dbname=messages port=5432 sslmode=disable"
    db, err = connectToDB(dsn)
    if err != nil {
        log.Fatalf("could not connect to database: %v", err)
    }

	db.AutoMigrate(&models.Message{})

	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	app := fiber.New()

	app.Post("/start", StartSending)
	app.Post("/stop", StopSending)
	app.Get("/sent-messages", GetSentMessages)
	app.Post("/send", SendMessageHandler)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	c = cron.New()
	_, err = c.AddFunc("@every 2m", SendMessages)
	if err != nil {
		log.Fatalf("could not schedule cron job: %v", err)
	}
	c.Start()

	log.Fatal(app.Listen(":8080"))
}

// StartSending godoc
// @Summary Start sending messages
// @Description Starts the automatic message sending process
// @Tags Messaging
// @Produce json
// @Success 200 {object} models.StartSendingResponse
// @Router /start [post]
func StartSending(c *fiber.Ctx) error {
	if isSending {
		return c.JSON(fiber.Map{"status": "Message sending is already running"})
	}
	isSending = true
	return c.JSON(fiber.Map{"status": "Message sending started"})
}

// StopSending godoc
// @Summary Stop sending messages
// @Description Stops the automatic message sending process
// @Tags Messaging
// @Produce json
// @Success 200 {object} models.StopSendingResponse
// @Router /stop [post]
func StopSending(c *fiber.Ctx) error {
	if !isSending {
		return c.JSON(fiber.Map{"status": "You should start message sending first"})
	}
	isSending = false
	return c.JSON(fiber.Map{"status": "Message sending stopped"})
}

// GetSentMessages godoc
// @Summary Get sent messages
// @Description Retrieves all messages with status 'sent'
// @Tags Messaging
// @Produce json
// @Success 200 {object} models.GetSentMessagesResponse
// @Router /sent-messages [get]
func GetSentMessages(c *fiber.Ctx) error {
	var messages []models.Message
	if err := db.Where("status = ?", "sent").Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"sentMessages": messages})
}

// SendMessageHandler godoc
// @Summary Send a message
// @Description Sends a message to a specified recipient
// @Tags Messaging
// @Accept json
// @Produce json
// @Param body body models.SendMessageRequest true "Message Content and Recipient"
// @Success 202 {object} models.SendMessageHandlerResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /send [post]
func SendMessageHandler(c *fiber.Ctx) error {
	var reqBody map[string]string
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	content, contentExists := reqBody["content"]
	to, toExists := reqBody["to"]
	if !contentExists || !toExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	messageID := uuid.New().String()

	message := models.Message{
		MessageID: messageID,
		Content:   content,
		Recipient: to,
		Status:    "pending",
	}

	if err := db.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := rdb.Set(ctx, messageID, time.Now().Format(time.RFC3339), 0).Err(); err != nil {
		log.Printf("Error caching message ID: %v", err)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":   "Accepted",
		"messageId": messageID,
	})
}

func SendMessages() {
	if !isSending {
		return
	}

	var messages []models.Message
    if err := db.Where("status = ?", "pending").Limit(2).Find(&messages).Error; err != nil {
        log.Printf("Error retrieving messages: %v", err)
        return
    }

	for _, msg := range messages {
		log.Printf("Processing message: %s with content: %s", msg.Recipient, msg.Content)

		msg.Status = "sent"
		if err := db.Save(&msg).Error; err != nil {
			log.Printf("Error updating message status: %v", err)
		}
	}
}

func connectToDB(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		log.Printf("Database connection failed: %v. Retrying in 3 seconds...", err)
		time.Sleep(3 * time.Second)
	}
	return nil, fmt.Errorf("could not connect to database after multiple attempts: %v", err)
}