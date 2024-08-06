package models

import "time"

// Message represents a message in the database
// @Description Message object
// @ID Message
// @Property MessageID string true "ID of the message"
// @Property Content string true "Content of the message"
// @Property Recipient string true "Recipient of the message"
// @Property Status string true "Status of the message"
// @Property SentAt string true "Time when the message was sent"
type Message struct {
	MessageID string    `gorm:"primaryKey;size:36"`
	Content   string    `gorm:"size:255"`
	Recipient string    `gorm:"size:255"`
	Status    string    `gorm:"size:50"`
	SentAt    time.Time
}

type StartSendingResponse struct {
	Status string `json:"status"`
}

type StopSendingResponse struct {
	Status string `json:"status"`
}

type GetSentMessagesResponse struct {
	SentMessages []Message `json:"sentMessages"`
}

type SendMessageRequest struct {
	Content string `json:"content" example:"hey there!"`
	To      string `json:"to" example:"+905551111111"`
}

type SendMessageHandlerResponse struct {
	Message   string `json:"message" example:"Accepted"`
	MessageId string `json:"messageId" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}