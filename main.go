package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/bpsoos/notion-dropbox-syncer/internal/logging"
	"github.com/bpsoos/notion-dropbox-syncer/pkg/adapters/dropbox"
	"github.com/bpsoos/notion-dropbox-syncer/pkg/adapters/notion"
	configpkg "github.com/bpsoos/notion-dropbox-syncer/pkg/config"
	"github.com/bpsoos/notion-dropbox-syncer/pkg/syncer"
	"golang.org/x/oauth2"
)

const dropboxAuthEndpoint = "https://api.dropbox.com/oauth2/authorize"
const dropboxTokenEndpoint = "https://api.dropbox.com/oauth2/token"

func main() {
	config := configpkg.MustLoadConfig()
	logger := logging.GetLogger(config.LogLevel)
	oauthConf := oauth2.Config{
		ClientID:     config.Dropbox.AppKey,
		ClientSecret: config.Dropbox.AppSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  dropboxAuthEndpoint,
			TokenURL: dropboxTokenEndpoint,
		},
	}
	ctx := context.Background()
	dropboxHttpClient := oauthConf.Client(ctx, &oauth2.Token{
		Expiry:       time.Now(),
		RefreshToken: config.Dropbox.RefreshToken,
	})

	dropboxAdapter := dropbox.NewDropboxAdapter(
		dropboxHttpClient,
		&dropbox.DropboxAdapterConfig{
			ListFolderPageSize:      config.Dropbox.ListFolderPageSize,
			ListFolderPagingDelayMs: config.Dropbox.ListFolderPagingDelayMs,
		},
	)
	notionAdapter := notion.NewNotionAdapter(
		http.DefaultClient,
		&notion.NotionAdapterConfig{
			ListDatabaseRowsPageSize:      config.Dropbox.ListFolderPageSize,
			ListDatabaseRowsPagingDelayMs: config.Dropbox.ListFolderPagingDelayMs,
			Token:                         config.Notion.ApiKey,
		})

	syncer := syncer.NewSyncer(
		&syncer.SyncerDeps{
			DropboxAdapter: dropboxAdapter,
			NotionAdapter:  notionAdapter,
			Logger:         logger,
		},
		&syncer.SyncerConfig{SyncDelayMs: config.SyncerSyncDelayMs},
	)

	logger.Info(
		"starting sync with path",
		"dropbox_folder_path", config.Dropbox.FolderPath,
	)
	for {
		err := syncer.Sync(config.Notion.DatabaseId, config.Dropbox.FolderPath)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		logger.Info("sleeping")
		time.Sleep(time.Duration(config.WorkSleepSeconds) * time.Second)
	}
}
