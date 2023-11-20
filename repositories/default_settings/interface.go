package default_settings

import (
	"context"
	"kinshi_vision_bot/entities"
)

type Repository interface {
	Upsert(ctx context.Context, setting *entities.DefaultSettings) (*entities.DefaultSettings, error)
	GetByMemberID(ctx context.Context, memberID string) (*entities.DefaultSettings, error)
}
