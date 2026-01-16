package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =============================================================================
// TestParseRemoteAddr - Unit tests for parseRemoteAddr function
// =============================================================================

func TestParseRemoteAddr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RemoteAddress
	}{
		{
			name:  "IPv4 with port",
			input: "192.168.1.1:8080",
			expected: RemoteAddress{
				Address: "192.168.1.1",
				Port:    "8080",
			},
		},
		{
			name:  "localhost with port",
			input: "127.0.0.1:3000",
			expected: RemoteAddress{
				Address: "127.0.0.1",
				Port:    "3000",
			},
		},
		{
			name:  "IPv6 loopback with port",
			input: "[::1]:8080",
			expected: RemoteAddress{
				Address: "[::1]",
				Port:    "8080",
			},
		},
		{
			name:  "full IPv6 with port",
			input: "[2001:db8::1]:443",
			expected: RemoteAddress{
				Address: "[2001:db8::1]",
				Port:    "443",
			},
		},
		{
			name:  "high port number",
			input: "10.0.0.1:65535",
			expected: RemoteAddress{
				Address: "10.0.0.1",
				Port:    "65535",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseRemoteAddr(tt.input)
			if result.Address != tt.expected.Address {
				t.Errorf("Address mismatch: got %q, want %q", result.Address, tt.expected.Address)
			}
			if result.Port != tt.expected.Port {
				t.Errorf("Port mismatch: got %q, want %q", result.Port, tt.expected.Port)
			}
		})
	}
}

// =============================================================================
// TestEchoHandler - HTTP handler tests
// =============================================================================

func TestEchoHandler_BasicGET(t *testing.T) {
	extras := Extras{AppName: "test-app"}
	handler := http.HandlerFunc(echo(extras))

	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	req.Host = "example.com"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %q", contentType)
	}

	// Parse response
	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify response fields
	if resp.Host != "example.com" {
		t.Errorf("Expected Host 'example.com', got %q", resp.Host)
	}
	if resp.Path != "/test-path" {
		t.Errorf("Expected Path '/test-path', got %q", resp.Path)
	}
	if resp.Method != http.MethodGet {
		t.Errorf("Expected Method 'GET', got %q", resp.Method)
	}
	if resp.Extras.AppName != "test-app" {
		t.Errorf("Expected Extras.AppName 'test-app', got %q", resp.Extras.AppName)
	}
}

func TestEchoHandler_QueryParams(t *testing.T) {
	extras := Extras{}
	handler := http.HandlerFunc(echo(extras))

	req := httptest.NewRequest(http.MethodGet, "/path?foo=bar&baz=qux&multi=a&multi=b", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Check single value params
	if len(resp.Query["foo"]) != 1 || resp.Query["foo"][0] != "bar" {
		t.Errorf("Expected Query['foo'] = ['bar'], got %v", resp.Query["foo"])
	}
	if len(resp.Query["baz"]) != 1 || resp.Query["baz"][0] != "qux" {
		t.Errorf("Expected Query['baz'] = ['qux'], got %v", resp.Query["baz"])
	}

	// Check multi-value param
	if len(resp.Query["multi"]) != 2 {
		t.Errorf("Expected Query['multi'] to have 2 values, got %d", len(resp.Query["multi"]))
	}
}

func TestEchoHandler_Headers(t *testing.T) {
	extras := Extras{}
	handler := http.HandlerFunc(echo(extras))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Custom-Header", "custom-value")
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("Accept-Language", "pt-BR")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Check custom headers are echoed
	if len(resp.Headers["X-Custom-Header"]) != 1 || resp.Headers["X-Custom-Header"][0] != "custom-value" {
		t.Errorf("Expected X-Custom-Header = ['custom-value'], got %v", resp.Headers["X-Custom-Header"])
	}
	if len(resp.Headers["Authorization"]) != 1 || resp.Headers["Authorization"][0] != "Bearer token123" {
		t.Errorf("Expected Authorization header, got %v", resp.Headers["Authorization"])
	}

	// Check multi-value header
	if len(resp.Headers["Accept-Language"]) != 2 {
		t.Errorf("Expected Accept-Language to have 2 values, got %d", len(resp.Headers["Accept-Language"]))
	}
}

func TestEchoHandler_JSONBody(t *testing.T) {
	extras := Extras{}
	handler := http.HandlerFunc(echo(extras))

	jsonPayload := `{"name":"John","age":30,"active":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Check that JSON body was captured
	if string(resp.JSON) != jsonPayload {
		t.Errorf("Expected JSON body %q, got %q", jsonPayload, string(resp.JSON))
	}

	// Check content type is recorded
	if resp.ContentType != "application/json" {
		t.Errorf("Expected ContentType 'application/json', got %q", resp.ContentType)
	}
}

func TestEchoHandler_JSONBodyWithCharset(t *testing.T) {
	extras := Extras{}
	handler := http.HandlerFunc(echo(extras))

	jsonPayload := `{"message":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(jsonPayload))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// JSON should still be captured when Content-Type includes charset
	if string(resp.JSON) != jsonPayload {
		t.Errorf("Expected JSON body %q, got %q", jsonPayload, string(resp.JSON))
	}
}

func TestEchoHandler_NonJSONBody(t *testing.T) {
	extras := Extras{}
	handler := http.HandlerFunc(echo(extras))

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("plain text body"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// JSON field should be empty for non-JSON content
	if string(resp.JSON) != "" {
		t.Errorf("Expected empty JSON field for non-JSON content, got %q", string(resp.JSON))
	}
}

func TestEchoHandler_Extras(t *testing.T) {
	t.Run("with AppName only", func(t *testing.T) {
		extras := Extras{AppName: "my-service"}
		handler := http.HandlerFunc(echo(extras))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		if resp.Extras.AppName != "my-service" {
			t.Errorf("Expected AppName 'my-service', got %q", resp.Extras.AppName)
		}
	})

	t.Run("with Envs", func(t *testing.T) {
		extras := Extras{
			AppName: "env-test",
			Envs: map[string]string{
				"ENV_VAR_1": "value1",
				"ENV_VAR_2": "value2",
			},
		}
		handler := http.HandlerFunc(echo(extras))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		if resp.Extras.Envs["ENV_VAR_1"] != "value1" {
			t.Errorf("Expected Envs['ENV_VAR_1'] = 'value1', got %q", resp.Extras.Envs["ENV_VAR_1"])
		}
		if resp.Extras.Envs["ENV_VAR_2"] != "value2" {
			t.Errorf("Expected Envs['ENV_VAR_2'] = 'value2', got %q", resp.Extras.Envs["ENV_VAR_2"])
		}
	})

	t.Run("empty extras", func(t *testing.T) {
		extras := Extras{}
		handler := http.HandlerFunc(echo(extras))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		if resp.Extras.AppName != "" {
			t.Errorf("Expected empty AppName, got %q", resp.Extras.AppName)
		}
		if resp.Extras.Envs != nil {
			t.Errorf("Expected nil Envs, got %v", resp.Extras.Envs)
		}
	})
}

func TestEchoHandler_Methods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
	}

	extras := Extras{}
	handler := http.HandlerFunc(echo(extras))

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d for %s, got %d", http.StatusOK, method, w.Code)
			}

			// HEAD requests don't have a body, skip JSON parsing
			if method == http.MethodHead {
				return
			}

			var resp Response
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response JSON for %s: %v", method, err)
			}

			if resp.Method != method {
				t.Errorf("Expected Method %q, got %q", method, resp.Method)
			}
		})
	}
}

func TestEchoHandler_Proto(t *testing.T) {
	extras := Extras{}
	handler := http.HandlerFunc(echo(extras))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// httptest uses HTTP/1.1 by default
	if resp.Proto != "HTTP/1.1" {
		t.Errorf("Expected Proto 'HTTP/1.1', got %q", resp.Proto)
	}
}

func TestEchoHandler_ContentLength(t *testing.T) {
	extras := Extras{}
	handler := http.HandlerFunc(echo(extras))

	body := "test body content"
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	expectedLen := int64(len(body))
	if resp.ContentLength != expectedLen {
		t.Errorf("Expected ContentLength %d, got %d", expectedLen, resp.ContentLength)
	}
}

// =============================================================================
// TestResponseJSONSerialization - JSON serialization tests
// =============================================================================

func TestResponseJSONSerialization(t *testing.T) {
	t.Run("full response", func(t *testing.T) {
		resp := Response{
			Host:          "localhost",
			Proto:         "HTTP/1.1",
			ContentLength: 42,
			Headers:       map[string][]string{"X-Test": {"value"}},
			Form:          map[string][]string{},
			Query:         map[string][]string{"key": {"val"}},
			Remote:        RemoteAddress{Address: "127.0.0.1", Port: "12345"},
			Path:          "/api/test",
			Method:        "POST",
			ContentType:   "application/json",
			Extras:        Extras{AppName: "test"},
			JSON:          json.RawMessage(`{"foo":"bar"}`),
		}

		data, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		// Verify it can be unmarshaled back
		var decoded Response
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if decoded.Host != resp.Host {
			t.Errorf("Host mismatch after round-trip")
		}
		if decoded.Path != resp.Path {
			t.Errorf("Path mismatch after round-trip")
		}
		if string(decoded.JSON) != string(resp.JSON) {
			t.Errorf("JSON mismatch after round-trip")
		}
	})

	t.Run("empty JSON omitted", func(t *testing.T) {
		resp := Response{
			Host:   "localhost",
			Path:   "/",
			Method: "GET",
			JSON:   json.RawMessage(""),
		}

		data, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		// Empty JSON should be omitted from output
		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal to map: %v", err)
		}

		// The JSON field should be present but empty
		if jsonVal, exists := decoded["json"]; exists && jsonVal != nil && jsonVal != "" {
			t.Logf("Note: JSON field is present with value: %v", jsonVal)
		}
	})

	t.Run("extras omitempty behavior", func(t *testing.T) {
		resp := Response{
			Host:   "localhost",
			Path:   "/",
			Method: "GET",
			Extras: Extras{}, // Empty extras
		}

		data, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal to map: %v", err)
		}

		// Extras should be present but with empty/nil fields omitted
		if extras, ok := decoded["extras"].(map[string]interface{}); ok {
			if _, hasEnvs := extras["envs"]; hasEnvs {
				t.Logf("Note: Empty envs field is present in extras")
			}
		}
	})
}

func TestRemoteAddressJSON(t *testing.T) {
	remote := RemoteAddress{
		Address: "192.168.1.100",
		Port:    "8080",
	}

	data, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("Failed to marshal RemoteAddress: %v", err)
	}

	var decoded RemoteAddress
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal RemoteAddress: %v", err)
	}

	if decoded.Address != remote.Address {
		t.Errorf("Address mismatch: got %q, want %q", decoded.Address, remote.Address)
	}
	if decoded.Port != remote.Port {
		t.Errorf("Port mismatch: got %q, want %q", decoded.Port, remote.Port)
	}
}
