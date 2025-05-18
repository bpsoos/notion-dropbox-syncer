package notion

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const createPageUrl = "https://api.notion.com/v1/pages/"

func (na *NotionAdapter) AddDatabaseRow(databaseId, name, link string) error {
	body := &AddRowRequest{
		Parent: AddRowParent{
			DatabaseID: databaseId,
		},
		Properties: AddRowProperties{
			Name: NameProperty{Title: []TitleType{{
				Text: TextType{
					Content: name,
				},
			}}},
			Link: LinkProperty{
				URL: link,
			},
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		createPageUrl,
		bytes.NewReader(jsonBody),
	)

	if err != nil {
		return err
	}
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+na.token)
	headers.Add("Content-Type", "application/json")
	headers.Add("Notion-Version", NotionApiVersion)
	req.Header = headers

	res, err := na.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return errors.New("bad status " + res.Status + " " + string(body))
	}

	return nil
}
