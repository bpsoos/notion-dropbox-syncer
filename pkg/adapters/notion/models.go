package notion

type DatabaseQueryRequest struct {
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size"`
}

type DatabaseQueryResponse struct {
	Object  string `json:"object"`
	Results []struct {
		Properties struct {
			Name NameProperty `json:"Name"`
		} `json:"properties"`
	} `json:"results"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

type TitleType struct {
	Text TextType `json:"text"`
}

type TextType struct {
	Content string `json:"content"`
}

type NameProperty struct {
	Title []TitleType `json:"title"`
}

type AddRowParent struct {
	DatabaseID string `json:"database_id"`
}

type LinkProperty struct {
	URL string `json:"url"`
}

type AddRowProperties struct {
	Name NameProperty `json:"Name"`
	Link LinkProperty `json:"Link"`
}

type AddRowRequest struct {
	Parent     AddRowParent     `json:"parent"`
	Properties AddRowProperties `json:"properties"`
}
