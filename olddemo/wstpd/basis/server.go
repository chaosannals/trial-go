package basis

import (
	"net/http"

	"github.com/bwmarrin/snowflake"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 业务服务器
type BizServer struct {
	upgrader        websocket.Upgrader
	engine          *gin.Engine
	snowflakeNode   *snowflake.Node
	FrontendManager *FrontendClientManager
	SoftwareManager *SoftwareClientManager
}

// 创建业务服务器
func NewBizServer() (*BizServer, error) {
	snowflakeNode, err := snowflake.NewNode(10)
	if err != nil {
		return nil, err
	}

	frontendManager, err := NewFrontendClientManager()
	if err != nil {
		return nil, err
	}

	softwareManager, err := NewSoftwareClientManager()
	if err != nil {
		return nil, err
	}

	return &BizServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		engine:          gin.Default(),
		snowflakeNode:   snowflakeNode,
		FrontendManager: frontendManager,
		SoftwareManager: softwareManager,
	}, nil
}

func (i *BizServer) NewSnowFlake() int64 {
	return i.snowflakeNode.Generate().Int64()
}

// 挂载路径处理器
func (i *BizServer) Attach(rpath string, handler gin.HandlerFunc) {
	i.engine.GET(rpath, handler)
}

// HTTP 协变 WebSocket
func (i *BizServer) Upgrade(c *gin.Context) (*websocket.Conn, error) {
	return i.upgrader.Upgrade(c.Writer, c.Request, nil)
}

// 启动服务器
func (i *BizServer) Run() {
	log.Infoln("服务器启动")
	i.engine.Run(":44444")
}
