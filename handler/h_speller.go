package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/devstackq/nexign/config"
	"github.com/devstackq/nexign/entity"
	"github.com/gin-gonic/gin"
)

type textgears struct {
	Sentences []string `json:"texts"`
	text      *string
}

type Resp struct {
	msg  string
	code int
}

func (c *handler) hSpeller(g *gin.Context) {
	var (
		err error
		tg  = &textgears{}
	)

	if err = g.ShouldBindJSON(tg); err != nil {
		c.lg.Info(err.Error())
		statusWithError(g, http.StatusBadRequest, err.Error())
		return
	}

	tg.convertToStr()

	if err = tg.spelling(c.cfg); err != nil {
		c.lg.Info(err.Error())
		statusWithError(g, http.StatusInternalServerError, err.Error())
		return
	}

	g.JSON(http.StatusOK, tg.Sentences)
}

func (tg *textgears) convertToStr() {
	str := ""
	for _, text := range tg.Sentences {
		str += text + "|"
	}
	tg.text = &str
}

func (tg *textgears) spelling(cfg *config.Config) error {
	var (
		response = entity.Result{}
		url      string
		err      error
	)
	url = fmt.Sprint(cfg.SpellerCfg.Url, "/spelling?key=", cfg.Key, "&language=", cfg.Language, "&text=", *tg.text)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bb, &response)
	if err != nil {
		return err
	}
	temp := ""

	if len(response.Errors) > 0 {
		for _, word := range response.Errors {
			temp = strings.Replace(*tg.text, word.Bad, word.Better[0], -1)
			tg.text = &temp
		}
	}

	tempList := strings.Split(*tg.text, "|")
	tg.Sentences = tempList
	return nil
}

func statusWithError(g *gin.Context, statusCode int, msg string) {
	resp := Resp{}
	resp.code = statusCode
	resp.msg = msg
	g.AbortWithStatusJSON(http.StatusBadRequest, resp)
}
