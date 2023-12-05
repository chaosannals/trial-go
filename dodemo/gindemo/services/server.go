package services

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"google.golang.org/grpc"
)

type GinService struct {
	router     *gin.Engine
	logger     *slog.Logger
	conf       *ConfService
	grpcServer *grpc.Server
}

func NewGinService(i *do.Injector) (*GinService, error) {
	router := gin.New()
	logger := do.MustInvoke[*slog.Logger](i)
	conf := do.MustInvoke[*ConfService](i)
	grpcServer := do.MustInvoke[*grpc.Server](i)

	router.Use(func(ctx *gin.Context) {
		if ctx.Request.ProtoMajor == 2 && strings.HasPrefix(ctx.GetHeader("Content-Type"), "application/grpc") {
			ctx.Status(http.StatusOK)
			grpcServer.ServeHTTP(ctx.Writer, ctx.Request)
			ctx.Abort()
			return
		}
		ctx.Next()
	})

	return &GinService{
		router:     router,
		logger:     logger,
		conf:       conf,
		grpcServer: grpcServer,
	}, nil
}

func (i *GinService) Run() {
	address := fmt.Sprintf("%s:%d", i.conf.Host, i.conf.Port)

	i.router.Run(address)

	// router.RunTLS(address)
}
