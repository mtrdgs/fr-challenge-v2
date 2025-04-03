package main

import (
	"net/http"
	"net/http/httptest"
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
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test #1 - valid json",
			args: args{
				w:      httptest.NewRecorder(),
				err:    nil,
				status: []int{http.StatusAccepted},
			},
			wantErr: false,
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
