package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Rule struct {
	Text              string            `json:"ruleText"`
	ID                string            `json:"ruleId,omitempty"`
	VersionID         string            `json:"versionId,omitempty"`
	Name              string            `json:"ruleName,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	Type              string            `json:"ruleType,omitempty"`
	VersionCreateTime string            `json:"versionCreateTime,omitempty"`
	CompilationState  string            `json:"compilationState,omitempty"`
	CompilationError  string            `json:"compilationError,omitempty"`
	LiveEnabled       bool              `json:"liveRuleEnabled,omitempty"`
	AlertingEnabled   bool              `json:"alertingEnabled,omitempty"`
}

type YARALValidation struct {
	Valid   bool   `json:"success"`
	Context string `json:"context"`
}

func (cli *Client) GetRule(id string) (*Rule, error) {
	url := fmt.Sprintf("%s/%s", cli.RuleBasePath, id)

	err := cli.rateLimiters.DetectionGetRule.Wait(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while getting rule %s", id))
	}

	res, err := sendRequest(cli, cli.backstoryAPIClient, "GET", cli.userAgent, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting rule")
	}

	var rule Rule
	err = json.Unmarshal(res, &rule)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal subject response")
	}

	return &rule, nil
}

func (cli *Client) CreateRule(rule Rule) (string, error) {
	url := cli.RuleBasePath

	err := cli.rateLimiters.DetectionCreateRule.Wait(context.Background())
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while creating rule %v", rule))
	}

	res, err := sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, rule)
	if err != nil {
		return "", errors.Wrap(err, "failed creating rule")
	}

	var ruleRes Rule
	err = json.Unmarshal(res, &ruleRes)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshal subject response")
	}

	return ruleRes.ID, nil
}

func (cli *Client) CreateRuleVersion(rule Rule) error {
	url := fmt.Sprintf("%s/%s:createVersion", cli.RuleBasePath, rule.ID)

	err := cli.rateLimiters.DetectionCreateRuleVersion.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while creating rule version %v", rule))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, rule)
	if err != nil {
		return errors.Wrap(err, "failed creating rule")
	}

	return nil
}

func (cli *Client) ChangeAlertingRule(id string, alertingEnabled bool) error {
	var operation string
	if alertingEnabled {
		operation = "enableAlerting"
	} else {
		operation = "disableAlerting"
	}
	url := fmt.Sprintf("%s/%s:%s", cli.RuleBasePath, id, operation)

	err := cli.rateLimiters.DetectionEnableAlertingRule.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while change alerting enabled on rule %s", id))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed changing alerting in rule")
	}

	return nil
}

func (cli *Client) ChangeLiveRule(id string, liveEnabled bool) error {
	var operation string
	if liveEnabled {
		operation = "enableLiveRule"
	} else {
		operation = "disableLiveRule"
	}
	url := fmt.Sprintf("%s/%s:%s", cli.RuleBasePath, id, operation)

	err := cli.rateLimiters.DetectionEnableLiveRule.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while change live enabled on rule %s", id))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed changing live in rule")
	}

	return nil
}

func (cli *Client) DeleteRule(id string) error {
	url := fmt.Sprintf("%s/%s", cli.RuleBasePath, id)

	err := cli.rateLimiters.DetectionDeleteRule.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while deleting rule %s", id))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "DELETE", cli.userAgent, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed deleting rule")
	}

	return nil
}

func (cli *Client) VerifyYARARule(yaraRule string) (bool, error) {
	url := fmt.Sprintf("%s:verifyRule", cli.RuleBasePath)
	body := map[string]string{
		"ruleText": yaraRule,
	}

	err := cli.rateLimiters.DetectionVerifyYARARule.Wait(context.Background())
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while verifying rule %s", yaraRule))
	}

	res, err := sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, body)
	if err != nil {
		return false, errors.Wrap(err, "failed verifying rule")
	}

	var validation YARALValidation
	err = json.Unmarshal(res, &validation)
	if err != nil {
		return false, errors.Wrap(err, "could not unmarshal validation response")
	}

	if validation.Valid {
		return true, nil
	}

	return false, fmt.Errorf("compilation error: %s", validation.Context)
}
