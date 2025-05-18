package dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const createSharedLinkWithSettingsUrl = "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"

func (da *DropboxAdapter) GetLink(filePath string) (string, error) {
	body := CreateSharedLinkWithSettingsRequest{
		Path:     filePath,
		Settings: SharedLinkCreateSettings{Access: "default"},
	}
	bodyJson, err := json.Marshal(&body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		createSharedLinkWithSettingsUrl,
		bytes.NewReader(bodyJson),
	)
	if err != nil {
		return "", err
	}
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	req.Header = headers
	res, err := da.client.Do(req)
	if err != nil {
		return "", err
	}

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != http.StatusOK {
		if res.StatusCode != http.StatusConflict {
			body, _ := io.ReadAll(res.Body)
			return "", errors.New("bad status " + res.Status + " " + string(body))
		}
		parsedResponse := new(CreateSharedLinkWithSettingsConflictResponse)
		err = decoder.Decode(parsedResponse)
		if err != nil {
			return "", err
		}
		return parsedResponse.Error.SharedLinkAlreadyExists.Metadata.URL, nil
	}
	parsedResponse := new(SharedLinkResponse)
	err = decoder.Decode(parsedResponse)
	if err != nil {
		return "", err
	}
	return parsedResponse.URL, nil
}
