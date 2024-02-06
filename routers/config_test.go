package routers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"bytes"
	
	model "github.com/mazeForGit/WordlistPageExtractor/model"
)
// use
// go test -v ./routers -run Test
//
func Test_ConfigGET_emptyData(t *testing.T) {
	router := gin.Default()
	router.GET("/config", ConfigGET)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/config", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("You received a %v error.", w.Code)
	}

	json.Unmarshal(w.Body.Bytes(), &model.GlobalConfig)

	if model.GlobalConfig.WordListUrl != "" {
		t.Errorf("expected empty config")
	}
}
func Test_ConfigPOST_validData(t *testing.T) {
	router := gin.Default()
	router.POST("/config", ConfigPOST)

	c := model.Config{
			RequestExecution: false, 
			WordListUrl: "test"}

	payload, err := json.Marshal(c)

	if err == nil {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/config", bytes.NewBuffer(payload))
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("You received a %v error.", w.Code)
		}

		var s model.Status
		json.Unmarshal(w.Body.Bytes(), &s)

		if s.Code != 200 && s.Text != "entity added" {
			t.Errorf("got wrong response .. code=%v, text=%s", s.Code, s.Text)
		}
	}
}
