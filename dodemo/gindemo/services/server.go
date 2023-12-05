package services

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type GinService struct {
	server     *http.Server
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
	helloClient := do.MustInvoke[*ApisHelloClientService](i)

	router.Use(func(ctx *gin.Context) {
		if ctx.Request.ProtoMajor == 2 && strings.HasPrefix(ctx.GetHeader("Content-Type"), "application/grpc") {
			ctx.Status(http.StatusOK)
			grpcServer.ServeHTTP(ctx.Writer, ctx.Request)
			ctx.Abort()
			return
		}
		ctx.Next()
	})

	router.Any("/hello", func(ctx *gin.Context) {

		r, err := helloClient.SayHello("aaaa")

		if err != nil {
			ctx.Status(http.StatusBadRequest)
			ctx.Writer.WriteString(err.Error())
			ctx.Abort()
			return
		}
		ctx.Writer.WriteString(r.Message)

		// ctx.Writer.WriteString("hello")
		ctx.Status(http.StatusOK)
		ctx.Abort()
	})

	address := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	server := &http.Server{
		Addr:           address,
		Handler:        h2c.NewHandler(router.Handler(), &http2.Server{}),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1000000,
	}

	return &GinService{
		server:     server,
		router:     router,
		logger:     logger,
		conf:       conf,
		grpcServer: grpcServer,
	}, nil
}

func (i *GinService) Run() {
	//address := fmt.Sprintf("%s:%d", i.conf.Host, i.conf.Port)

	//i.router.Run(address)
	// router.RunTLS(address, "certFile", "keyFile")

	if err := i.server.ListenAndServe(); err == nil {
		i.logger.Error("serve failed", "error", err)
	}

	// if err := i.server.ListenAndServeTLS("certFile", "keyFile"); err == nil {
	// 	i.logger.Error("serve failed", "error", err)
	// }
}
