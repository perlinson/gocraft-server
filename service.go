package main

import (
	"sync"
	"fmt"
	"github.com/perlinson/gocraft-server/proto"
)

type PlayerService struct {
	mutex   sync.Mutex
	server  *Server
	players map[string]proto.PlayerState
}

func NewPlayerService(server *Server) *PlayerService {
	s := &PlayerService{
		server:  server,
		players: make(map[string]proto.PlayerState),
	}
	server.SetPlayerCallback(s.onPlayerCallback)
	return s
}

func (s *PlayerService) UpdateState(req *proto.UpdateStateRequest, rep *proto.UpdateStateResponse) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.players[req.Id]; !ok {
		return nil
	}
	s.players[req.Id] = req.State
	rep.Players = make(map[string]proto.PlayerState)
	for id, state := range s.players {
		if id == req.Id {
			continue
		}
		rep.Players[id] = state
	}
	return nil
}

func (s *PlayerService) onPlayerCallback(action string, id int32) {
	switch action {
	case "online":
		s.addPlayer(fmt.Sprintf("%d", id))
	case "offline":
		s.removePlayer(fmt.Sprintf("%d", id))
	}
}

func (s *PlayerService) removePlayer(pid string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.players, pid)
	req := &proto.RemovePlayerRequest{
		Id: pid,
	}
	s.server.RangeSession(func(id int32, sess *Session) {
		if fmt.Sprintf("%d", id) == pid {
			return
		}
		sess.Go("Player.RemovePlayer", req, new(proto.RemovePlayerResponse), nil)
	})
}

func (s *PlayerService) addPlayer(pid string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.players[pid] = proto.PlayerState{}
}
