package webpage

import (
	"cli/internal/http"

	"net/url"
)

type Webpage struct {
	URL          string `json:"url"`
	ResponseBody string `json:"response_body"`
	ResponseSize int    `json:"response_size"`
	Error        error  `json:"error"`
}

func NewWebpage(inputUrl string) *Webpage {
	_, err := url.ParseRequestURI(inputUrl)
	if err != nil {
		return &Webpage{
			URL:   inputUrl,
			Error: err,
		}
	}

	return &Webpage{
		URL: inputUrl,
	}
}

func (wp *Webpage) Fetch() {
	if wp.Error != nil {
		return
	}

	wp.ResponseBody, wp.ResponseSize, wp.Error = http.Get(wp.URL)
}
