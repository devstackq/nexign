package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/devstackq/nexign/config"
	"github.com/devstackq/nexign/entity"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const url = "https://api.textgears.com/spelling"

type handler struct {
	lg  *zap.Logger
	cfg *config.Config
	// usc
}

// type Speller interface {
// 	spell()
// }

type textgears struct {
	Text []*string `json:"texts"`
}

func New(lg *zap.Logger, cfg *config.Config) *handler {
	return &handler{lg: lg, cfg: cfg}
}

func (c *handler) hSpeller(g *gin.Context) {
	var (
		err error
		tg  = &textgears{}
	)
	if err = g.ShouldBindJSON(tg); err != nil {
		c.lg.Info(err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	tg.concurrentRequest(c.cfg)

	g.JSON(http.StatusOK, tg.Text)
}

func (s *textgears) spelling(cfg *config.Config, sentence *string, result chan entity.Result, errCh chan error) {
	var (
		response = entity.Result{}
		url      string
		err      error
	)
	url = fmt.Sprint(cfg.SpellerCfg.Url, "/spelling?key=", cfg.Key, "&language=", cfg.Language, "&text=", *sentence)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		errCh <- err
		log.Println(err)
		return
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errCh <- err
		log.Println(err)
		return
	}

	err = json.Unmarshal(bb, &response)
	if err != nil {
		errCh <- err
		log.Println(err)
		return
	}
	errCh <- nil
	result <- response
}

// register handler
func (c *handler) Route() http.Handler {
	r := gin.New()
	r.POST("/speller", c.hSpeller)
	return r
}
