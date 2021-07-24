package t7Error

import (
	"net/http"
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Code    string
		Message string
		Detail  string
		Status  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				Code: "code",
				Message: "message",
				Detail: "",
			},
			want: "error code: code, message: message",
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Detail:  tt.fields.Detail,
				status:  tt.fields.Status,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_GetStatus(t *testing.T) {
	type fields struct {
		Code    string
		Message string
		Detail  string
		Status  int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "default",
			fields: fields{},
			want: http.StatusBadRequest,
		},
		{
			name: "500",
			fields: fields{
				Status: http.StatusInternalServerError,
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Detail:  tt.fields.Detail,
				status:  tt.fields.Status,
			}
			if got := e.GetStatus(); got != tt.want {
				t.Errorf("GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithDetail(t *testing.T) {
	type fields struct {
		Code    string
		Message string
		Detail  string
		Status  int
	}
	type args struct {
		d string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "have detail",
			fields: fields{},
			args: args{
				d: "detail",
			},
			want: &Error{
				Detail: "detail",
			},
		},
		{
			name: "have no detail",
			fields: fields{},
			args: args{},
			want: &Error{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Detail:  tt.fields.Detail,
				status:  tt.fields.Status,
			}
			if got := e.WithDetail(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithDetailAndStatus(t *testing.T) {
	type fields struct {
		Code    string
		Message string
		Detail  string
		Status  int
	}
	type args struct {
		d string
		s int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "empty error",
			want: &Error{},
		},
		{
			name: "normal error with default status",
			fields: fields{
				Code: "code",
				Message: "message",
			},
			args: args{
				d: "detail",
			},
			want: &Error{
				Code: "code",
				Message: "message",
				Detail: "detail",
			},
		},
		{
			name: "normal error with custom status",
			fields: fields{
				Code: "code",
				Message: "message",
			},
			args: args{
				d: "detail",
				s: http.StatusInternalServerError,
			},
			want: &Error{
				Code:    "code",
				Message: "message",
				Detail:  "detail",
				status:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Detail:  tt.fields.Detail,
				status:  tt.fields.Status,
			}
			if got := e.WithDetailAndStatus(tt.args.d, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDetailAndStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithStatus(t *testing.T) {
	type fields struct {
		Code    string
		Message string
		Detail  string
		Status  int
	}
	type args struct {
		s int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "empty error",
			want: &Error{},
		},
		{
			name: "normal error with default status",
			fields: fields{
				Code: "code",
				Message: "message",
			},
			want: &Error{
				Code: "code",
				Message: "message",
			},
		},
		{
			name: "normal error with custom status",
			fields: fields{
				Code: "code",
				Message: "message",
			},
			args: args{
				s: http.StatusInternalServerError,
			},
			want: &Error{
				Code:    "code",
				Message: "message",
				status:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Detail:  tt.fields.Detail,
				status:  tt.fields.Status,
			}
			if got := e.WithStatus(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPointerReceiver(t *testing.T) {
	testError := InvalidBody
	withDetail := testError.WithDetail("detail")
	if withDetail.Code != testError.Code || withDetail.Message != testError.Message {
		t.Errorf("message and code not match")
	}
}