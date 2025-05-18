package notion

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func (na *NotionAdapter) ListDatabaseRows(databaseId string) ([]string, error) {
	url := "https://api.notion.com/v1/databases/" + databaseId + "/query"
	body := DatabaseQueryRequest{
		StartCursor: "",
		PageSize:    na.listDatabaseRowsPageSize,
	}
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+na.token)
	headers.Add("Content-Type", "application/json")
	headers.Add("Notion-Version", NotionApiVersion)
	rows := make([]string, 0)
	for {
		result, err := na.listDatabaseRows(url, body, headers)
		if err != nil {
			return nil, err
		}
		rows = append(rows, result.rows...)
		if !result.HasMore {
			break
		}
		body = DatabaseQueryRequest{
			StartCursor: result.Cursor,
			PageSize:    3,
		}
		time.Sleep(time.Duration(na.listDatabaseRowsPagingDelayMs) * time.Millisecond)
	}

	return rows, nil
}

type listResult struct {
	HasMore bool
	Cursor  string
	rows    []string
}

func (na *NotionAdapter) listDatabaseRows(
	url string,
	body DatabaseQueryRequest,
	headers http.Header,
) (*listResult, error) {
	jsonBody, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	res, err := na.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, errors.New("bad status " + res.Status + " " + string(body))
	}

	decoder := json.NewDecoder(res.Body)
	parsedResponse := new(DatabaseQueryResponse)
	err = decoder.Decode(parsedResponse)
	if err != nil {
		return nil, err
	}
	rowTitles := make([]string, len(parsedResponse.Results))
	for i := range parsedResponse.Results {
		rowTitles[i] = parsedResponse.Results[i].Properties.Name.Title[0].Text.Content
	}

	return &listResult{
		HasMore: parsedResponse.HasMore,
		Cursor:  parsedResponse.NextCursor,
		rows:    rowTitles,
	}, nil
}
