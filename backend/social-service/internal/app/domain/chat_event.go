package domain

type ChatEventKind int

const (
	ChatEventKindMessage ChatEventKind = iota
	ChatEventKindTyping
	ChatEventKindUnsend
)

type ChatEvent struct {
	Kind    ChatEventKind
	Message *ChatEventMessage
	Typing  *ChatEventTyping
	Unsend  *ChatEventUnsend
}

type ChatEventMessage struct {
	Message *Message
}

type ChatEventTyping struct {
	SenderID   int32
	ReceiverID int32
	IsTyping   bool
}

type ChatEventUnsend struct {
	MessageID  int32
	SenderID   int32
	ReceiverID int32
}

func (e *ChatEvent) SenderID() int32 {
	switch e.Kind {
	case ChatEventKindMessage:
		return e.Message.Message.SenderID
	case ChatEventKindTyping:
		return e.Typing.SenderID
	case ChatEventKindUnsend:
		return e.Unsend.SenderID
	}
	return 0
}

func (e *ChatEvent) ReceiverID() int32 {
	switch e.Kind {
	case ChatEventKindMessage:
		return e.Message.Message.ReceiverID
	case ChatEventKindTyping:
		return e.Typing.ReceiverID
	case ChatEventKindUnsend:
		return e.Unsend.ReceiverID
	}
	return 0
}
