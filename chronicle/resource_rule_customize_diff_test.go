package chronicle

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testValidRuleText = `rule singleEventRule2{meta:      author = "securityuser"      description = "single event rule that should generate detections TEST"
    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"

const testInvalidRuleText = "invalid yara-l\n"

func testUpdatedValidRuleText() string {
	return strings.Replace(testValidRuleText, "securityuser", "newAuthor", 1)
}

type verifyRuleHandler struct {
	verifyCalls *atomic.Int32
	response    map[string]interface{}
	statusCode  int
}

func (h *verifyRuleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && strings.Contains(r.URL.Path, ":verifyRule") {
		h.verifyCalls.Add(1)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var request map[string]string
		if err := json.Unmarshal(body, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if request["ruleText"] == "" {
			http.Error(w, "missing ruleText", http.StatusBadRequest)
			return
		}

		if h.statusCode == 0 {
			h.statusCode = http.StatusOK
		}
		w.WriteHeader(h.statusCode)
		_ = json.NewEncoder(w).Encode(h.response)
		return
	}

	http.NotFound(w, r)
}

func newTestRuleClient(t *testing.T, handler http.Handler) *chronicle.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return chronicle.NewTestClient(server.Client(), server.URL+"/v2/detect/rules")
}

func newVerifyRuleTestServer(t *testing.T, response map[string]interface{}) (*chronicle.Client, *atomic.Int32) {
	t.Helper()

	verifyCalls := &atomic.Int32{}
	handler := &verifyRuleHandler{
		verifyCalls: verifyCalls,
		response:    response,
	}

	return newTestRuleClient(t, handler), verifyCalls
}

func TestResourceRuleCustomizeDiff_NewResource_ValidRuleText(t *testing.T) {
	t.Parallel()

	client, verifyCalls := newVerifyRuleTestServer(t, map[string]interface{}{
		"success": true,
		"context": "identified no known errors",
	})

	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"rule_text": testValidRuleText,
	})

	_, err := resourceRule().Diff(context.Background(), nil, config, client)
	if err != nil {
		t.Fatalf("expected plan to succeed, got: %v", err)
	}

	if verifyCalls.Load() != 1 {
		t.Fatalf("expected VerifyRule to be called once, got %d", verifyCalls.Load())
	}
}

func TestResourceRuleCustomizeDiff_NewResource_InvalidRuleText(t *testing.T) {
	t.Parallel()

	client, verifyCalls := newVerifyRuleTestServer(t, map[string]interface{}{
		"context": "generic::invalid_argument invalid token line: 1-2, column: 3-4",
	})

	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"rule_text": testInvalidRuleText,
	})

	_, err := resourceRule().Diff(context.Background(), nil, config, client)
	if err == nil {
		t.Fatal("expected plan to fail for invalid rule_text")
	}

	if !strings.Contains(err.Error(), "error verifying YARA-L 2.0 rule during plan") {
		t.Fatalf("expected plan-time verification error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "generic::invalid_argument invalid token line: 1-2, column: 3-4") {
		t.Fatalf("expected Chronicle compilation error in plan output, got: %v", err)
	}

	if verifyCalls.Load() != 1 {
		t.Fatalf("expected VerifyRule to be called once, got %d", verifyCalls.Load())
	}
}

func TestResourceRuleCustomizeDiff_UpdateRuleText(t *testing.T) {
	t.Parallel()

	client, verifyCalls := newVerifyRuleTestServer(t, map[string]interface{}{
		"success": true,
		"context": "identified no known errors",
	})

	state := &terraform.InstanceState{
		ID: "rule-123",
		Attributes: map[string]string{
			"rule_text":        testValidRuleText,
			"live_enabled":     "false",
			"alerting_enabled": "false",
		},
	}
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"rule_text":        testUpdatedValidRuleText(),
		"live_enabled":     false,
		"alerting_enabled": false,
	})

	_, err := resourceRule().Diff(context.Background(), state, config, client)
	if err != nil {
		t.Fatalf("expected plan to succeed, got: %v", err)
	}

	if verifyCalls.Load() != 1 {
		t.Fatalf("expected VerifyRule to be called once when rule_text changes, got %d", verifyCalls.Load())
	}
}

func TestResourceRuleCustomizeDiff_UpdateLiveEnabledOnly(t *testing.T) {
	t.Parallel()

	client, verifyCalls := newVerifyRuleTestServer(t, map[string]interface{}{
		"success": true,
		"context": "identified no known errors",
	})

	state := &terraform.InstanceState{
		ID: "rule-123",
		Attributes: map[string]string{
			"rule_text":        testValidRuleText,
			"live_enabled":     "false",
			"alerting_enabled": "false",
		},
	}
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"rule_text":        testValidRuleText,
		"live_enabled":     true,
		"alerting_enabled": false,
	})

	_, err := resourceRule().Diff(context.Background(), state, config, client)
	if err != nil {
		t.Fatalf("expected plan to succeed, got: %v", err)
	}

	if verifyCalls.Load() != 0 {
		t.Fatalf("expected VerifyRule not to be called when only live_enabled changes, got %d calls", verifyCalls.Load())
	}
}

func TestResourceRuleCustomizeDiff_UpdateAlertingEnabledOnly(t *testing.T) {
	t.Parallel()

	client, verifyCalls := newVerifyRuleTestServer(t, map[string]interface{}{
		"success": true,
		"context": "identified no known errors",
	})

	state := &terraform.InstanceState{
		ID: "rule-123",
		Attributes: map[string]string{
			"rule_text":        testValidRuleText,
			"live_enabled":     "false",
			"alerting_enabled": "false",
		},
	}
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"rule_text":        testValidRuleText,
		"live_enabled":     false,
		"alerting_enabled": true,
	})

	_, err := resourceRule().Diff(context.Background(), state, config, client)
	if err != nil {
		t.Fatalf("expected plan to succeed, got: %v", err)
	}

	if verifyCalls.Load() != 0 {
		t.Fatalf("expected VerifyRule not to be called when only alerting_enabled changes, got %d calls", verifyCalls.Load())
	}
}

func TestResourceRuleCreate_CallsVerifyYARARule(t *testing.T) {
	t.Parallel()

	verifyCalls := &atomic.Int32{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, ":verifyRule"):
			verifyCalls.Add(1)
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"context": "identified no known errors",
			})
		case r.Method == http.MethodPost && !strings.Contains(r.URL.Path, ":"):
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"ruleId": "rule-123"})
		case r.Method == http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"ruleText":        testValidRuleText,
				"ruleId":          "rule-123",
				"liveRuleEnabled": false,
				"alertingEnabled": false,
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := chronicle.NewTestClient(server.Client(), server.URL+"/v2/detect/rules")
	d := schema.TestResourceDataRaw(t, resourceRule().Schema, map[string]interface{}{
		"rule_text": testValidRuleText,
	})

	if err := resourceRuleCreate(d, client); err != nil {
		t.Fatalf("expected create to succeed, got: %v", err)
	}

	if verifyCalls.Load() != 1 {
		t.Fatalf("expected VerifyRule to be called once during create, got %d", verifyCalls.Load())
	}
}

func TestResourceRuleUpdate_CallsVerifyYARARuleWhenRuleTextChanges(t *testing.T) {
	t.Parallel()

	verifyCalls := &atomic.Int32{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, ":verifyRule"):
			verifyCalls.Add(1)
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"context": "identified no known errors",
			})
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, ":createVersion"):
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		case r.Method == http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"ruleText":        testUpdatedValidRuleText(),
				"ruleId":          "rule-123",
				"liveRuleEnabled": false,
				"alertingEnabled": false,
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := chronicle.NewTestClient(server.Client(), server.URL+"/v2/detect/rules")

	state := &terraform.InstanceState{
		ID: "rule-123",
		Attributes: map[string]string{
			"rule_text":        testValidRuleText,
			"live_enabled":     "false",
			"alerting_enabled": "false",
		},
	}
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"rule_text":        testUpdatedValidRuleText(),
		"live_enabled":     false,
		"alerting_enabled": false,
	})

	diff, err := resourceRule().Diff(context.Background(), state, config, client)
	if err != nil {
		t.Fatalf("expected plan diff to succeed, got: %v", err)
	}

	_, diags := resourceRule().Apply(context.Background(), state, diff, client)
	if diags.HasError() {
		t.Fatalf("expected update to succeed, got: %v", diags)
	}

	if verifyCalls.Load() != 2 {
		t.Fatalf("expected VerifyRule to be called during plan and update, got %d", verifyCalls.Load())
	}
}

func TestResourceRuleUpdate_SkipsVerifyWhenOnlyLiveEnabledChanges(t *testing.T) {
	t.Parallel()

	verifyCalls := &atomic.Int32{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, ":verifyRule"):
			verifyCalls.Add(1)
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"context": "identified no known errors",
			})
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, ":enableLiveRule"):
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"ruleText":        testValidRuleText,
				"ruleId":          "rule-123",
				"liveRuleEnabled": true,
				"alertingEnabled": false,
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := chronicle.NewTestClient(server.Client(), server.URL+"/v2/detect/rules")

	state := &terraform.InstanceState{
		ID: "rule-123",
		Attributes: map[string]string{
			"rule_text":        testValidRuleText,
			"live_enabled":     "false",
			"alerting_enabled": "false",
		},
	}
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"rule_text":        testValidRuleText,
		"live_enabled":     true,
		"alerting_enabled": false,
	})

	diff, err := resourceRule().Diff(context.Background(), state, config, client)
	if err != nil {
		t.Fatalf("expected plan diff to succeed, got: %v", err)
	}

	_, diags := resourceRule().Apply(context.Background(), state, diff, client)
	if diags.HasError() {
		t.Fatalf("expected update to succeed, got: %v", diags)
	}

	if verifyCalls.Load() != 0 {
		t.Fatalf("expected VerifyRule not to be called when only live_enabled changes, got %d", verifyCalls.Load())
	}
}
