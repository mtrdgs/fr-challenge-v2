package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mtrdgs/fr/data"
)

func TestConfig_writeJSON(t *testing.T) {
	type fields struct {
		Models data.Models
	}
	type args struct {
		w       http.ResponseWriter
		status  int
		data    any
		headers []http.Header
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantStatus int
		wantBody   string
	}{
		{
			name:   "test #1 - valid json",
			fields: fields{Models: data.Models{}},
			args: args{
				w:       httptest.NewRecorder(),
				status:  http.StatusOK,
				data:    map[string]string{"message": "success!"},
				headers: nil,
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"success!"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Config{
				Models: tt.fields.Models,
			}
			if err := app.writeJSON(tt.args.w, tt.args.status, tt.args.data, tt.args.headers...); (err != nil) != tt.wantErr {
				t.Errorf("Config.writeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_readJSON(t *testing.T) {
	type fields struct {
		Models data.Models
	}
	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		data any
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		wantBody string
	}{
		{
			name:   "test #1 - valid json",
			fields: fields{Models: data.Models{}},
			args: args{
				w:    httptest.NewRecorder(),
				r:    httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"key":"value"}`)),
				data: &map[string]string{},
			},
			wantErr:  false,
			wantBody: `{"key":"value"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Config{
				Models: tt.fields.Models,
			}
			if err := app.readJSON(tt.args.w, tt.args.r, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Config.readJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_errorJSON(t *testing.T) {
	type fields struct {
		Models data.Models
	}
	type args struct {
		w      http.ResponseWriter
		err    error
		status []int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantStatus int
		wantBody   string
	}{
		{
			name: "test #1 - valid error",
			args: args{
				w:      httptest.NewRecorder(),
				err:    errors.New("error"),
				status: []int{http.StatusBadGateway},
			},
			wantErr:    false,
			wantStatus: http.StatusBadGateway,
			wantBody:   `{"error":true,"message":"something went wrong"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Config{
				Models: tt.fields.Models,
			}
			if err := app.errorJSON(tt.args.w, tt.args.err, tt.args.status...); (err != nil) != tt.wantErr {
				t.Errorf("Config.errorJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_checkRequest(t *testing.T) {
	type fields struct {
		Models data.Models
	}
	type args struct {
		req requestQuote
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantArgs []string
	}{
		{
			name: "test #1 - invalid zipcode",
			args: args{
				req: requestQuote{
					Recipient: recipient{
						Address: address{
							Zipcode: "",
						},
					},
					Volumes: []volume{
						{
							Category: "electronics",
							Amount:   1,
							Price:    100.0,
							Sku:      "SKU123",
							Height:   10.0,
							Width:    5.0,
							Length:   20.0,
						},
					},
				},
			},
			wantArgs: []string{"Zipcode is required"},
		},
		{
			name: "test #2 - invalid volumes",
			args: args{
				req: requestQuote{
					Recipient: recipient{
						Address: address{
							Zipcode: "1234",
						},
					},
					Volumes: []volume{},
				},
			},
			wantArgs: []string{"Volumes are required"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Config{
				Models: tt.fields.Models,
			}
			if gotArgs := app.checkRequest(tt.args.req); !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("Config.checkRequest() = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
