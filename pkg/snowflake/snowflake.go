package snowflake

import (
	"sync"
	"time"
)

const epoch int64 = 1704067200000 // 2024-01-01 00:00:00 UTC

type Generator struct {
	mu       sync.Mutex
	lastTime int64
	sequence int64
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Next() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & 0xFFF
		if g.sequence == 0 {
			for now <= g.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now
	return ((now - epoch) << 12) | g.sequence
}
