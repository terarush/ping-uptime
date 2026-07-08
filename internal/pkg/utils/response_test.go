package utils

import (
	"encoding/json"
	"testing"
)

// TestSuccessResponseModel_JSONMarshaling verifies that SuccessResponseModel
// (added for swagger documentation) marshals to JSON using the expected
// field names and includes all fields regardless of zero values.
func TestSuccessResponseModel_JSONMarshaling(t *testing.T) {
	model := SuccessResponseModel{
		Data:    map[string]interface{}{"id": float64(1)},
		Message: "ok",
		Error:   "",
	}

	b, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("unexpected error marshaling SuccessResponseModel: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unexpected error unmarshaling result: %v", err)
	}

	if _, ok := decoded["data"]; !ok {
		t.Errorf("expected 'data' key in marshaled JSON, got: %s", string(b))
	}
	if decoded["message"] != "ok" {
		t.Errorf("expected message 'ok', got %v", decoded["message"])
	}
	if decoded["error"] != "" {
		t.Errorf("expected empty error, got %v", decoded["error"])
	}
}

// TestSuccessResponseModel_NilData verifies the model marshals correctly
// even when Data is nil (common in list/empty responses).
func TestSuccessResponseModel_NilData(t *testing.T) {
	model := SuccessResponseModel{
		Data:    nil,
		Message: "no content",
		Error:   "",
	}

	b, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("unexpected error marshaling SuccessResponseModel with nil data: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unexpected error unmarshaling result: %v", err)
	}

	if v, ok := decoded["data"]; !ok || v != nil {
		t.Errorf("expected 'data' key with nil value, got: %v (present: %v)", v, ok)
	}
}

// TestErrorResponseModel_JSONMarshaling verifies that ErrorResponseModel
// (added for swagger documentation) marshals to JSON with the expected shape.
func TestErrorResponseModel_JSONMarshaling(t *testing.T) {
	model := ErrorResponseModel{
		Data:    nil,
		Message: "",
		Error:   "something went wrong",
	}

	b, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("unexpected error marshaling ErrorResponseModel: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unexpected error unmarshaling result: %v", err)
	}

	if decoded["error"] != "something went wrong" {
		t.Errorf("expected error 'something went wrong', got %v", decoded["error"])
	}
	if decoded["message"] != "" {
		t.Errorf("expected empty message, got %v", decoded["message"])
	}
	if _, ok := decoded["data"]; !ok {
		t.Errorf("expected 'data' key to be present in marshaled JSON, got: %s", string(b))
	}
}

// TestErrorResponseModel_Unmarshal verifies round-tripping of the model
// via JSON unmarshal, ensuring field tags map correctly.
func TestErrorResponseModel_Unmarshal(t *testing.T) {
	raw := []byte(`{"data":{"field":"name"},"message":"validation failed","error":"bad request"}`)

	var model ErrorResponseModel
	if err := json.Unmarshal(raw, &model); err != nil {
		t.Fatalf("unexpected error unmarshaling into ErrorResponseModel: %v", err)
	}

	if model.Message != "validation failed" {
		t.Errorf("expected message 'validation failed', got %q", model.Message)
	}
	if model.Error != "bad request" {
		t.Errorf("expected error 'bad request', got %q", model.Error)
	}
	if model.Data == nil {
		t.Errorf("expected non-nil data")
	}
}

// TestSuccessAndErrorResponseModel_SameShape verifies both new swagger
// documentation models share an identical JSON field shape, since they are
// intended to represent the same generic success/error envelope structure.
func TestSuccessAndErrorResponseModel_SameShape(t *testing.T) {
	success := SuccessResponseModel{Data: "x", Message: "m", Error: "e"}
	failure := ErrorResponseModel{Data: "x", Message: "m", Error: "e"}

	sb, err := json.Marshal(success)
	if err != nil {
		t.Fatalf("unexpected error marshaling success model: %v", err)
	}
	fb, err := json.Marshal(failure)
	if err != nil {
		t.Fatalf("unexpected error marshaling failure model: %v", err)
	}

	if string(sb) != string(fb) {
		t.Errorf("expected identical JSON shape, got success=%s failure=%s", string(sb), string(fb))
	}
}