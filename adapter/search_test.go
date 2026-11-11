package adapter_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/adapter"
	"github.com/jodi-ivan/numbered-notation-xml/svc/entity"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/webserver"
	"github.com/julienschmidt/httprouter"
	ubermock "go.uber.org/mock/gomock"
)

func TestSearchAdapterHTTP_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uberCtrl := ubermock.NewController(t)
	defer uberCtrl.Finish()
	type args struct {
		r  *http.Request
		ps httprouter.Params
	}
	tests := []struct {
		name string // description of this test case
		args args

		initMock                   func(ctrl *ubermock.Controller) *usecase.MockUsecase
		initHTTPResponseWriterMock func(ctrl *gomock.Controller) *webserver.MockResponseWriter
		wantStatusCode             int    // desired HTTP status code
		wantBody                   string // desired response body content (optional)
	}{
		{
			name: "successful search",
			args: args{
				r:  &http.Request{URL: &url.URL{RawQuery: "q=test"}},
				ps: httprouter.Params([]httprouter.Param{httprouter.Param{Key: "number", Value: "1"}}),
			},
			initMock: func(ctrl *ubermock.Controller) *usecase.MockUsecase {
				mockUsecase := usecase.NewMockUsecase(ctrl)
				mockUsecase.EXPECT().Lookup(gomock.Any(), "test").Return(
					[]entity.SearchMatch{entity.SearchMatch{HymnNo: 1, Title: "Test Title", Content: "Test Content"}},
					nil,
				).AnyTimes()
				mockUsecase.EXPECT().Highlight(gomock.Any(), gomock.Any()).Return(
					[]usecase.HighlightSearchResult{{FullTitle: "Test Title"}},
				).AnyTimes()
				return mockUsecase
			},
			initHTTPResponseWriterMock: func(ctrl *gomock.Controller) *webserver.MockResponseWriter {
				res := webserver.NewMockResponseWriter(ctrl)
				header := http.Header{}
				res.EXPECT().Header().Return(header) //.Set("Content-Type", "application/json")

				res.EXPECT().WriteHeader(200)
				res.EXPECT().Write([]byte(`[{"title":{"snippet":"","spans":null},"full_title":"Test Title","tier":0}]`))
				return res
			},
		},
		{
			name: "lookup error",
			args: args{
				r:  &http.Request{URL: &url.URL{RawQuery: "q=error"}},
				ps: httprouter.Params([]httprouter.Param{httprouter.Param{Key: "number", Value: "1"}}),
			},
			initMock: func(ctrl *ubermock.Controller) *usecase.MockUsecase {
				mockUsecase := usecase.NewMockUsecase(ctrl)
				mockUsecase.EXPECT().Lookup(gomock.Any(), "error").Return(nil, errors.New("lookup failed")).AnyTimes()
				return mockUsecase
			},
			initHTTPResponseWriterMock: func(ctrl *gomock.Controller) *webserver.MockResponseWriter {
				res := webserver.NewMockResponseWriter(ctrl)
				res.EXPECT().WriteHeader(500)
				res.EXPECT().Write([]byte("Failed to search."))
				return res
			},
		},
		{
			name: "missing query parameter",
			args: args{
				r:  &http.Request{URL: &url.URL{RawQuery: ""}},
				ps: httprouter.Params([]httprouter.Param{httprouter.Param{Key: "number", Value: "1"}}),
			},
			initMock: func(ctrl *ubermock.Controller) *usecase.MockUsecase {
				return usecase.NewMockUsecase(ctrl)
			},
			initHTTPResponseWriterMock: func(ctrl *gomock.Controller) *webserver.MockResponseWriter {
				res := webserver.NewMockResponseWriter(ctrl)
				res.EXPECT().WriteHeader(400)
				res.EXPECT().Write([]byte("Query can't be empty"))

				return res
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sah := adapter.NewSearchAdapter(nil)
			if tt.initMock != nil {
				sah.Usecase = tt.initMock(uberCtrl)
			}

			w := tt.initHTTPResponseWriterMock(ctrl)
			sah.ServeHTTP(w, tt.args.r, tt.args.ps)
		})
	}
}
