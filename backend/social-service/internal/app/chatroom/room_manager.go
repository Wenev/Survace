package chatroom

import (
	"sync"

	"github.com/Wenev/Survace/proto/gen/dto"
)

type Client struct {
	ID     int32
	RoomID int64
	Ch     chan *dto.ChatEvent
}

type Room struct {
	clients    map[int32]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *dto.ChatEvent
	mu         sync.RWMutex
}

func newRoom() *Room {
	room := &Room{
		clients:    make(map[int32]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *dto.ChatEvent),
	}
	go room.run()
	return room
}

func (r *Room) run() {
	for {
		select {
		case c := <-r.Register:
			r.mu.Lock()
			r.clients[c.ID] = c
			r.mu.Unlock()
		case c := <-r.Unregister:
			r.mu.Lock()
			delete(r.clients, c.ID)
			close(c.Ch)
			r.mu.Unlock()
		case event := <-r.Broadcast:
			var receiverId, senderId int32
			switch e := event.Event.(type) {
			case *dto.ChatEvent_Message:
				receiverId = e.Message.ReceiverId
				senderId = e.Message.SenderId
			case *dto.ChatEvent_Typing:
				receiverId = e.Typing.ReceiverId
				senderId = e.Typing.SenderId
			case *dto.ChatEvent_Unsend:
				receiverId = e.Unsend.ReceiverId
				senderId = e.Unsend.SenderId
			}
			r.mu.RLock()
			if c, ok := r.clients[receiverId]; ok {
				c.Ch <- event
			}
			if senderId != receiverId {
				if c, ok := r.clients[senderId]; ok {
					c.Ch <- event
				}
			}
			r.mu.RUnlock()
		}
	}
}

type Manager struct {
	rooms map[int64]*Room
	mu    sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[int64]*Room),
	}
}

func (m *Manager) GetRoom(roomId int64) *Room {
	m.mu.RLock()
	room, ok := m.rooms[roomId]
	m.mu.RUnlock()
	if !ok {
		m.mu.Lock()
		room = newRoom()
		m.rooms[roomId] = room
		m.mu.Unlock()
	}
	return room
}

func CalculateRoomId(userId1, userId2 int32) int64 {
	var smallerId, largerId int32
	if userId1 < userId2 {
		smallerId = userId1
		largerId = userId2
	} else {
		smallerId = userId2
		largerId = userId1
	}
	return int64(smallerId)*100000 + int64(largerId)
}
