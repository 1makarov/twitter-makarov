package twt

import "time"

type AddRulesResponse struct {
	Meta struct {
		Sent    time.Time `json:"sent"`
		Summary struct {
			Created    int `json:"created"`
			NotCreated int `json:"not_created"`
			Valid      int `json:"valid"`
			Invalid    int `json:"invalid"`
		} `json:"summary"`
	} `json:"meta"`
	Errors []struct {
		Value string `json:"value"`
		ID    string `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"errors"`
}

type DeleteRulesResponse struct {
	Meta struct {
		Sent    time.Time `json:"sent"`
		Summary struct {
			Deleted    int `json:"deleted"`
			NotDeleted int `json:"not_deleted"`
		} `json:"summary"`
	} `json:"meta"`
	Errors []struct {
		Errors []struct {
			Parameters struct {
			} `json:"parameters"`
			Message string `json:"message"`
		} `json:"errors"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Type   string `json:"type"`
	} `json:"errors"`
}

type ValidateRulesResponse struct {
	Data []struct {
		Value string `json:"value"`
		Tag   string `json:"tag"`
		ID    string `json:"id"`
	} `json:"data"`
	Meta struct {
		Sent    time.Time `json:"sent"`
		Summary struct {
			Created    int `json:"created"`
			NotCreated int `json:"not_created"`
			Valid      int `json:"valid"`
			Invalid    int `json:"invalid"`
		} `json:"summary"`
	} `json:"meta"`
	Errors []struct {
		Value   string   `json:"value"`
		Details []string `json:"details"`
		Title   string   `json:"title"`
		Type    string   `json:"type"`
	} `json:"errors"`
}

type GetListRulesResponse struct {
	Data []Rule `json:"data"`
	Meta struct {
		Sent time.Time `json:"sent"`
	} `json:"meta"`
}
