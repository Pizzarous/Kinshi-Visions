package canvas_queue

import (
	"kinshi_vision_bot/entities"

	"github.com/bwmarrin/discordgo"
)

type Queue interface {
	AddCanvas(item *QueueItem) (int, error)
	StartPolling(botSession *discordgo.Session)
	GetBotDefaultSettings() (*entities.DefaultSettings, error)
	UpdateDefaultDimensions(width, height int) (*entities.DefaultSettings, error)
	UpdateDefaultBatch(batchCount, batchSize int) (*entities.DefaultSettings, error)
}
