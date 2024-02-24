package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gin-gonic/gin"
)

type IGqlHandler interface {
	GqlHdl() gin.HandlerFunc
}

type GqlHandler struct {
	hdlServer *handler.Server
}

func NewGqlHandler(hdlServer *handler.Server) IGqlHandler {
	return &GqlHandler{
		hdlServer: hdlServer,
	}
}

func (g *GqlHandler) GqlHdl() gin.HandlerFunc {
	return func(c *gin.Context) {
		g.hdlServer.AddTransport(transport.Options{})
		g.hdlServer.AddTransport(transport.GET{})
		g.hdlServer.AddTransport(transport.POST{})
		g.hdlServer.AddTransport(transport.MultipartForm{})
		g.hdlServer.ServeHTTP(c.Writer, c.Request)
	}
}
