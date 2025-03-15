package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	authpb "auth"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	mu       sync.RWMutex
	sessions map[string]*UserSession
}

type UserSession struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}

func NewAuthService(server interface{}) *AuthService {
	return &AuthService{
		sessions: make(map[string]*UserSession),
	}
}

// 生成随机令牌
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// 登录实现
func (s *AuthService) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
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

	return &authpb.LoginResponse{
		Token: token,
		User: &authpb.User{
			Id:   userID,
			Name: "测试用户",
		},
	}, nil
}

// 登出实现
func (s *AuthService) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	s.mu.Lock()
	delete(s.sessions, req.Token)
	s.mu.Unlock()

	return &authpb.LogoutResponse{}, nil
}
