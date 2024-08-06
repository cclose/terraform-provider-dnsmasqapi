package validation

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestValidateHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname interface{}
		wantErr  bool
	}{
		{"Valid Hostname", "example.com", false},
		{"Valid Subdomain Hostname", "sub.example.com", false},
		{"TooLongHostname", strings.Repeat("a", 254), true},
		{"StartingHyphen", "-example.com", true},
		{"Double Dot", "example..com-", true},
		{"Bad Type", make(chan int), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHostname(tt.hostname, cty.Path{})
			assert.Equal(t, tt.wantErr, err != nil, "ValidateHostname() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateIPAddress(t *testing.T) {
	tests := []struct {
		name    string
		ip      interface{}
		wantErr bool
	}{
		{"Valid IP", "192.168.0.1", false},
		{"Invalid IP", "1942.168.0.1", true},
		{"Total Gibberish", "fishface", true},
		{"Bad Type", make(chan int), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIPAddress(tt.ip, cty.Path{})
			assert.Equal(t, tt.wantErr, err != nil, "ValidateIPAddress() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
