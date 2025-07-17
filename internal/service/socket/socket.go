package socketservice

import (
	"fmt"
	"github.com/1612421/cinema-booking/internal/controller/httpapi"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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

func (s *SocketService) BroadcastMessageExceptID(message any, exceptedID uuid.UUID) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, conn := range s.clients {
		if i == exceptedID.String() {
			continue
		}

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

func (s *SocketService) HandleWebsocket(authService httpapi.IAuthService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload, err := authService.AuthenticateSocket(ctx)
		if err != nil {
			log.Bg().Error("Verify socket failed", zap.Error(err))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Bg().Error("Upgrade socket connection failed", zap.Error(err))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
		defer conn.Close()

		s.registerClient(authPayload.UserId, conn)
		log.Bg().Info(fmt.Sprintf("Socket user ID %s connected", authPayload.UserId.String()))
		defer s.unregisterClient(authPayload.UserId, conn)

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Bg().Info(fmt.Sprintf("%s disconnected", authPayload.UserId.String()))
				break // client disconnected
			}
		}
	}
}
