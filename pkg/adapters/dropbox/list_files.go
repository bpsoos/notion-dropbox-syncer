package dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/bpsoos/notion-dropbox-syncer/pkg/adapters/models"
)

const (
	listFolderContinueUrl = "https://api.dropboxapi.com/2/files/list_folder/continue"
	listFolderUrl         = "https://api.dropboxapi.com/2/files/list_folder"
)

func (da *DropboxAdapter) ListFolder(folderPath string) ([]models.FileInfo, error) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")

	result, err := da.listFolder(folderPath, headers)
	if err != nil {
		return nil, err
	}
	if !result.HasMore {
		return result.FileInfos, nil
	}
	cursor := result.Cursor
	fileInfos := result.FileInfos
	for {
		result, err = da.listFolderContinue(cursor, headers)
		if !result.HasMore {
			break
		}
		fileInfos = append(fileInfos, result.FileInfos...)
		cursor = result.Cursor
		time.Sleep(time.Duration(da.listFolderPagingDelayMs) * time.Millisecond)
	}

	return fileInfos, nil
}

type listFolderResult struct {
	HasMore   bool
	Cursor    string
	FileInfos []models.FileInfo
}

func (da *DropboxAdapter) listFolder(folderPath string, headers http.Header) (*listFolderResult, error) {
	bodyJson, err := json.Marshal(&ListFolderRequest{
		Recursive: false,
		Path:      folderPath,
		Limit:     da.listFolderPageSize,
	})
	if err != nil {
		return nil, err
	}
	res, err := da.doRequest(
		http.MethodPost,
		listFolderUrl,
		headers,
		bodyJson,
	)
	if err != nil {
		return nil, err
	}

	return parseResponse(res)
}

func (da *DropboxAdapter) listFolderContinue(cursor string, headers http.Header) (*listFolderResult, error) {
	body := &ListFolderContinueRequest{Cursor: cursor}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := da.doRequest(
		http.MethodPost,
		listFolderContinueUrl,
		headers,
		bodyJson,
	)
	if err != nil {
		return nil, err
	}

	return parseResponse(res)
}

func parseResponse(res *http.Response) (*listFolderResult, error) {
	parsedResponse := new(ListFolderResponse)
	err := json.NewDecoder(res.Body).Decode(parsedResponse)
	if err != nil {
		return nil, err
	}
	files := make([]models.FileInfo, 0)
	for _, entry := range parsedResponse.Entries {
		if entry.Tag != "file" {
			continue
		}
		files = append(files, models.FileInfo{
			Path: entry.PathDisplay,
			Name: entry.Name,
			Hash: entry.ContentHash,
		})
	}

	return &listFolderResult{
		HasMore:   parsedResponse.HasMore,
		Cursor:    parsedResponse.Cursor,
		FileInfos: files,
	}, nil
}

func (da *DropboxAdapter) doRequest(method, url string, headers http.Header, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	res, err := da.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, errors.New("bad status " + res.Status + " " + string(body))
	}

	return res, nil
}
