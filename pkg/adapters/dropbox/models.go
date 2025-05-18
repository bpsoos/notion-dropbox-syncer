package dropbox

type ListFolderRequest struct {
	Recursive bool   `json:"recursive"`
	Path      string `json:"path"`
	Limit     int    `json:"limit"`
}

type ListFolderContinueRequest struct {
	Cursor string `json:"cursor"`
}

type ListFolderResponse struct {
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
	Entries []struct {
		Tag         string `json:".tag"`
		Name        string `json:"name"`
		PathDisplay string `json:"path_display"`
		ContentHash string `json:"content_hash"`
	} `json:"entries"`
}

type CreateSharedLinkWithSettingsRequest struct {
	Path     string                   `json:"path"`
	Settings SharedLinkCreateSettings `json:"settings"`
}

type SharedLinkResponse struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type CreateSharedLinkWithSettingsConflictResponse struct {
	ErrorSummary string `json:"error_summary"`
	Error        struct {
		Tag                     string `json:".tag"`
		SharedLinkAlreadyExists struct {
			Metadata SharedLinkResponse `json:"metadata"`
		} `json:"shared_link_already_exists"`
	} `json:"error"`
}

type SharedLinkCreateSettings struct {
	Access string `json:"access"`
}
