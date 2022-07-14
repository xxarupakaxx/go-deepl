package deepl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type usageResponse struct {
	CharacterCount int `json:"character_count"`
	CharacterLimit int `json:"character_limit"`
}

func (c *Client) CheckCharacterCount() (int, error) {
	u, err := url.Parse(c.baseURL.String() + "usage")
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return 0, err
	}

	req.Header = c.header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	if res.StatusCode >= 400 {
		var errMessage ErrMessage
		if err = json.Unmarshal(body, &errMessage); err != nil {
			return 0, err
		}

		return 0, fmt.Errorf("%s", errMessage.DisplayMessage())
	}

	var data usageResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return data.CharacterCount, nil
}

func (c *Client) CheckCharacterLimit() (int, error) {
	u, err := url.Parse(c.baseURL.String() + "usage")
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return 0, err
	}

	req.Header = c.header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	if res.StatusCode >= 400 {
		var errMessage ErrMessage
		if err = json.Unmarshal(body, &errMessage); err != nil {
			return 0, err
		}

		return 0, fmt.Errorf("%s", errMessage.DisplayMessage())
	}

	var data usageResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return data.CharacterLimit, nil
}
