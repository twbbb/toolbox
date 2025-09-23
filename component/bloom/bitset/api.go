package bitset

import (
	"context"
)

// BitSetApi bit存储提供者
type BitSetApi interface {
	SetBit(ctx context.Context, location []uint) error
	// Exists 存在
	Exists(ctx context.Context, location []uint) (bool, error)
}
