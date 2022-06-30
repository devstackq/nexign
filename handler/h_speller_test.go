package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devstackq/nexign/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_hSpeller(t *testing.T) {
	testCases := []struct {
		name         string
		expectError  error
		expectStatus int
		url          string
		cfg          *config.Config
		texts        textgears
	}{
		{
			name:         "ok",
			texts:        textgears{Sentences: []string{"всем привет, как дела?", "что на счет завтрака?"}},
			expectStatus: http.StatusOK,

			cfg: &config.Config{
				Port: ":6969",
				SpellerCfg: config.SpellerCfg{
					Key:      "6gHksNL9Zjck2zEx",
					Url:      "https://api.textgears.com",
					Language: "ru-RU",
				},
			},
			expectError: nil,
			url:         "http://localhost:6969/speller",
		},
	}

	var assertion = assert.New(t)
	gin.SetMode(gin.TestMode)

	lg, err := zap.NewProduction()
	if err != nil {
		t.Error(err)
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			respRec := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(respRec)
			hr := New(lg, tc.cfg)

			router.POST("/speller", hr.hSpeller)

			bytesBuffer := &bytes.Buffer{}

			err = json.NewEncoder(bytesBuffer).Encode(tc.texts)

			ctx.Request, err = http.NewRequest("POST", tc.url, bytesBuffer)

			tg := textgears{}
			tg.convertToStr()
			tg.spelling(tc.cfg)

			router.ServeHTTP(respRec, ctx.Request)
			assertion.Equal(tc.expectError, err)
			assertion.EqualValues(tc.expectStatus, respRec.Code)
		})
	}
}
