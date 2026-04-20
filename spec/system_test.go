package spec

import (
	"encoding/json"
	"testing"
)

func TestLoadFromJSON(t *testing.T) {
	input := `{
		"name": "test-system",
		"services": {
			"api": {
				"image": {"name": "nginx", "tag": "1.25"}
			}
		}
	}`

	sys, err := LoadFromJSON([]byte(input))
	if err != nil {
		t.Fatalf("LoadFromJSON failed: %v", err)
	}

	if sys.Name != "test-system" {
		t.Errorf("expected name 'test-system', got %q", sys.Name)
	}

	if len(sys.Services) != 1 {
		t.Errorf("expected 1 service, got %d", len(sys.Services))
	}

	api, ok := sys.Services["api"]
	if !ok {
		t.Fatal("expected service 'api'")
	}

	if api.Image.Name != "nginx" {
		t.Errorf("expected image name 'nginx', got %q", api.Image.Name)
	}

	if api.Image.Tag != "1.25" {
		t.Errorf("expected image tag '1.25', got %q", api.Image.Tag)
	}
}

func TestContainerImageFullName(t *testing.T) {
	tests := []struct {
		name     string
		image    ContainerImage
		expected string
	}{
		{
			name:     "name only",
			image:    ContainerImage{Name: "nginx"},
			expected: "nginx",
		},
		{
			name:     "with tag",
			image:    ContainerImage{Name: "nginx", Tag: "1.25"},
			expected: "nginx:1.25",
		},
		{
			name:     "with digest",
			image:    ContainerImage{Name: "nginx", Digest: "sha256:abc123"},
			expected: "nginx@sha256:abc123",
		},
		{
			name:     "digest takes precedence over tag",
			image:    ContainerImage{Name: "nginx", Tag: "1.25", Digest: "sha256:abc123"},
			expected: "nginx@sha256:abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.image.FullName()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestValidateRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "missing name",
			input:   `{"services": {"api": {"image": {"name": "nginx"}}}}`,
			wantErr: true,
			errMsg:  "system name is required",
		},
		{
			name:    "empty services",
			input:   `{"name": "test", "services": {}}`,
			wantErr: true,
			errMsg:  "must have at least one service",
		},
		{
			name:    "missing image name",
			input:   `{"name": "test", "services": {"api": {"image": {}}}}`,
			wantErr: true,
			errMsg:  "image name is required",
		},
		{
			name:    "valid minimal spec",
			input:   `{"name": "test", "services": {"api": {"image": {"name": "nginx"}}}}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadFromJSON([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateConnections(t *testing.T) {
	// Connection to unknown service
	input := `{
		"name": "test",
		"services": {
			"api": {
				"image": {"name": "nginx"},
				"connections": {
					"unknown-service": {"port": 8080, "protocol": "http"}
				}
			}
		}
	}`

	_, err := LoadFromJSON([]byte(input))
	if err == nil {
		t.Error("expected error for connection to unknown service")
	}
}

func TestToJSON(t *testing.T) {
	sys := &System{
		Name: "test",
		Services: map[string]Service{
			"api": {
				Image: ContainerImage{Name: "nginx", Tag: "1.25"},
				Repo: &GitRepo{
					URL: "https://github.com/example/api",
				},
			},
		},
	}

	data, err := sys.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Parse back and verify
	var parsed System
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}

	if parsed.Name != "test" {
		t.Errorf("expected name 'test', got %q", parsed.Name)
	}
}
