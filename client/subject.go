package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Subject struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Roles []Role `json:"roles,omitempty"`
}

type Role struct {
	Name        string       `json:"name"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	CreateTime  string       `json:"createtime,omitempty"`
	IsDefault   string       `json:"isdefeault,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	CreateTime  string `json:"createtime,omitempty"`
}

func (cli *Client) GetSubject(name string) (*Subject, error) {
	url := fmt.Sprintf("%s/%s", cli.SubjectsBasePath, name)

	err := cli.rateLimiters.RBACGetSubject.Wait(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter getting subject %s", name))
	}

	res, err := sendRequest(cli, cli.backstoryAPIClient, "GET", cli.userAgent, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting subject")
	}

	var subject Subject
	err = json.Unmarshal(res, &subject)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal subject response")
	}

	return &subject, nil
}

func (cli *Client) CreateSubject(subject Subject) error {
	url := cli.SubjectsBasePath

	err := cli.rateLimiters.RBACCreateSubject.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter creating subject %v", subject))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, subject)
	if err != nil {
		return errors.Wrap(err, "failed creating subject")
	}

	return nil
}

func (cli *Client) UpdateSubject(subject Subject) error {
	url := fmt.Sprintf("%s/%s", cli.SubjectsBasePath, subject.Name)

	err := cli.rateLimiters.RBACUpdateSubject.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter updating subject %v", subject))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "PATCH", cli.userAgent, url, subject)
	if err != nil {
		return errors.Wrap(err, "failed updating subject")
	}

	return nil
}

func (cli *Client) DeleteSubject(name string) error {
	url := fmt.Sprintf("%s/%s", cli.SubjectsBasePath, name)

	err := cli.rateLimiters.RBACDeleteSubject.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter deleting subject %s", name))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "DELETE", cli.userAgent, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed deleting subject")
	}

	return nil
}
