package deepl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type DocumentParams struct {
	SourceLang lang
	TargetLang lang
	File       string
	Filename   string
}

type DocumentResponse struct {
	DocumentId  string `json:"document_id"`
	DocumentKey string `json:"document_key"`
}

func (c *Client) TranslateDocument(params DocumentParams) (*DocumentResponse, error) {
	u, err := url.Parse(c.baseURL.String() + "document")
	if err != nil {
		return nil, err
	}

	file, err := os.Open(params.File)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}

	mw := multipart.NewWriter(body)

	fw, err := mw.CreateFormFile("file", filepath.Base(params.File))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, file)
	contentType := mw.FormDataContentType()

	mw.Close()

	authkey, _ := c.GetAuthKey()
	q := u.Query()

	q.Add("target_lang", convertLang(params.TargetLang))
	q.Add("file", params.File)
	q.Add("auth_key", authkey)

	u.RawQuery = q.Encode()

	res, err := http.Post(u.String(), contentType, body)

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {

	}
	if res.StatusCode != 200 {
		var errMessage ErrMessage
		if err = json.Unmarshal(b, &errMessage); err != nil {
			return nil, err
		}

		fmt.Println(errMessage)
		return nil, nil
	}

	var data DocumentResponse
	if err = json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
