package adapter_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/adapter"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
	"github.com/jodi-ivan/numbered-notation-xml/utils/webserver"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	ubermock "go.uber.org/mock/gomock"
)

func TestDiagnosticHTTP_ServeHTTP(t *testing.T) {
	uberCtrl := ubermock.NewController(t)
	defer uberCtrl.Finish()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		r  *http.Request
		ps httprouter.Params
	}
	tests := []struct {
		name string // description of this test case
		args args

		initMock                   func(ctrl *ubermock.Controller) *usecase.MockUsecase
		initHTTPResponseWriterMock func(ctrl *ubermock.Controller) (*webserver.MockResponseWriteFlusher, http.Header)
		wantHeader                 http.Header
		withParam                  bool
	}{
		{
			name: "invalid scope",
			args: args{
				r:  &http.Request{},
				ps: httprouter.Params([]httprouter.Param{httprouter.Param{Key: "scope", Value: "@"}}),
			},
			initMock: func(ctrl *ubermock.Controller) *usecase.MockUsecase {
				mockUsecase := usecase.NewMockUsecase(ctrl)
				return mockUsecase
			},
			initHTTPResponseWriterMock: func(ctrl *ubermock.Controller) (*webserver.MockResponseWriteFlusher, http.Header) {
				res := webserver.NewMockResponseWriteFlusher(ctrl)
				h := http.Header{}
				res.EXPECT().WriteHeader(http.StatusBadRequest)
				res.EXPECT().Header().Return(h).Times(3)
				res.EXPECT().Write(gomock.Any())
				res.EXPECT().Flush()
				return res, h
			},
			wantHeader: map[string][]string{
				"Content-Type":  []string{"text/event-stream"},
				"Cache-Control": []string{"no-cache"},
				"Connection":    []string{"keep-alive"},
			},
		},
		{
			name: "invalid focus",
			args: args{
				r: &http.Request{
					Form: url.Values{
						"focus": {"nope"},
					},
				},
				ps: httprouter.Params([]httprouter.Param{httprouter.Param{Key: "scope", Value: "1a"}}),
			},
			initMock: func(ctrl *ubermock.Controller) *usecase.MockUsecase {
				mockUsecase := usecase.NewMockUsecase(ctrl)
				return mockUsecase
			},
			initHTTPResponseWriterMock: func(ctrl *ubermock.Controller) (*webserver.MockResponseWriteFlusher, http.Header) {
				res := webserver.NewMockResponseWriteFlusher(ctrl)
				h := http.Header{}
				res.EXPECT().WriteHeader(http.StatusBadRequest)
				res.EXPECT().Header().Return(h).Times(3)
				res.EXPECT().Write(gomock.Any())
				res.EXPECT().Flush()
				return res, h
			},
			wantHeader: map[string][]string{
				"Content-Type":  []string{"text/event-stream"},
				"Cache-Control": []string{"no-cache"},
				"Connection":    []string{"keep-alive"},
			},
		},
		{
			name: "everything went fine",
			args: args{
				r:  &http.Request{},
				ps: httprouter.Params([]httprouter.Param{httprouter.Param{Key: "scope", Value: "1a"}}),
			},
			initMock: func(ctrl *ubermock.Controller) *usecase.MockUsecase {
				mockUsecase := usecase.NewMockUsecase(ctrl)

				mockUsecase.EXPECT().RenderHymn(gomock.Any(), gomock.Any(), 1, "a")
				return mockUsecase
			},
			initHTTPResponseWriterMock: func(ctrl *ubermock.Controller) (*webserver.MockResponseWriteFlusher, http.Header) {
				res := webserver.NewMockResponseWriteFlusher(ctrl)
				h := http.Header{}
				res.EXPECT().WriteHeader(http.StatusOK)
				res.EXPECT().Header().Return(h).Times(3)
				res.EXPECT().Write(gomock.Any()).AnyTimes()
				res.EXPECT().Flush().AnyTimes()
				return res, h
			},
			wantHeader: map[string][]string{
				"Content-Type":  []string{"text/event-stream"},
				"Cache-Control": []string{"no-cache"},
				"Connection":    []string{"keep-alive"},
			},
			withParam: true,
		},
		{
			name: "timeout",
			args: args{
				r:  &http.Request{},
				ps: httprouter.Params([]httprouter.Param{httprouter.Param{Key: "scope", Value: "1a"}}),
			},
			initMock: func(ctrl *ubermock.Controller) *usecase.MockUsecase {
				mockUsecase := usecase.NewMockUsecase(ctrl)
				mockUsecase.EXPECT().RenderHymn(gomock.Any(), gomock.Any(), 1, "a")
				return mockUsecase
			},
			initHTTPResponseWriterMock: func(ctrl *ubermock.Controller) (*webserver.MockResponseWriteFlusher, http.Header) {
				res := webserver.NewMockResponseWriteFlusher(ctrl)
				h := http.Header{}
				res.EXPECT().WriteHeader(http.StatusOK)
				res.EXPECT().Header().Return(h).Times(3)
				res.EXPECT().Write(gomock.Any())
				res.EXPECT().Flush()
				return res, h
			},
			wantHeader: map[string][]string{
				"Content-Type":  []string{"text/event-stream"},
				"Cache-Control": []string{"no-cache"},
				"Connection":    []string{"keep-alive"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dh adapter.DiagnosticHTTP
			if tt.initMock != nil {
				dh.Usecase = tt.initMock(uberCtrl)
			}
			dh.Interrupt = make(chan os.Signal)
			prm := &params.Param{
				Verse: 3, // hardcoded to trigger the load verse mechanism
				Diagnostic: &params.DiagParam{
					VerseDiagnostic: make(chan params.VerseDiagnostic, 1),
					MapMtx:          &sync.Map{},
				},
			}
			w, h := tt.initHTTPResponseWriterMock(uberCtrl)
			if tt.withParam {
				adapter.NewParam = func(b bool) *params.Param {
					return prm
				}
				defer func() {
					adapter.NewParam = adapter.GetParam

				}()
			}
			r := tt.args.r
			rctx := params.NewParamContext(context.Background(), prm)

			r = tt.args.r.WithContext(rctx)
			time.AfterFunc(time.Millisecond*50, func() {
				prm.Diagnostic.VerseDiagnostic <- params.VerseDiagnostic{}
			})
			dh.ServeHTTP(w, r, tt.args.ps)

			// Assert Header match (though the mock implementation doesn't set specific headers, we keep the structure)
			assert.Equal(t, tt.wantHeader, h)

			// Verification logic to ensure RenderHymn was called with expected arguments (optional, as mocks handle expectations)
		})
	}
}
