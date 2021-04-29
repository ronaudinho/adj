package adj

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		server  *httptest.Server
		want    string
		wantErr string
	}{
		{
			name: "ok",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hello"))
			})),
			want: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name: "empty body",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			})),
			want: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name: "non 200",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})),
			want: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name: "json body",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := json.Marshal(struct {
					msg string `json:"msg"`
				}{
					msg: "hello",
				})
				w.Header().Set("Content-Type", "application/json")
				w.Write(b)
			})),
			want: "99914b932bd37a50b983c5e7c90ae93b",
		},
		{
			name: "json body no header",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := json.Marshal(struct {
					msg string `json:"msg"`
				}{
					msg: "hello",
				})
				w.Write(b)
			})),
			want: "99914b932bd37a50b983c5e7c90ae93b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			got, err := Get(tt.server.URL)
			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("want %s, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr {
					t.Errorf("want %s, got %s", tt.wantErr, err.Error())
					return
				}
			}
			if got.Response != tt.want {
				t.Errorf("want %s, got %s", tt.want, got.Response)
			}
		})
	}
}

func TestGet_refused(t *testing.T) {
	wantErr := "refused"
	_, err := Get("localhost:8080")
	if err == nil {
		t.Errorf("want %s, got nil", wantErr)
	}
	if !strings.Contains(err.Error(), wantErr) {
		t.Errorf("%s does not contain %s", err, wantErr)
	}
}
