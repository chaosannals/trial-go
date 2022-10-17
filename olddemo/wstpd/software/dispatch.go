package software

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"trial/wstpd/basis"
	"trial/wstpd/models"
)

// 调度器
type SoftwareDispatcher struct {
	server *basis.BizServer
}

// 调度
func NewSoftwareDispatcher(server *basis.BizServer) (*SoftwareDispatcher, error) {
	r := &SoftwareDispatcher{
		server: server,
	}
	return r, nil
}

// 调度软件请求
func (i *SoftwareDispatcher) DispatchRequest(c *gin.Context) {
	ws, err := i.server.Upgrade(c)
	if err != nil {
		log.Error(err)
		return
	}
	defer ws.Close()

	// 登录验证
	cid, err := i.login(ws)
	if err != nil {
		log.Error(err)
		return
	}
	client := i.server.SoftwareManager.NewClient(ws, cid)
	defer i.server.SoftwareManager.DelClient(client)

	logic := &SoftwareLogic{}
	ltv := reflect.ValueOf(logic)

	// 读取请求
	go func() {
		for {
			request, err := i.watch(ws)
			if err != nil {
				log.Error(err)
				client.SetAble(false)
				client.CloseChannel <- true
				break
			}
			client.RequestChannel <- request
		}
	}()

	// 业务请求循环处理。
	// 判断是否断开。。。。TODO
RLOOP:
	for client.IsAble() {
		select {
		case request := <-client.RequestChannel:
			if request.Type == "frontend" {
				// 服务 直发前端
				fcs := i.server.FrontendManager.GetCompanyClients(client.CompanyId)
				for i := range fcs {
					fc := fcs[i]
					fc.ResponseChannel <- request.Data.(*basis.FrontendResponse)
				}
			} else {
				// 处理功能性 服务请求 调度。
				m := ltv.MethodByName(request.Type)
				if !m.IsValid() {
					log.Error(request)
					ws.WriteJSON(&basis.SoftwareResponse{
						Id:   request.Id,
						Code: -1,
						Data: "没有该操作类型",
					})
					continue RLOOP
				}
				ma := make([]reflect.Value, 1)
				ma[0] = reflect.ValueOf(request)
				rs := m.Call(ma)
				log.Infoln("调用：", request.Type, " 成功\n\r", rs)
				response := rs[0].Interface()
				ws.WriteJSON(response)
			}
		case response := <-client.ResponseChannel:
			// 前端直接转发的请求，直接响应给服务。
			ws.WriteJSON(response)
		case request := <-client.CloseChannel:
			log.Infoln("关闭请求", request)
			if request {
				break RLOOP
			}
		}
	}
}

// 获取请求参数。
func (i *SoftwareDispatcher) watch(ws *websocket.Conn) (*basis.SoftwareRequest, error) {
	// 获取登录信息。
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Error(err)
		ws.WriteJSON(&basis.SoftwareResponse{
			Id:   0,
			Code: -1,
			Data: "读取请求失败",
		})
		return nil, err
	}
	request := &basis.SoftwareRequest{}
	err = json.Unmarshal(msg, request)
	if err != nil {
		log.Warning(err)
		ws.WriteJSON(&basis.SoftwareResponse{
			Id:   0,
			Code: -1,
			Data: "请求参数有误",
		})
	}
	return request, nil
}

// 登录软件账户。
func (i *SoftwareDispatcher) login(ws *websocket.Conn) (uint32, error) {
	// 获取请求参数。
	loginRequest, err := i.watch(ws)
	if err != nil {
		return 0, err
	}

	// 验证登录类型。
	if loginRequest.Type != "login" {
		ws.WriteJSON(&basis.SoftwareResponse{
			Id:   0,
			Code: -1,
			Data: "第一个请求必须是 login 类型",
		})
	}

	// 获取查询信息。
	db, err := models.OpenCommonDb()
	if err != nil {
		return 0, err
	}
	cs := &models.CompanySoftware{}
	ld := loginRequest.Data.(map[string]interface{})
	err = db.First(cs, "appkey = ?", ld["appkey"].(string)).Error
	if err != nil {
		ws.WriteJSON(&basis.SoftwareResponse{
			Id:   loginRequest.Id,
			Code: -1,
			Data: "不是有效的 appkey",
		})
		return 0, err
	}

	// 验证 TOKEN 可否解密。
	rtoken := ld["token"].(string)
	token, err := basis.FromString([]byte(cs.Appsecret), rtoken)
	if err != nil {
		ws.WriteJSON(&basis.SoftwareResponse{
			Id:   loginRequest.Id,
			Code: -1,
			Data: fmt.Sprintf("token 无效 %s", rtoken),
		})
		return 0, err
	}

	// 验证 TOKEN 的公司 ID 。
	if token.CompanyId != cs.CompanyId {
		ws.WriteJSON(&basis.SoftwareResponse{
			Id:   loginRequest.Id,
			Code: -1,
			Data: "token id 无效",
		})
		return token.CompanyId, errors.New("token id 无效")
	}

	// 验证 TOKEN 有效期。
	et := token.CreateAt.Add(2 * time.Hour)
	nt := time.Now()
	if et.Before(nt) {
		ermsg := fmt.Sprintf("token 失效, %v %v", et, nt)
		ws.WriteJSON(&basis.SoftwareResponse{
			Id:   loginRequest.Id,
			Code: -1,
			Data: ermsg,
		})
		return token.CompanyId, errors.New(ermsg)
	}

	// 登录成功。
	ws.WriteJSON(&basis.SoftwareResponse{
		Id:   loginRequest.Id,
		Code: 0,
		Data: "登录成功",
	})
	log.Infoln(cs.CompanyId, " 登录")
	return token.CompanyId, nil
}
