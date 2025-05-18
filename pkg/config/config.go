package config

import (
	"log/slog"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

const (
	defaultWorkSleepSeconds                    = 30
	defaultNotionListDatabaseRowsPageSize      = 100
	defaultNotionListDatabaseRowsPagingDelayMs = 200
	defaultDropboxListFolderPageSize           = 100
	defaultDropboxListFolderPagingDelayMs      = 200
	defaultSyncerSyncDelayMs                   = 200
)

type NotionConfig struct {
	DatabaseId                    string
	ApiKey                        string
	ListDatabaseRowsPageSize      int
	ListDatabaseRowsPagingDelayMs int
}

type DropboxConfig struct {
	FolderPath              string
	AppKey                  string
	AppSecret               string
	RefreshToken            string
	ListFolderPageSize      int
	ListFolderPagingDelayMs int
}

type Config struct {
	Notion            *NotionConfig
	Dropbox           *DropboxConfig
	WorkSleepSeconds  int
	LogLevel          slog.Level
	SyncerSyncDelayMs int
}

func defaultInt(k *koanf.Koanf, key string, defVal int) int {
	if !k.Exists(key) {
		return defVal
	}

	return k.Int(key)
}

func MustLoadConfig() *Config {
	var k = koanf.New(".")
	k.Load(env.Provider(
		"NDSYNCER_",
		"__",
		func(s string) string {
			return strings.TrimPrefix(s, "NDSYNCER_")
		},
	), nil)

	return &Config{
		Notion: &NotionConfig{
			DatabaseId: k.MustString("NOTION.DATABASE_ID"),
			ApiKey:     k.MustString("NOTION.API_KEY"),
			ListDatabaseRowsPageSize: defaultInt(
				k,
				"NOTION.LIST_DATABASE_ROWS_PAGE_SIZE",
				defaultNotionListDatabaseRowsPageSize,
			),
			ListDatabaseRowsPagingDelayMs: defaultInt(
				k,
				"NOTION.LIST_DATABASE_ROWS_PAGING_DELAY_MS",
				defaultNotionListDatabaseRowsPagingDelayMs,
			),
		},
		Dropbox: &DropboxConfig{
			FolderPath:   k.MustString("DROPBOX.FOLDER_PATH"),
			AppKey:       k.MustString("DROPBOX.APP_KEY"),
			AppSecret:    k.MustString("DROPBOX.APP_SECRET"),
			RefreshToken: k.MustString("DROPBOX.REFRESH_TOKEN"),
			ListFolderPageSize: defaultInt(
				k,
				"DROPBOX.LIST_FOLDER_PAGE_SIZE",
				defaultDropboxListFolderPageSize,
			),
			ListFolderPagingDelayMs: defaultInt(k,
				"DROPBOX.LIST_FOLDER_PAGING_DELAY_MS",
				defaultDropboxListFolderPagingDelayMs,
			),
		},
		WorkSleepSeconds: defaultInt(
			k,
			"WORK_SLEEP_SECONDS",
			defaultWorkSleepSeconds,
		),
		LogLevel: mustLoadLogLevel(k),
		SyncerSyncDelayMs: defaultInt(
			k,
			"SYNCER_SYNC_DELAY_MS",
			defaultSyncerSyncDelayMs,
		),
	}
}

func mustLoadLogLevel(k *koanf.Koanf) slog.Level {
	logLevelString := k.String("LOGLEVEL")
	var logLevel slog.Level
	switch logLevelString {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO", "":
		logLevel = slog.LevelInfo
	default:
		panic("invalid loglevel: " + logLevelString)
	}
	return logLevel
}
