package adapter

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
	"github.com/jodi-ivan/numbered-notation-xml/utils/webserver"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRenderHTTP_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		r  *http.Request
		ps httprouter.Params
	}
	tests := []struct {
		name string
		args args

		initMock                   func(ctrl *gomock.Controller) *usecase.MockUsecase
		initHTTPResponseWriterMock func(ctrl *gomock.Controller) *webserver.MockResponseWriter
	}{
		{
			name: "invalid parameter",
			initHTTPResponseWriterMock: func(ctrl *gomock.Controller) *webserver.MockResponseWriter {
				res := webserver.NewMockResponseWriter(ctrl)
				res.EXPECT().WriteHeader(http.StatusBadRequest)
				res.EXPECT().Write([]byte("Invalid URL"))
				return res
			},
		},
		{
			name: "everything went fine",
			args: args{
				ps: httprouter.Params([]httprouter.Param{
					httprouter.Param{
						Key:   "number",
						Value: "1",
					},
				}),
				r: &http.Request{},
			},
			initHTTPResponseWriterMock: func(ctrl *gomock.Controller) *webserver.MockResponseWriter {
				res := webserver.NewMockResponseWriter(ctrl)
				return res
			},
			initMock: func(ctrl *gomock.Controller) *usecase.MockUsecase {
				res := usecase.NewMockUsecase(ctrl)
				res.EXPECT().RenderHymn(gomock.Any(), gomock.Any(), int(1))
				return res
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &RenderHTTP{}
			if tt.initMock != nil {
				rh.usecase = tt.initMock(ctrl)
			}

			w := tt.initHTTPResponseWriterMock(ctrl)
			rh.ServeHTTP(w, tt.args.r, tt.args.ps)
		})
	}
}

func TestCanvasDelegatorHTTP_OnBeforeStartWrite(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                       string
		initHTTPResponseWriterMock func(ctrl *gomock.Controller) (*webserver.MockResponseWriter, http.Header)
	}{
		{
			name: "default",
			initHTTPResponseWriterMock: func(ctrl *gomock.Controller) (*webserver.MockResponseWriter, http.Header) {
				res := webserver.NewMockResponseWriter(ctrl)
				res.EXPECT().WriteHeader(http.StatusOK)
				header := http.Header{}
				res.EXPECT().Header().Return(header)
				return res, header
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := tt.initHTTPResponseWriterMock(ctrl)
			cdh := CanvasDelegatorHTTP{
				w: w,
			}
			cdh.OnBeforeStartWrite()

			if !assert.Equal(t, http.Header(map[string][]string{
				"Content-Type": []string{"image/svg+xml"},
			}), h) {
				t.Fail()
			}
		})
	}
}

func TestCanvasDelegatorHTTP_OnError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want canvas.DelegatorErrorFlowControl

		initHTTPResponseWriterMock func(ctrl *gomock.Controller) *webserver.MockResponseWriter
	}{
		{
			name: "error hymn not found",
			args: args{
				err: repository.ErrHymnNotFound,
			},
			initHTTPResponseWriterMock: func(ctrl *gomock.Controller) *webserver.MockResponseWriter {
				res := webserver.NewMockResponseWriter(ctrl)
				return res
			},
			want: canvas.DelegatorErrorFlowControlIgnore,
		},
		{
			name: "other error",
			args: args{
				err: errors.New("nope"),
			},
			initHTTPResponseWriterMock: func(ctrl *gomock.Controller) *webserver.MockResponseWriter {
				res := webserver.NewMockResponseWriter(ctrl)
				res.EXPECT().WriteHeader(http.StatusInternalServerError)
				res.EXPECT().Write([]byte("nope"))
				return res
			},
			want: canvas.DelegatorErrorFlowControlStop,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.initHTTPResponseWriterMock(ctrl)

			cdh := CanvasDelegatorHTTP{
				w: w,
			}
			if got := cdh.OnError(tt.args.err); !assert.Equal(t, tt.want, got) {
				t.Errorf("CanvasDelegatorHTTP.OnError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		u usecase.Usecase
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "default",
			args: args{
				u: usecase.New(config.Config{}, nil, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := New(tt.args.u)

			if !assert.NotNil(t, res) {
				t.Fail()
			}
		})
	}
}
