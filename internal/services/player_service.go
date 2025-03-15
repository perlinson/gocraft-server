package services

import (
	"context"
	"sync"

	playerpb "github.com/perlinson/gocraft-server/internal/proto/player"
)

type PlayerService struct {
	playerpb.UnimplementedPlayerServiceServer
	mu      sync.RWMutex
	players map[string]*playerpb.PlayerState
}

func NewPlayerService(server interface{}) *PlayerService {
	return &PlayerService{
		players: make(map[string]*playerpb.PlayerState),
	}
}

// 实现 gRPC 服务接口
func (s *PlayerService) UpdateState(ctx context.Context, req *playerpb.UpdateStateRequest) (*playerpb.UpdateStateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 更新玩家状态
	s.players[req.Id] = req.State

	// 准备响应数据
	resp := &playerpb.UpdateStateResponse{
		Players: make(map[string]*playerpb.PlayerState),
	}

	// 收集其他玩家状态（排除自己）
	for id, state := range s.players {
		if id != req.Id {
			resp.Players[id] = state
		}
	}
	return resp, nil
}

func (s *PlayerService) RemovePlayer(ctx context.Context, req *playerpb.RemovePlayerRequest) (*playerpb.RemovePlayerResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.players, req.Id)
	return &playerpb.RemovePlayerResponse{}, nil
}
