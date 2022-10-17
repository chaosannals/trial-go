package frontend

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"trial/wstpd/basis"
	"trial/wstpd/models"
)

type FrontendDispatcher struct {
	server *basis.BizServer
}

// 初始化。
func NewFrontendDispatcher(server *basis.BizServer) (*FrontendDispatcher, error) {
	r := &FrontendDispatcher{
		server: server,
	}
	return r, nil
}

// 调度网页前端请求
func (i *FrontendDispatcher) DispatchRequest(c *gin.Context) {
	ws, err := i.server.Upgrade(c)
	if err != nil {
		log.Error(err)
		return
	}
	defer ws.Close()

	// 获取请求参数。
	request, err := i.watch(ws)
	if err != nil {
		log.Error(err)
		return
	}

	// 验证。
	cid, err := i.validate(ws, request)
	if err != nil {
		log.Error(err)
		return
	}

	// 初始化用户。
	client := i.server.FrontendManager.NewClient(cid)
	defer i.server.FrontendManager.DelClient(client)

	go func() {
		for {
			log.Infoln("转发请求：", request.Type, ".", request.Action, "来自：", cid)
			// 转发过程。
			sr := &basis.SoftwareResponse{
				Id:   i.server.NewSnowFlake(),
				Code: 4444,
				Data: request,
			}
			sclients := i.server.SoftwareManager.GetCompanyClients(cid)
			for i := range sclients {
				sclient := sclients[i]
				sclient.ResponseChannel <- sr
			}

			// 获取并验证下一个请求。
			request, err = i.watch(ws)
			if err != nil {
				log.Error(err)
				continue
			}
			cid, err = i.validate(ws, request)
			if err != nil {
				log.Error(err)
				continue
			}
			client.RequestChannel <- request
		}
	}()

	// 循环处理业务逻辑
FRLOOP:
	for {
		log.Infoln("转发请求：", request.Type, ".", request.Action, "来自：", cid)

		// 转发过程。
		sr := &basis.SoftwareResponse{
			Id:   i.server.NewSnowFlake(),
			Code: 4444,
			Data: request,
		}
		sclients := i.server.SoftwareManager.GetCompanyClients(cid)
		for i := range sclients {
			sclient := sclients[i]
			sclient.ResponseChannel <- sr
		}

		// 获取并验证下一个请求。
		request, err = i.watch(ws)
		if err != nil {
			log.Error(err)
			continue FRLOOP
		}
		cid, err = i.validate(ws, request)
		if err != nil {
			log.Error(err)
			continue FRLOOP
		}
	}
}

func (i *FrontendDispatcher) watch(ws *websocket.Conn) (*basis.FrontendRequest, error) {
	// 读入信息。
	_, msg, err := ws.ReadMessage()
	if err != nil {
		ws.WriteJSON(&basis.FrontendResponse{
			RequestId:     "--",
			ResponseError: "读取请求失败",
		})
		return nil, err
	}

	// 解析请求。
	request := &basis.FrontendRequest{}
	err = json.Unmarshal(msg, request)
	if err != nil {
		ws.WriteJSON(&basis.FrontendResponse{
			RequestId:     "--",
			ResponseError: "请求参数有误",
		})
		return nil, err
	}

	return request, nil
}

func (i *FrontendDispatcher) validate(ws *websocket.Conn, request *basis.FrontendRequest) (uint32, error) {
	// 验证请求 appkey
	db, err := models.OpenCommonDb()
	if err != nil {
		log.Errorln("数据库连接失败")
		return 0, err
	}
	cs := &models.CompanySoftware{}
	err = db.First(cs, "appkey = ?", request.Appkey).Error
	if err != nil {
		ws.WriteJSON(&basis.FrontendResponse{
			RequestId:     request.Id,
			ResponseError: "不是有效的 appkey",
		})
		return 0, err
	}

	// 验证 TOKEN 可否解密。
	token, err := basis.FromString([]byte(cs.Appsecret), request.Token)
	if err != nil {
		ws.WriteJSON(&basis.FrontendResponse{
			RequestId:     request.Id,
			ResponseError: fmt.Sprintf("token 无效 %s", request.Token),
		})
		return 0, err
	}

	// 验证 TOKEN 的公司 ID 。
	if token.CompanyId != cs.CompanyId {
		ws.WriteJSON(&basis.FrontendResponse{
			RequestId:     request.Id,
			ResponseError: "token id 无效",
		})
		return 0, errors.New("token id 无效")
	}

	// 验证 TOKEN 有效期。
	et := token.CreateAt.Add(2 * time.Hour)
	nt := time.Now()
	if et.Before(nt) {
		ermsg := fmt.Sprintf("token 失效, %v %v", et, nt)
		ws.WriteJSON(&basis.FrontendResponse{
			RequestId:     request.Id,
			ResponseError: ermsg,
		})
		return token.CompanyId, errors.New(ermsg)
	}
	return token.CompanyId, nil
}
