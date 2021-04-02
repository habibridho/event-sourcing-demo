package handler

import "time"

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
