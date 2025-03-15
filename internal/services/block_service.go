package services

import (
	"context"
	"fmt"
	"sync"

	blockpb "github.com/perlinson/gocraft-server/internal/proto/block"
)

type BlockService struct {
	blockpb.UnimplementedBlockServiceServer
	mu      sync.RWMutex
	chunks  map[string]*ChunkData
	version int64
}

type ChunkData struct {
	blocks  []int32
	version int64
}

func NewBlockService() *BlockService {
	return &BlockService{
		chunks: make(map[string]*ChunkData),
	}
}

// 区块坐标转字符串键
func chunkKey(p, q int32) string {
	return fmt.Sprintf("%d:%d", p, q)
}

// 实现 FetchChunk RPC
func (s *BlockService) FetchChunk(ctx context.Context, req *blockpb.FetchChunkRequest) (*blockpb.FetchChunkResponse, error) {
	key := chunkKey(req.P, req.Q)

	s.mu.RLock()
	defer s.mu.RUnlock()

	// 返回现有区块或初始化新区块
	chunk, exists := s.chunks[key]
	if !exists {
		return &blockpb.FetchChunkResponse{
			Blocks:  make([]int32, 16*16*16), // 16x16x16区块
			Version: 0,
		}, nil
	}

	return &blockpb.FetchChunkResponse{
		Blocks:  chunk.blocks,
		Version: chunk.version,
	}, nil
}

// 实现 UpdateBlock RPC
func (s *BlockService) UpdateBlock(ctx context.Context, req *blockpb.UpdateBlockRequest) (*blockpb.UpdateBlockResponse, error) {
	key := chunkKey(req.P, req.Q)

	s.mu.Lock()
	defer s.mu.Unlock()

	// 确保区块存在
	if _, exists := s.chunks[key]; !exists {
		s.chunks[key] = &ChunkData{
			blocks: make([]int32, 16*16*16),
		}
	}

	// 更新具体方块
	index := req.Y*16*16 + req.Z*16 + req.X
	s.chunks[key].blocks[index] = req.W
	s.version++
	s.chunks[key].version = s.version

	return &blockpb.UpdateBlockResponse{
		Version: s.version,
	}, nil
}
