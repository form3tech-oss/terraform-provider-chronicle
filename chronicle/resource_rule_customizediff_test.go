package chronicle

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	customizeDiffValidRuleText = `rule singleEventRule2{meta:      author = "securityuser"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"

	customizeDiffOtherValidRuleText = `rule singleEventRule2{meta:      author = "newAuthor"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"

	customizeDiffInvalidContext = `generic::invalid_argument invalid token line: 1-2, column: 3-4`
)

// verifyRuleMock is a test double for Chronicle's "rules:verifyRule" endpoint.
// It always answers with the configured validation result and records every
// rule_text it was sent with, so tests can assert whether (and how often)
// VerifyYARARule was actually called.
type verifyRuleMock struct {
	valid   bool
	context string

	mu       sync.Mutex
	requests []string
}

func (m *verifyRuleMock) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, ":verifyRule") {
			t.Fatalf("unexpected request to mock rule endpoint: %s %s", r.Method, r.URL.String())
		}

		var body struct {
			RuleText string `json:"ruleText"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed decoding verifyRule request body: %s", err)
		}

		m.mu.Lock()
		m.requests = append(m.requests, body.RuleText)
		m.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": m.valid,
			"context": m.context,
		})
	}
}

func (m *verifyRuleMock) callCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.requests)
}

// newTestRuleClient returns a *chronicle.Client whose rule endpoint points at
// a local httptest server driven by mock. This uses the same WithRuleBasePath
// mechanism the provider uses for the "rule_custom_endpoint" setting, so
// VerifyYARARule exercises real HTTP request/response handling.
func newTestRuleClient(t *testing.T, mock *verifyRuleMock) *chronicle.Client {
	t.Helper()

	server := httptest.NewServer(mock.handler(t))
	t.Cleanup(server.Close)

	return newTestRuleClientAt(t, server.URL)
}

// newTestRuleClientAt returns a *chronicle.Client with a fake access token
// (so no real Google auth call is made) whose RuleBasePath is redirected to
// baseURL, mirroring what the provider does for "rule_custom_endpoint".
func newTestRuleClientAt(t *testing.T, baseURL string, opts ...chronicle.Option) *chronicle.Client {
	t.Helper()

	opts = append([]chronicle.Option{chronicle.WithBackstoryAPIAccessToken("test-token")}, opts...)

	cli, err := chronicle.NewClient(
		chronicle.RegionEurope,
		"terraform-provider-chronicle-test",
		context.Background(),
		opts...,
	)
	if err != nil {
		t.Fatalf("failed creating test client: %s", err)
	}

	cli.WithRuleBasePath(baseURL + "/v2/detect/rules")

	return cli
}

// diffRule computes the plan-time diff for resourceRule, which is where
// CustomizeDiff (resourceRuleCustomizeDiff) runs.
func diffRule(t *testing.T, state *terraform.InstanceState, rawConfig map[string]interface{}, meta interface{}) (*terraform.InstanceDiff, error) {
	t.Helper()
	config := terraform.NewResourceConfigRaw(rawConfig)
	return resourceRule().Diff(context.Background(), state, config, meta)
}

// existingRuleState simulates the prior state of an already-created
// chronicle_rule, so tests can plan an update instead of a creation.
func existingRuleState(t *testing.T, ruleText string, liveEnabled, alertingEnabled bool) *terraform.InstanceState {
	t.Helper()

	setupClient := newTestRuleClient(t, &verifyRuleMock{valid: true, context: "identified no known errors"})

	createDiff, err := diffRule(t, nil, map[string]interface{}{
		"rule_text":        ruleText,
		"live_enabled":     liveEnabled,
		"alerting_enabled": alertingEnabled,
	}, setupClient)
	if err != nil {
		t.Fatalf("failed preparing prior state: %s", err)
	}

	state := new(terraform.InstanceState).MergeDiff(createDiff)
	state.ID = "test-rule-id"

	return state
}

func TestResourceRuleCustomizeDiff_VerifiesOnNewResource(t *testing.T) {
	mock := &verifyRuleMock{valid: true, context: "identified no known errors"}
	cli := newTestRuleClient(t, mock)

	_, err := diffRule(t, nil, map[string]interface{}{
		"rule_text":        customizeDiffValidRuleText,
		"live_enabled":     false,
		"alerting_enabled": false,
	}, cli)
	if err != nil {
		t.Fatalf("unexpected error planning new rule: %s", err)
	}

	if got := mock.callCount(); got != 1 {
		t.Fatalf("expected exactly 1 call to verifyRule, got %d", got)
	}
	if mock.requests[0] != customizeDiffValidRuleText {
		t.Fatalf("expected verifyRule to be called with configured rule_text, got %q", mock.requests[0])
	}
}

func TestResourceRuleCustomizeDiff_FailsPlanOnInvalidRule(t *testing.T) {
	mock := &verifyRuleMock{valid: false, context: customizeDiffInvalidContext}
	cli := newTestRuleClient(t, mock)

	_, err := diffRule(t, nil, map[string]interface{}{
		"rule_text":        customizeDiffValidRuleText,
		"live_enabled":     false,
		"alerting_enabled": false,
	}, cli)

	if err == nil {
		t.Fatal("expected terraform plan to fail for a rule Chronicle rejects, got nil error")
	}
	if !strings.Contains(err.Error(), customizeDiffInvalidContext) {
		t.Fatalf("expected error to surface Chronicle's compilation context, got: %s", err)
	}
	if got := mock.callCount(); got != 1 {
		t.Fatalf("expected exactly 1 call to verifyRule, got %d", got)
	}
}

func TestResourceRuleCustomizeDiff_VerifiesWhenRuleTextChanges(t *testing.T) {
	priorState := existingRuleState(t, customizeDiffValidRuleText, false, false)

	mock := &verifyRuleMock{valid: true, context: "identified no known errors"}
	cli := newTestRuleClient(t, mock)

	_, err := diffRule(t, priorState, map[string]interface{}{
		"rule_text":        customizeDiffOtherValidRuleText,
		"live_enabled":     false,
		"alerting_enabled": false,
	}, cli)
	if err != nil {
		t.Fatalf("unexpected error planning rule_text update: %s", err)
	}

	if got := mock.callCount(); got != 1 {
		t.Fatalf("expected exactly 1 call to verifyRule when rule_text changes, got %d", got)
	}
}

func TestResourceRuleCustomizeDiff_SkipsWhenOnlyLiveEnabledChanges(t *testing.T) {
	priorState := existingRuleState(t, customizeDiffValidRuleText, false, false)

	mock := &verifyRuleMock{valid: true, context: "identified no known errors"}
	cli := newTestRuleClient(t, mock)

	_, err := diffRule(t, priorState, map[string]interface{}{
		"rule_text":        customizeDiffValidRuleText,
		"live_enabled":     true,
		"alerting_enabled": false,
	}, cli)
	if err != nil {
		t.Fatalf("unexpected error planning live_enabled toggle: %s", err)
	}

	if got := mock.callCount(); got != 0 {
		t.Fatalf("expected verifyRule NOT to be called when only live_enabled changes, got %d calls", got)
	}
}

func TestResourceRuleCustomizeDiff_SkipsWhenOnlyAlertingEnabledChanges(t *testing.T) {
	priorState := existingRuleState(t, customizeDiffValidRuleText, false, false)

	mock := &verifyRuleMock{valid: true, context: "identified no known errors"}
	cli := newTestRuleClient(t, mock)

	_, err := diffRule(t, priorState, map[string]interface{}{
		"rule_text":        customizeDiffValidRuleText,
		"live_enabled":     false,
		"alerting_enabled": true,
	}, cli)
	if err != nil {
		t.Fatalf("unexpected error planning alerting_enabled toggle: %s", err)
	}

	if got := mock.callCount(); got != 0 {
		t.Fatalf("expected verifyRule NOT to be called when only alerting_enabled changes, got %d calls", got)
	}
}

func TestResourceRuleCustomizeDiff_SkipsWhenRuleTextUnknown(t *testing.T) {
	mock := &verifyRuleMock{valid: true, context: "identified no known errors"}
	cli := newTestRuleClient(t, mock)

	// "74D93920-ED26-11E3-AC10-0800200C9A66" is terraform-plugin-sdk's
	// hcl2shim.UnknownVariableValue sentinel (it lives in an internal
	// package, so it can't be imported directly). NewResourceConfigRaw
	// documents this literal as the supported way to mark a raw config
	// value as not known until apply, e.g. rule_text interpolated from
	// another resource's computed attribute.
	const unknownValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

	_, err := diffRule(t, nil, map[string]interface{}{
		"rule_text":        unknownValue,
		"live_enabled":     false,
		"alerting_enabled": false,
	}, cli)
	if err != nil {
		t.Fatalf("unexpected error planning with unknown rule_text: %s", err)
	}

	if got := mock.callCount(); got != 0 {
		t.Fatalf("expected verifyRule NOT to be called when rule_text is unknown, got %d calls", got)
	}
}

func TestResourceRuleCustomizeDiff_SkipsWhenRuleTextBecomesUnknownOnUpdate(t *testing.T) {
	// Prior state has a concrete rule_text, so planning an update to an
	// unknown value produces HasChange=true while the new value is not known.
	// This is the case where diff.NewValueKnown is load-bearing: HasChange
	// alone would let it through.
	priorState := existingRuleState(t, customizeDiffValidRuleText, false, false)

	mock := &verifyRuleMock{valid: true, context: "identified no known errors"}
	cli := newTestRuleClient(t, mock)

	const unknownValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

	_, err := diffRule(t, priorState, map[string]interface{}{
		"rule_text":        unknownValue,
		"live_enabled":     false,
		"alerting_enabled": false,
	}, cli)
	if err != nil {
		t.Fatalf("unexpected error planning update to unknown rule_text: %s", err)
	}

	if got := mock.callCount(); got != 0 {
		t.Fatalf("expected verifyRule NOT to be called when new rule_text is unknown, got %d calls", got)
	}
}

func TestResourceRuleCustomizeDiff_FailsPlanWhenVerifyResponseOmitsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Chronicle's verifyRule response omits "success" entirely for some
		// invalid rules. The YARALValidation.Valid zero value (false) must
		// still be treated as invalid rather than panicking or being ignored.
		_, _ = w.Write([]byte(`{"context":"` + customizeDiffInvalidContext + `"}`))
	}))
	t.Cleanup(server.Close)

	cli := newTestRuleClientAt(t, server.URL)

	_, err := diffRule(t, nil, map[string]interface{}{
		"rule_text":        customizeDiffValidRuleText,
		"live_enabled":     false,
		"alerting_enabled": false,
	}, cli)

	if err == nil {
		t.Fatal(`expected plan to fail when the verifyRule response omits "success"`)
	}
	if !strings.Contains(err.Error(), customizeDiffInvalidContext) {
		t.Fatalf("expected error to surface Chronicle's compilation context, got: %s", err)
	}
}

func TestResourceRuleCustomizeDiff_FailsPlanOnAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	cli := newTestRuleClientAt(t, server.URL, chronicle.WithRequestAttempts(1))

	_, err := diffRule(t, nil, map[string]interface{}{
		"rule_text":        customizeDiffValidRuleText,
		"live_enabled":     false,
		"alerting_enabled": false,
	}, cli)

	if err == nil {
		t.Fatal("expected plan to fail when the verifyRule API call itself errors")
	}
	if !strings.Contains(err.Error(), "error verifying YARA-L 2.0 rule during plan") {
		t.Fatalf("expected a clearly-labelled plan-time verification error, got: %s", err)
	}
}

func TestResourceRuleCreate_StillVerifiesBeforeCreating(t *testing.T) {
	mock := &verifyRuleMock{valid: false, context: customizeDiffInvalidContext}
	cli := newTestRuleClient(t, mock)

	d := schema.TestResourceDataRaw(t, resourceRule().Schema, map[string]interface{}{
		"rule_text": customizeDiffValidRuleText,
	})

	err := resourceRuleCreate(d, cli)
	if err == nil {
		t.Fatal("expected resourceRuleCreate to fail when VerifyYARARule fails")
	}
	if !strings.Contains(err.Error(), "error verifying YARA-L 2.0 rule") {
		t.Fatalf("expected verification error, got: %s", err)
	}
	// Only the verifyRule call should have happened; CreateRule must not run.
	if got := mock.callCount(); got != 1 {
		t.Fatalf("expected exactly 1 API call (verifyRule only), got %d", got)
	}
}

func TestResourceRuleUpdate_StillVerifiesBeforeCreatingVersion(t *testing.T) {
	mock := &verifyRuleMock{valid: false, context: customizeDiffInvalidContext}
	cli := newTestRuleClient(t, mock)

	d := schema.TestResourceDataRaw(t, resourceRule().Schema, map[string]interface{}{
		"rule_text": customizeDiffValidRuleText,
	})
	d.SetId("existing-rule-id")

	err := resourceRuleUpdate(d, cli)
	if err == nil {
		t.Fatal("expected resourceRuleUpdate to fail when VerifyYARARule fails")
	}
	if !strings.Contains(err.Error(), "error verifying YARA-L 2.0 rule") {
		t.Fatalf("expected verification error, got: %s", err)
	}
	// Only the verifyRule call should have happened; CreateRuleVersion must not run.
	if got := mock.callCount(); got != 1 {
		t.Fatalf("expected exactly 1 API call (verifyRule only), got %d", got)
	}
}
