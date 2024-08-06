package client

import (
	"github.com/cclose/dnsmasq-api/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestClient_DeleteDNSEntry(t *testing.T) {
	tests := []struct {
		name           string
		hostname       string
		mockStatusCode int
		mockResponse   string
		wantErr        bool
	}{
		{
			name:           "Successful Deletion",
			hostname:       "example.com",
			mockStatusCode: http.StatusOK,
			mockResponse:   "",
			wantErr:        false,
		},
		{
			name:           "Non-existent Entry",
			hostname:       "nonexistent.com",
			mockStatusCode: http.StatusNotFound,
			mockResponse:   `{"error":"not found"}`,
			wantErr:        true,
		},
		{
			name:           "Blank Hostname",
			mockStatusCode: 0,
			mockResponse:   "",
			wantErr:        true,
		},
		{
			name:           "Server Error",
			hostname:       "example.com",
			mockStatusCode: http.StatusInternalServerError,
			mockResponse:   `{"error":"Internal Server Error"}`,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected method DELETE, got %v", r.Method)
				}
				if r.URL.Path != "/dns/"+tt.hostname {
					t.Errorf("expected path /dns/%s, got %v", tt.hostname, r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			c := &Client{
				httpClient: mockServer.Client(),
				APIAddress: mockServer.URL,
				SSLVerify:  true,
			}

			err := c.DeleteDNSEntry(tt.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteDNSEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestClient_GetAllDNSEntries tests the GetAllDNSEntries method
func TestClient_GetAllDNSEntries(t *testing.T) {
	tests := []struct {
		name           string
		mockStatusCode int
		mockResponse   string
		wantErr        bool
		want           []model.DNSRecord
	}{
		{
			name:           "Successful Retrieval",
			mockStatusCode: http.StatusOK,
			mockResponse:   `[{"hostname":"example.com","ip":"192.168.1.1"}]`,
			wantErr:        false,
			want: []model.DNSRecord{
				{Hostname: "example.com", IP: "192.168.1.1"},
			},
		},
		{
			name:           "Multiple Entries",
			mockStatusCode: http.StatusOK,
			mockResponse:   `[{"hostname":"example.com","ip":"192.168.1.1"},{"hostname":"example.com","ip":"192.168.1.2"}]`,
			wantErr:        false,
			want: []model.DNSRecord{
				{Hostname: "example.com", IP: "192.168.1.1"},
				{Hostname: "example.com", IP: "192.168.1.2"},
			},
		},
		{
			name:           "Server Error",
			mockStatusCode: http.StatusInternalServerError,
			mockResponse:   `{"error":"internal server error"}`,
			wantErr:        true,
			want:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected method GET, got %v", r.Method)
				}
				if r.URL.Path != "/dns" {
					t.Errorf("expected path /dns, got %v", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			c := &Client{
				httpClient: mockServer.Client(),
				APIAddress: mockServer.URL,
				SSLVerify:  true,
			}

			entries, err := c.GetAllDNSEntries()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllDNSEntries() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, entries)
		})
	}
}

// TestClient_GetDNSEntry tests the GetDNSEntry method
func TestClient_GetDNSEntry(t *testing.T) {
	tests := []struct {
		name           string
		hostname       string
		mockStatusCode int
		mockResponse   string
		wantErr        bool
		want           []model.DNSRecord
	}{
		{
			name:           "Successful Retrieval",
			hostname:       "example.com",
			mockStatusCode: http.StatusOK,
			mockResponse:   `[{"hostname":"example.com","ip":"192.168.1.1"}]`,
			wantErr:        false,
			want: []model.DNSRecord{
				{Hostname: "example.com", IP: "192.168.1.1"},
			},
		},
		{
			name:           "Multiple Entries",
			hostname:       "example.com",
			mockStatusCode: http.StatusOK,
			mockResponse:   `[{"hostname":"example.com","ip":"192.168.1.1"},{"hostname":"example.com","ip":"192.168.1.2"}]`,
			wantErr:        false,
			want: []model.DNSRecord{
				{Hostname: "example.com", IP: "192.168.1.1"},
				{Hostname: "example.com", IP: "192.168.1.2"},
			},
		},
		{
			name:           "Blank Hostname",
			mockStatusCode: 0,
			mockResponse:   "",
			wantErr:        true,
		},
		{
			name:           "Non-existent Entry",
			hostname:       "nonexistent.com",
			mockStatusCode: http.StatusNotFound,
			mockResponse:   `{"error":"not found"}`,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected method GET, got %v", r.Method)
				}
				if r.URL.Path != "/dns/"+tt.hostname {
					t.Errorf("expected path /dns/%s, got %v", tt.hostname, r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			c := &Client{
				httpClient: mockServer.Client(),
				APIAddress: mockServer.URL,
				SSLVerify:  true,
			}

			entries, err := c.GetDNSEntry(tt.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDNSEntry() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, entries)
		})
	}
}

// TestClient_PostDNSEntry tests the PostDNSEntry method
func TestClient_PostDNSEntry(t *testing.T) {
	tests := []struct {
		name           string
		hostname       string
		ip             string
		append         bool
		mockStatusCode int
		mockResponse   string
		wantErr        bool
		want           []model.DNSRecord
	}{
		{
			name:           "Successful Creation",
			hostname:       "example.com",
			ip:             "192.168.1.1",
			append:         false,
			mockStatusCode: http.StatusOK,
			mockResponse:   `[{"hostname":"example.com","ip":"192.168.1.1"}]`,
			wantErr:        false,
			want: []model.DNSRecord{
				{Hostname: "example.com", IP: "192.168.1.1"},
			},
		},
		{
			name:           "Append Creation",
			hostname:       "example.com",
			ip:             "192.168.1.2",
			append:         true,
			mockStatusCode: http.StatusOK,
			mockResponse:   `[{"hostname":"example.com","ip":"192.168.1.1"},{"hostname":"example.com","ip":"192.168.1.2"}]`,
			wantErr:        false,
			want: []model.DNSRecord{
				{Hostname: "example.com", IP: "192.168.1.1"},
				{Hostname: "example.com", IP: "192.168.1.2"},
			},
		},
		{
			name:           "Invalid IP",
			hostname:       "example.com",
			ip:             "invalid-ip",
			append:         false,
			mockStatusCode: http.StatusBadRequest,
			mockResponse:   `{"error":"invalid IP"}`,
			wantErr:        true,
		},
		{
			name:           "Blank Hostname",
			mockStatusCode: 0,
			mockResponse:   "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected method POST, got %v", r.Method)
				}
				if r.URL.Path != "/dns/"+tt.hostname {
					t.Errorf("expected path /dns/%s, got %v", tt.hostname, r.URL.Path)
				}
				if tt.append && !strings.Contains(r.URL.RawQuery, "append=true") {
					t.Errorf("expected append=true in query, got %v", r.URL.RawQuery)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			c := &Client{
				httpClient: mockServer.Client(),
				APIAddress: mockServer.URL,
				SSLVerify:  true,
			}

			entries, err := c.PostDNSEntry(tt.hostname, tt.ip, tt.append)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostDNSEntry() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, entries)
		})
	}
}

func TestClient_decodeDNSRecords(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		contentType string
		want        []model.DNSRecord
		wantErr     bool
	}{
		{
			name:        "Successful decode",
			body:        `[{"hostname":"example.com","ip":"192.168.1.1"}]`,
			contentType: "application/json",
			want: []model.DNSRecord{
				{Hostname: "example.com", IP: "192.168.1.1"},
			},
			wantErr: false,
		},
		{
			name:        "Successful decode - multiple",
			body:        `[{"hostname":"example.com","ip":"192.168.1.1"},{"hostname":"example.com","ip":"192.168.1.2"}]`,
			contentType: "application/json",
			want: []model.DNSRecord{
				{Hostname: "example.com", IP: "192.168.1.1"},
				{Hostname: "example.com", IP: "192.168.1.2"},
			},
			wantErr: false,
		},
		{
			name:        "Failed Decode - Bad Body",
			body:        `[{"hostname":"example.com`,
			contentType: "application/json",
			want:        nil,
			wantErr:     true,
		},
		{
			name: "Failed Decode - Invalid Content Type",
			body: `---
- hostname: example.com
  ip: 192.168.1.1`,
			contentType: "application/x-yaml",
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// client won't actually do anything here
			c := &Client{}

			res := http.Response{
				Body: io.NopCloser(strings.NewReader(tt.body)),
				Header: map[string][]string{
					"Content-Type": {tt.contentType},
				},
			}

			got, err := c.decodeDNSRecords(&res)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equalf(t, tt.want, got, "decodeDNSRecords('%s')", tt.body)
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		apiURL     string
		apiPort    int
		sslVerify  bool
		wantClient *Client
		wantDiags  diag.Diagnostics
		wantErr    bool
	}{
		{
			name:      "Valid URL with port",
			apiURL:    "https://api.example.com",
			apiPort:   8080,
			sslVerify: true,
			wantClient: &Client{
				APIAddress: "https://api.example.com:8080",
				SSLVerify:  true,
			},
			wantDiags: nil,
			wantErr:   false,
		},
		{
			name:      "URL without scheme",
			apiURL:    "api.example.com",
			apiPort:   0,
			sslVerify: true,
			wantClient: &Client{
				APIAddress: "https://api.example.com",
				SSLVerify:  true,
			},
			wantDiags: diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "URL Lacks Protocol",
					Detail:   "API URL lacks 'http://' or 'https://'. Defaulting to 'https://'.\n",
				},
			},
			wantErr: false,
		},
		{
			name:       "Invalid URL",
			apiURL:     ":invalid-url",
			apiPort:    0,
			sslVerify:  true,
			wantClient: nil,
			wantDiags: diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "invalid API URL: parse \"https://:invalid-url\": invalid port \":invalid-url\" after host",
				},
			},
			wantErr: true,
		},
		{
			name:       "URL with existing port",
			apiURL:     "https://api.example.com:443",
			apiPort:    0,
			sslVerify:  true,
			wantClient: nil,
			wantDiags: diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "API URL should not contain a port. Use the separate field api_port",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, gotDiags := NewClient(tt.apiURL, tt.apiPort, tt.sslVerify)
			if !reflect.DeepEqual(gotClient, tt.wantClient) {
				t.Errorf("NewClient() gotClient = %v, want %v", gotClient, tt.wantClient)
			}
			if !reflect.DeepEqual(gotDiags, tt.wantDiags) {
				t.Errorf("NewClient() gotDiags = %v, want %v", gotDiags, tt.wantDiags)
			}
		})
	}
}

func TestClient_doRequest(t *testing.T) {
	type fields struct {
		APIAddress string
		SSLVerify  bool
	}
	type args struct {
		method string
		route  string
		body   interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Valid POST request with body",
			fields: fields{
				APIAddress: "https://api.example.com",
				SSLVerify:  true,
			},
			args: args{
				method: http.MethodPost,
				route:  "dns/entry",
				body: map[string]interface{}{
					"hostname": "example.com",
					"ip":       "192.168.1.1",
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Invalid URL in client",
			fields: fields{
				APIAddress: ":invalid-url",
				SSLVerify:  true,
			},
			args: args{
				method: http.MethodGet,
				route:  "dns/entry",
				body:   nil,
			},
			wantStatus: 0,
			wantErr:    true,
		},
		{
			name: "Failed to marshal body",
			fields: fields{
				APIAddress: "https://api.example.com",
				SSLVerify:  true,
			},
			args: args{
				method: http.MethodPost,
				route:  "dns/entry",
				body:   make(chan int),
			},
			wantStatus: 0,
			wantErr:    true,
		},
		{
			name: "Request failed with status code",
			fields: fields{
				APIAddress: "https://api.example.com",
				SSLVerify:  true,
			},
			args: args{
				method: http.MethodGet,
				route:  "dns/entry",
				body:   nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP server
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.wantStatus != http.StatusOK {
					http.Error(w, "error", tt.wantStatus)
				} else {
					w.WriteHeader(tt.wantStatus)
				}
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			// Override APIAddress with the mock server URL
			c := &Client{
				APIAddress: server.URL,
				SSLVerify:  tt.fields.SSLVerify,
			}
			if tt.fields.APIAddress == ":invalid-url" {
				c.APIAddress = tt.fields.APIAddress
			}

			got, err := c.doRequest(tt.args.method, tt.args.route, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("doRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.StatusCode != tt.wantStatus {
				t.Errorf("doRequest() got = %v, want %v", got.StatusCode, tt.wantStatus)
			}
		})
	}
}
