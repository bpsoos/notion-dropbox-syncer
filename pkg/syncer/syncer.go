package syncer

import (
	"github.com/bpsoos/notion-dropbox-syncer/pkg/adapters/models"
)

type Syncer struct {
	syncDelayMs int

	dropboxAdapter DropboxAdapter
	notionAdapter  NotionAdapter
	logger         Logger
}

type DropboxAdapter interface {
	GetLink(filePath string) (string, error)
	ListFolder(folderPath string) ([]models.FileInfo, error)
}

type NotionAdapter interface {
	AddDatabaseRow(databaseId, name, link string) error
	ListDatabaseRows(databaseId string) ([]string, error)
}

type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
}

type SyncerConfig struct {
	SyncDelayMs int
}

type SyncerDeps struct {
	DropboxAdapter DropboxAdapter
	NotionAdapter  NotionAdapter
	Logger         Logger
}

func NewSyncer(deps *SyncerDeps, config *SyncerConfig) *Syncer {
	return &Syncer{
		syncDelayMs:    config.SyncDelayMs,
		dropboxAdapter: deps.DropboxAdapter,
		notionAdapter:  deps.NotionAdapter,
		logger:         deps.Logger,
	}
}
