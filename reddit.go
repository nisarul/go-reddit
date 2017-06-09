// Package reddit implements a basic client for the Reddit API.
package reddit

import (
	"net/http"
	"fmt"
	"encoding/json"
	"errors"
)

// Item describes a Reddit Item.
type Item struct {
	Title string
	URL string
	Comments int `json:"num_comments"`
}

type Response struct {
	Data struct {
		Children []struct {
			Data Item
		}
	}
}

//Get fetches the most recent Items posted to the specific subreddit.
func (i Item) String() string {
	comm := ""
	switch i.Comments {
	case 0:
		//nothing
	case 1:
		comm = " (1 comment)"
	default:
		comm = fmt.Sprintf(" (%d comments)", i.Comments)
	}
	return fmt.Sprintf("==========\n%s%s\n----------\n%s\n",
	i.Title, comm, i.URL)

}

func Get(reddit string) ([]Item, error) {
	url := fmt.Sprintf("http://reddit.com/r/%s.json", reddit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	items := make([]Item, len(r.Data.Children))
	for i, child := range r.Data.Children {
		items[i] = child.Data
	}
	return items, nil
}
