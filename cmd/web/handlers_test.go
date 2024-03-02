package main

import (
	"letsgo.bepo1337/internal/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.initializeRoutes())
	defer ts.Close()

	statusCode, _, bodyString := ts.get(t, "/ping")
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, bodyString, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.initializeRoutes())
	defer ts.Close()
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "valid id",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "tale",
		},
		{
			name:     "non existing id",
			urlPath:  "/snippet/view/1337",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		code, _, body := ts.get(t, tt.urlPath)
		assert.Equal(t, code, tt.wantCode)
		if tt.wantBody != "" {
			assert.Contains(t, tt.wantBody, body)
		}
	}
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.initializeRoutes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		statusCode, header, _ := ts.get(t, "/snippet/create")
		assert.Equal(t, statusCode, http.StatusFound)
		assert.Equal(t, header.Get("location"), "/user/login")
	})

}

func extractToken(body string) string {
	return ""
}
