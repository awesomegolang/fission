/*
Copyright 2016 The Fission Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/client-go/1.5/pkg/api"

	"github.com/fission/fission/tpr"
)

func (c *Client) PackageCreate(f *tpr.Package) (*api.ObjectMeta, error) {

	reqbody, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.url("packages"), "application/json", bytes.NewReader(reqbody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := c.handleCreateResponse(resp)
	if err != nil {
		return nil, err
	}

	var m api.ObjectMeta
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (c *Client) PackageGet(m *api.ObjectMeta) (*tpr.Package, error) {
	relativeUrl := fmt.Sprintf("packages/%v", m.Name)
	relativeUrl += fmt.Sprintf("?namespace=%v", m.Namespace)

	resp, err := http.Get(c.url(relativeUrl))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var f tpr.Package
	err = json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (c *Client) PackageUpdate(f *tpr.Package) (*api.ObjectMeta, error) {
	reqbody, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	relativeUrl := fmt.Sprintf("packages/%v", f.Metadata.Name)

	resp, err := c.put(relativeUrl, "application/json", reqbody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var m api.ObjectMeta
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *Client) PackageDelete(m *api.ObjectMeta) error {
	relativeUrl := fmt.Sprintf("packages/%v", m.Name)
	relativeUrl += fmt.Sprintf("?namespace=%v", m.Namespace)
	return c.delete(relativeUrl)
}

func (c *Client) PackageList() ([]tpr.Package, error) {
	resp, err := http.Get(c.url("packages"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	funcs := make([]tpr.Package, 0)
	err = json.Unmarshal(body, &funcs)
	if err != nil {
		return nil, err
	}

	return funcs, nil
}