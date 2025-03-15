package services

import (
	"fmt"
	"sync"
	"github.com/perlinson/gocraft-server/proto"
)

type PlayerService struct {
	mu      sync.RWMutex
	server  *Server
	players map[string]*proto.PlayerState
}

func NewPlayerService(server *Server) *PlayerService {
	return &PlayerService{
		server:  server,
		players: make(map[string]*proto.PlayerState),
	}
}

// 实现 gRPC 服务接口
func (s *PlayerService) UpdateState(ctx context.Context, req *proto.UpdateStateRequest) (*proto.UpdateStateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 更新玩家状态
	s.players[req.Id] = req.State

	// 准备响应数据
	resp := &proto.UpdateStateResponse{
		Players: make(map[string]*proto.PlayerState),
	}
	
	// 收集其他玩家状态（排除自己）
	for id, state := range s.players {
		if id != req.Id {
			resp.Players[id] = state
		}
	}
	return resp, nil
}

func (s *PlayerService) RemovePlayer(ctx context.Context, req *proto.RemovePlayerRequest) (*proto.RemovePlayerResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.players, req.Id)
	return &proto.RemovePlayerResponse{}, nil
}

// 玩家上下线处理
func (s *PlayerService) handlePlayerConnection(action string, id int32) {
	playerID := fmt.Sprintf("%d", id)
	
	s.mu.Lock()
	defer s.mu.Unlock()

	switch action {
	case "online":
		s.players[playerID] = &proto.PlayerState{}
	case "offline":
		delete(s.players, playerID)
	}
}