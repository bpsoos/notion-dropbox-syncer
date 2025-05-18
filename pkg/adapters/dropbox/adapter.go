package dropbox

import "net/http"

type DropboxAdapterConfig struct {
	ListFolderPageSize      int
	ListFolderPagingDelayMs int
}

type DropboxAdapter struct {
	listFolderPageSize      int
	listFolderPagingDelayMs int

	client *http.Client
}

func NewDropboxAdapter(client *http.Client, config *DropboxAdapterConfig) *DropboxAdapter {
	return &DropboxAdapter{
		listFolderPageSize:      config.ListFolderPageSize,
		listFolderPagingDelayMs: config.ListFolderPagingDelayMs,

		client: client,
	}
}
