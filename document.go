package deepl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type DocumentParams struct {
	TargetLang lang
	File       string
	Filename   string
}

type documentResponse struct {
	DocumentId  string `json:"document_id"`
	DocumentKey string `json:"document_key"`
}

func (d *documentResponse) GetDocumentKey() string {
	return d.DocumentKey
}

func (d *documentResponse) GetDocumentID() string {
	return d.DocumentId
}

func (c *Client) TranslateDocument(params DocumentParams) (string, string, error) {
	u, err := url.Parse(c.baseURL.String() + "document")
	if err != nil {
		return "", "", err
	}

	err = validateExt(params.File)
	if err != nil {
		return "", "", err
	}

	file, err := os.Open(params.File)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}

	mw := multipart.NewWriter(body)

	fw, err := mw.CreateFormFile("file", filepath.Base(params.File))
	if err != nil {
		return "", "", err
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
		return "", "", err
	}
	if res.StatusCode != 200 {
		var errMessage ErrMessage
		if err = json.Unmarshal(b, &errMessage); err != nil {
			return "", "", err
		}

		return "", "", errMessage.Error()
	}

	var data documentResponse
	if err = json.Unmarshal(b, &data); err != nil {
		return "", "", err
	}

	return data.GetDocumentID(), data.GetDocumentKey(), nil
}

type statusResponse struct {
	DocumentId       string `json:"document_id"`
	Status           string `json:"status"`
	SecondsRemaining int    `json:"seconds_remaining"`
}

func (c *Client) GetTranslateStatus(documentID, documentKey string) (string, error) {
	u, err := url.Parse(c.baseURL.String() + "document/" + documentID)
	if err != nil {
		return "", err
	}

	authkey, _ := c.GetAuthKey()
	q := u.Query()
	q.Add("auth_key", authkey)
	q.Add("document_key", documentKey)

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		var errMessage ErrMessage
		if err = json.Unmarshal(body, &errMessage); err != nil {
			return "", err
		}

		return "", errMessage.Error()
	}

	var d statusResponse
	if err = json.Unmarshal(body, &d); err != nil {
		return "", err
	}

	return d.Status, nil
}

func (c *Client) GetTranslationRemainingTime(documentID, documentKey string) (int, error) {
	u, err := url.Parse(c.baseURL.String() + "document/" + documentID)
	if err != nil {
		return 0, err
	}

	authkey, _ := c.GetAuthKey()
	q := u.Query()
	q.Add("auth_key", authkey)
	q.Add("document_key", documentKey)

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	if res.StatusCode != 200 {
		var errMessage ErrMessage
		if err = json.Unmarshal(body, &errMessage); err != nil {
			return 0, err
		}

		return 0, errMessage.Error()
	}

	var d statusResponse
	if err = json.Unmarshal(body, &d); err != nil {
		return 0, err
	}

	return d.SecondsRemaining, nil
}

func (c *Client) GetTranslatedDocumentSentence(params DocumentParams) (string, error) {
	documentID, documentKey, err := c.TranslateDocument(params)
	if err != nil {
		return "", err
	}

	ch := make(chan int)
	go func(chan<- int) {
		for true {
			status, err := c.GetTranslateStatus(documentID, documentKey)
			if err != nil {
				return
			}

			if status == "done" {
				ch <- 0
				break
			} else if status == "error" {
				return
			}
			fmt.Println(status)

			time.Sleep(1000 * time.Millisecond)
		}
	}(ch)

	select {
	case <-ch:
	case <-c.ctx.Done():
		return "", c.ctx.Err()
	}

	u, err := url.Parse(c.baseURL.String() + "document/" + documentID + "/result")
	if err != nil {
		return "", err
	}

	authkey, _ := c.GetAuthKey()
	q := u.Query()
	q.Add("auth_key", authkey)
	q.Add("document_key", documentKey)

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		var errMessage ErrMessage
		if err = json.Unmarshal(body, &errMessage); err != nil {
			return "", err
		}

		return "", errMessage.Error()
	}

	return string(body), nil
}

func (c *Client) GetTranslatedDocument(params DocumentParams) error {
	id, key, err := c.TranslateDocument(params)
	if err != nil {
		return err
	}

	ch := make(chan int)
	go func(chan<- int) {
		for true {
			status, err := c.GetTranslateStatus(id, key)
			if err != nil {
				return
			}

			if status == "done" {
				ch <- 0
				break
			} else if status == "error" {
				return
			}

			time.Sleep(1000 * time.Millisecond)
		}
	}(ch)

	select {
	case <-ch:
	case <-c.ctx.Done():
		return c.ctx.Err()
	}

	return c.getTranslatedDocument(id, key)
}

func (c *Client) getTranslatedDocument(documentID, documentKey string) error {
	u, err := url.Parse(c.baseURL.String() + "document/" + documentID + "/result")
	if err != nil {
		return err
	}

	authkey, _ := c.GetAuthKey()
	q := u.Query()
	q.Add("auth_key", authkey)
	q.Add("document_key", documentKey)

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var errMessage ErrMessage
		if err = json.Unmarshal(body, &errMessage); err != nil {
			return err
		}

		return errMessage.Error()
	}

	if _, err = os.Stat("deepl"); os.IsNotExist(err) {
		os.Mkdir("deepl", 0777)
		err = os.Chmod("deepl", 0777)
		if err != nil {
			return err
		}

	}
	max := 0
	err = filepath.Walk("deepl", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, "deepl/translated_")
		path = strings.TrimSuffix(path, ".txt")
		i, _ := strconv.Atoi(path)
		if i >= max {
			max = i
		}

		return nil
	})
	if err != nil {
		return err
	}

	file := "deepl/translated_" + strconv.Itoa(max+1) + ".txt"
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		_, err = f.WriteString(scanner.Text())
		if err != nil {
			return err
		}
	}

	return nil
}

func validateExt(file string) error {
	ext := filepath.Ext(file)
	if ext != ".docx" && ext != ".pptx" && ext != ".pdf" && ext != ".html" && ext != ".txt" {
		return errors.New("invalid extension")
	}

	return nil
}
