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
	type args struct {
		w       http.ResponseWriter
		status  int
		data    any
		headers []http.Header
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
		wantBody   string
	}{
		{
			name: "test #1 - valid json",
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
			app := &Config{}

			err := app.writeJSON(tt.args.w, tt.args.status, tt.args.data, tt.args.headers...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Config.writeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			rr := tt.args.w.(*httptest.ResponseRecorder)

			if rr.Code != tt.wantStatus {
				t.Errorf("Config.errorJSON() status code = %v, want %v", rr.Code, tt.wantStatus)
			}

			gotBody := strings.TrimSpace(rr.Body.String())
			if gotBody != tt.wantBody {
				t.Errorf("Config.errorJSON() body = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestConfig_readJSON(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		data any
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
		// wantBody string
	}{
		{
			name: "test #1 - valid json",
			args: args{
				w:    httptest.NewRecorder(),
				r:    httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"key":"value"}`)),
				data: &map[string]string{},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
			// wantBody:   `{"key":"value"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Config{}

			err := app.readJSON(tt.args.w, tt.args.r, tt.args.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Config.readJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			rr := tt.args.w.(*httptest.ResponseRecorder)

			if rr.Code != tt.wantStatus {
				t.Errorf("Config.errorJSON() status code = %v, want %v", rr.Code, tt.wantStatus)
			}

			// gotBody := strings.TrimSpace(rr.Body.String())
			// if gotBody != tt.wantBody {
			// 	t.Errorf("Config.errorJSON() body = %v, want %v", gotBody, tt.wantBody)
			// }
		})
	}
}

func TestConfig_errorJSON(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		err    error
		status []int
	}
	tests := []struct {
		name       string
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
			wantBody:   `{"error":true,"message":"error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Config{}

			err := app.errorJSON(tt.args.w, tt.args.err, tt.args.status...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Config.errorJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			rr := tt.args.w.(*httptest.ResponseRecorder)

			if rr.Code != tt.wantStatus {
				t.Errorf("Config.errorJSON() status code = %v, want %v", rr.Code, tt.wantStatus)
			}

			gotBody := strings.TrimSpace(rr.Body.String())
			if gotBody != tt.wantBody {
				t.Errorf("Config.errorJSON() body = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestConfig_checkRequest(t *testing.T) {
	type args struct {
		req requestQuote
	}
	tests := []struct {
		name     string
		args     args
		wantArgs []string
	}{
		{
			name: "test #1 - invalid zipcode",
			args: args{
				req: requestQuote{
					Recipient: recipientQuote{
						Address: address{
							Zipcode: "",
						},
					},
					Volumes: []volume{
						{
							Category: 1,
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
					Recipient: recipientQuote{
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
			app := &Config{}

			gotArgs := app.checkRequest(tt.args.req)

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("Config.checkRequest() = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestConfig_buildRequestAPI(t *testing.T) {
	type args struct {
		reqQuote requestQuote
	}
	tests := []struct {
		name       string
		args       args
		wantReqAPI requestAPI
	}{
		{
			name: "test #1 - valid request",
			args: args{
				reqQuote: requestQuote{
					Recipient: recipientQuote{
						Address: address{
							Zipcode: "12345",
						},
					},
					Volumes: []volume{
						{
							Category: 1,
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
			wantReqAPI: requestAPI{
				Recipient: recipientApi{
					Type:    0,
					Country: "BRA",
					Zipcode: 12345,
				},
				Dispatchers: []dispatcher{
					{
						RegisteredNumber: "",
						Zipcode:          12345,
						Volumes: []volumeApi{
							{
								Category:      "1",
								Amount:        1,
								UnitaryWeight: 0,
								Price:         100.0,
								UnitaryPrice:  100,
								Sku:           "SKU123",
								Height:        10.0,
								Width:         5.0,
								Length:        20.0,
							},
						},
					},
				},
				SimulationType: []int{0},
				Returns: returns{
					Composition:  false,
					Volumes:      false,
					AppliedRules: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Config{}

			gotReqAPI := app.buildRequestAPI(tt.args.reqQuote)

			if !reflect.DeepEqual(gotReqAPI, tt.wantReqAPI) {
				t.Errorf("Config.buildRequestAPI() = %v, want %v", gotReqAPI, tt.wantReqAPI)
			}
		})
	}
}

func TestConfig_formatResponseAPI(t *testing.T) {
	type args struct {
		entry responseAPI
	}
	tests := []struct {
		name       string
		args       args
		wantResult data.QuoteEntry
	}{
		{
			name: "test #1 - valid json",
			args: args{
				entry: responseAPI{
					Dispatchers: []dispatcherAPI{
						{
							ID: "test",
							Offers: []offer{
								{
									Modal:      "test",
									FinalPrice: 1.5,
									Carrier: carrier{
										Name: "test",
									},
									CarrierOriginalDeliveryTime: carrierOriginalDeliveryTime{
										Days: 1,
									},
								},
							},
						},
					},
				},
			},
			wantResult: data.QuoteEntry{
				Carrier: []data.Carrier{
					{
						Name:     "test",
						Service:  "test",
						Deadline: 1,
						Price:    1.5,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Config{}

			gotResult := app.formatResponseAPI(tt.args.entry)

			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Config.formatResponseAPI() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
