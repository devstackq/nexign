package handler

import (
	"net/http"

	"github.com/devstackq/nexign/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handler struct {
	lg  *zap.Logger
	cfg *config.Config
	// usc
}

func New(lg *zap.Logger, cfg *config.Config) *handler {
	return &handler{lg: lg, cfg: cfg}
}

func (c *handler) Route() http.Handler {
	r := gin.New()
	r.POST("/speller", c.hSpeller)
	return r
}
