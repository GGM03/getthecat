package main

import (
	"fmt"
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
)

type GoogleResponse struct {
	Kind string `json:"kind"`
	URL  struct {
		Type     string `json:"type"`
		Template string `json:"template"`
	} `json:"url"`
	Queries struct {
		Request []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
			SearchType     string `json:"searchType"`
		} `json:"request"`
		NextPage []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
			SearchType     string `json:"searchType"`
		} `json:"nextPage"`
	} `json:"queries"`
	Context struct {
		Title string `json:"title"`
	} `json:"context"`
	SearchInformation struct {
		SearchTime            float64 `json:"searchTime"`
		FormattedSearchTime   string  `json:"formattedSearchTime"`
		TotalResults          string  `json:"totalResults"`
		FormattedTotalResults string  `json:"formattedTotalResults"`
	} `json:"searchInformation"`
	Items []struct {
		Kind        string `json:"kind"`
		Title       string `json:"title"`
		HTMLTitle   string `json:"htmlTitle"`
		Link        string `json:"link"`
		DisplayLink string `json:"displayLink"`
		Snippet     string `json:"snippet"`
		HTMLSnippet string `json:"htmlSnippet"`
		Mime        string `json:"mime"`
		Image       struct {
			ContextLink     string `json:"contextLink"`
			Height          int    `json:"height"`
			Width           int    `json:"width"`
			ByteSize        int    `json:"byteSize"`
			ThumbnailLink   string `json:"thumbnailLink"`
			ThumbnailHeight int    `json:"thumbnailHeight"`
			ThumbnailWidth  int    `json:"thumbnailWidth"`
		} `json:"image"`
	} `json:"items"`

	Error struct {
		Errors []struct {
			Domain       string `json:"domain"`
			Reason       string `json:"reason"`
			Message      string `json:"message"`
			ExtendedHelp string `json:"extendedHelp"`
		} `json:"errors"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type GoogleAPI struct {
	ApiKey string
	apiUrl string
	params req.Param
}

func NewGoogleAPI(key string, cx string) GoogleAPI {
	api := GoogleAPI{ApiKey: key,
		apiUrl: "https://www.googleapis.com/customsearch/v1",
		params: req.Param{"key": key, "cx": cx, "searchType": "image"}}
	return api
}

func (g GoogleAPI) SearchImages(query string) ([]ImgInfo, error) {
	rnd := randrange(5, 20)
	log.Tracef("[GoogleApi] Sending to google: query:\"%s\" page:\"%d\"", query, rnd)
	resp, err := req.Get(g.apiUrl, g.params, req.Param{"q": query,
		"filter":    "1",
		"lowRange":  rnd - 5,
		"highRange": rnd,
	})
	if err != nil {
		log.Infof("[GoogleApi] error requesting google: %v", err)
		return []ImgInfo{}, err
	}
	log.Tracef("[GoogleApi] Google respond: %d", resp.Response().StatusCode)

	var jsonData GoogleResponse
	resp.ToJSON(&jsonData)

	if len(jsonData.Error.Errors) > 0 {
		log.Infof("[GoogleApi] Google respond error: %s", jsonData.Error.Errors[0].Reason)
		return []ImgInfo{}, fmt.Errorf(jsonData.Error.Errors[0].Reason)
	}

	var imgUrls []ImgInfo
	for _, itm := range jsonData.Items {
		imgUrls = append(imgUrls, ImgInfo{Origin: itm.Link, Width: itm.Image.Width, Height: itm.Image.Height})
	}

	return imgUrls, nil
}
