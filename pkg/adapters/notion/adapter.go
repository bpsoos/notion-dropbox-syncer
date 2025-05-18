package notion

import "net/http"

const NotionApiVersion = "2022-06-28"

type NotionAdapter struct {
	listDatabaseRowsPageSize      int
	listDatabaseRowsPagingDelayMs int
	token                         string

	client *http.Client
}

type NotionAdapterConfig struct {
	ListDatabaseRowsPageSize      int
	ListDatabaseRowsPagingDelayMs int
	Token                         string
}

func NewNotionAdapter(client *http.Client, config *NotionAdapterConfig) *NotionAdapter {
	return &NotionAdapter{
		listDatabaseRowsPageSize:      config.ListDatabaseRowsPageSize,
		listDatabaseRowsPagingDelayMs: config.ListDatabaseRowsPagingDelayMs,
		token:                         config.Token,
		client:                        client,
	}
}
