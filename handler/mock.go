package handler

import (
	"log"
	"time"
)

type MockHandler struct {
}

func (m *MockHandler) SendNotification(userID uint, message string) error {
	time.Sleep(150 * time.Millisecond)
	return nil
}

func (m *MockHandler) SendEmail(userID uint, message string) error {
	time.Sleep(150 * time.Millisecond)
	return nil
}

func (m *MockHandler) Handle(msg []byte) error {
	log.Printf("message: %v", string(msg))
	return nil
}
