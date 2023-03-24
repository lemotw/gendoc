package conflunce

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/lemotw/gendoc/model"
)

func doConflunceAPI(method, uri string, data io.Reader, auth *model.ConflunceAccount) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, data)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(auth.Username, auth.Password)
	req.Header.Set("Content-Type", "application/json")

	// do request
	client := &http.Client{}
	return client.Do(req)
}

// CRUD API

func FetchConfluncePage(url, pageid string, auth *model.ConflunceAccount) (*model.Page, error) {
	domain := url
	if !strings.HasSuffix(url, "/") {
		domain = url + "/"
	}

	res, err := doConflunceAPI(http.MethodGet, domain+"rest/api/content/"+pageid+"?expand=body.storage", nil, auth)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		// parse page
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		page := &model.Page{}
		err = json.Unmarshal(body, page)
		if err != nil {
			return nil, err
		}

		return page, nil
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return nil, errors.New("response body: " + string(body) + ", status code: " + res.Status)
}

func NewConfluncePage(url string, page *model.Page, auth *model.ConflunceAccount) error {
	domain := url
	if !strings.HasSuffix(url, "/") {
		domain = url + "/"
	}

	pagePayload, err := json.Marshal(page)
	if err != nil {
		return err
	}

	res, err := doConflunceAPI(http.MethodPost, domain+"rest/api/content", bytes.NewBuffer(pagePayload), auth)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return errors.New("response body: " + string(body) + ", status code: " + res.Status)
	}

	return nil
}

func PutConfluncePage(url, pageId string, page *model.Page, auth *model.ConflunceAccount) error {
	domain := url
	if !strings.HasSuffix(url, "/") {
		domain = url + "/"
	}

	pagePayload, err := json.Marshal(page)
	if err != nil {
		return err
	}

	res, err := doConflunceAPI(http.MethodPut, domain+"rest/api/content/"+pageId, bytes.NewBuffer(pagePayload), auth)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return errors.New("response body: " + string(body) + ", status code: " + res.Status)
	}

	return nil
}

// fetch child

func FetchConflunceChildPage(url, pageid string, auth *model.ConflunceAccount) (*model.ChildPageResponse, error) {
	domain := url
	if !strings.HasSuffix(url, "/") {
		domain = url + "/"
	}

	res, err := doConflunceAPI(http.MethodGet, domain+"rest/api/content/"+pageid+"/child/page", nil, auth)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		// parse page
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		childPage := &model.ChildPageResponse{}
		err = json.Unmarshal(body, childPage)
		if err != nil {
			return nil, err
		}

		return childPage, nil
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return nil, errors.New("response body: " + string(body) + ", status code: " + res.Status)
}

func WalkAllChildPage(url, pageid string, fn func(c *model.ChildPage) bool, auth *model.ConflunceAccount) error {
	pageStack := []string{pageid}

	for len(pageStack) > 0 {
		pid := pageStack[len(pageStack)-1]
		pageStack = pageStack[:len(pageStack)-1]

		childPage, err := FetchConflunceChildPage(url, pid, auth)
		if err != nil {
			return err
		}

		for i := 0; i < len(childPage.Results); i++ {
			fn(&childPage.Results[i])
		}
	}

	return nil
}
