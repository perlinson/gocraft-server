package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	blockpb "github.com/perlinson/gocraft-server/internal/proto/block"
	Store "github.com/perlinson/gocraft-server/internal/store"
)

type BlockService struct {
	blockpb.UnimplementedBlockServiceServer
	mu      sync.RWMutex
	chunks  map[string]*ChunkData
	version int64
	store   *Store.Store
}

type ChunkData struct {
	blocks  []int32
	version int64
}

func NewBlockService(store *Store.Store) *BlockService {
	return &BlockService{
		store:  store,
		chunks: make(map[string]*ChunkData),
	}
}

// 区块坐标转字符串键
func chunkKey(p, q int32) string {
	return fmt.Sprintf("%d:%d", p, q)
}

// 实现 FetchChunk RPC
func (s *BlockService) FetchChunk(ctx context.Context, req *blockpb.FetchChunkRequest) (*blockpb.FetchChunkResponse, error) {
	id := Store.Vec3{X: req.P, Y: 0, Z: req.Q}

	s.mu.RLock()
	defer s.mu.RUnlock()

	version := s.store.GetChunkVersion(id)

	response := &blockpb.FetchChunkResponse{
		Version: version,
	}

	if req.Version == version {
		return response, nil
	}
	blocks := make([]*blockpb.Block, 0)
	s.store.RangeBlocks(id, func(bid Store.Vec3, w int) {
		blocks = append(blocks, &blockpb.Block{
			X: bid.X,
			Y: bid.Y,
			Z: bid.Z,
			W: int32(w),
		})
	})

	response.Blocks = blocks
	return response, nil
}

// 实现 UpdateBlock RPC
func (s *BlockService) UpdateBlock(ctx context.Context, req *blockpb.UpdateBlockRequest) (*blockpb.UpdateBlockResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("UpdateBlock: %v", req)
	version := Store.GenerateChunkVersion()

	// 更新方块和区块版本
	s.store.UpdateBlock(Store.Vec3{X: req.X, Y: req.Y, Z: req.Z}, int(req.W))
	s.store.UpdateChunkVersion(Store.Vec3{X: req.P, Y: 0, Z: req.Q}, version)

	// 创建响应
	response := &blockpb.UpdateBlockResponse{
		Version: version,
	}

	// TODO: 广播更新给其他玩家的逻辑需要重新设计
	// 可能需要使用某种发布订阅机制或事件系统

	return response, nil
}

func (s *BlockService) StreamChunk(req *blockpb.ChunkRequest, stream blockpb.BlockService_StreamChunkServer) error {
	// 实现流区块数据的逻辑
	// 创建一个定时器用于定期发送区块更新
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// 记录最后发送的版本号
	lastVersion := int64(0)

	for {
		select {
		case <-stream.Context().Done():
			// 如果客户端断开连接，退出循环
			return stream.Context().Err()
		case <-ticker.C:
			s.mu.RLock()
			// 检查是否有新的区块更新
			if s.version > lastVersion {
				// 遍历所有区块
				for key, chunk := range s.chunks {
					if chunk.version > lastVersion {
						// 解析区块坐标
						var p, q int32
						fmt.Sscanf(key, "%d:%d", &p, &q)

						// 发送更新的区块数据
						response := &blockpb.ChunkUpdate{
							P:       p,
							Q:       q,
							Blocks:  chunk.blocks,
							Version: fmt.Sprintf("%d", chunk.version),
						}
						if err := stream.Send(response); err != nil {
							s.mu.RUnlock()
							return err
						}
					}
				}
				lastVersion = s.version
			}
			s.mu.RUnlock()
		}
	}
}
