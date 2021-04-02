package worker

type MessageHandler interface {
	Handle(msg []byte) error
}
