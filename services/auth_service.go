package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
	
	"github.com/perlinson/gocraft-server/proto"
)

type AuthService struct {
	mu         sync.RWMutex
	sessions   map[string]*UserSession
	server     *Server
}

type UserSession struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}

func NewAuthService(server *Server) *AuthService {
	return &AuthService{
		sessions: make(map[string]*UserSession),
		server:   server,
	}
}

// 登录实现
func (s *AuthService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// TODO: 实际用户验证逻辑
	userID := "user123" // 示例用户ID
	
	// 生成会话令牌
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	// 创建会话
	session := &UserSession{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.mu.Lock()
	s.sessions[token] = session
	s.mu.Unlock()

	return &proto.LoginResponse{
		Token:   token,
		Expires: session.ExpiresAt.Unix(),
	}, nil
}

// 登出实现
func (s *AuthService) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, req.Token)
	return &proto.LogoutResponse{}, nil
}

// 生成随机令牌
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}