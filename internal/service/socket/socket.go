package socketservice

import (
	"fmt"
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type SocketService struct {
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
}

func NewSocketService() *SocketService {
	return &SocketService{
		clients: make(map[string]*websocket.Conn),
	}
}

func (s *SocketService) registerClient(userID uuid.UUID, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[userID.String()] = conn
}

func (s *SocketService) unregisterClient(userID uuid.UUID, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, userID.String())
}

func (s *SocketService) BroadcastMessage(message any) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, conn := range s.clients {
		if err := conn.WriteJSON(message); err != nil {
			_ = conn.WriteJSON(message)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *SocketService) HandleWebsocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	userSessionValue := authservice.GetUserSessionValue(ctx).(*authservice.AuthPayload)
	if userSessionValue == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorx.New(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	s.registerClient(userSessionValue.UserId, conn)
	fmt.Printf("%s connected\n", userSessionValue.UserId.String())
	defer s.unregisterClient(userSessionValue.UserId, conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			fmt.Printf("%s disconnected\n", userSessionValue.UserId.String())
			break // client disconnected
		}
	}

}
