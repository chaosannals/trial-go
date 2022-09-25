package basis

import (
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
)

// 软件客户端
type SoftwareClient struct {
	Id              int64
	CompanyId       uint32
	Socket          *websocket.Conn
	RequestChannel  chan *SoftwareRequest
	ResponseChannel chan *SoftwareResponse
	CloseChannel    chan bool
	Able            bool
	AbleMutex       *sync.RWMutex
}

func (i *SoftwareClient) SetAble(v bool) {
	i.AbleMutex.Lock()
	i.Able = v
	i.AbleMutex.Unlock()
}

func (i *SoftwareClient) IsAble() bool {
	i.AbleMutex.RLock()
	result := i.Able
	i.AbleMutex.RUnlock()
	return result
}

// 软件客户端管理器。
type SoftwareClientManager struct {
	clients       map[uint32]map[int64]*SoftwareClient
	snowflakeNode *snowflake.Node
	rwm           *sync.RWMutex
}

// 创建一个软件客户端管理器
func NewSoftwareClientManager() (*SoftwareClientManager, error) {
	snode, err := snowflake.NewNode(2)
	if err != nil {
		return nil, err
	}
	return &SoftwareClientManager{
		clients:       make(map[uint32]map[int64]*SoftwareClient),
		snowflakeNode: snode,
		rwm:           new(sync.RWMutex),
	}, nil
}

// 创建一个软件客户端
func (i *SoftwareClientManager) NewClient(ws *websocket.Conn, companyId uint32) *SoftwareClient {
	i.rwm.Lock()
	defer i.rwm.Unlock()
	c := &SoftwareClient{
		Id:              i.snowflakeNode.Generate().Int64(),
		CompanyId:       companyId,
		Socket:          ws,
		RequestChannel:  make(chan *SoftwareRequest),
		ResponseChannel: make(chan *SoftwareResponse),
		CloseChannel: make(chan bool),
		Able:            true,
		AbleMutex:       new(sync.RWMutex),
	}
	if i.clients[c.CompanyId] == nil {
		i.clients[c.CompanyId] = make(map[int64]*SoftwareClient)
	}
	i.clients[c.CompanyId][c.Id] = c
	return c
}

// 删除软件客户端
func (i *SoftwareClientManager) DelClient(c *SoftwareClient) {
	i.rwm.Lock()
	defer i.rwm.Unlock()
	delete(i.clients[c.CompanyId], c.Id)
}

// 获取指定公司的软件客户端
func (i *SoftwareClientManager) GetCompanyClients(companyId uint32) []*SoftwareClient {
	i.rwm.RLock()
	defer i.rwm.Unlock()
	m := i.clients[companyId]
	result := make([]*SoftwareClient, 0, len(m))
	for k := range m {
		result = append(result, m[k])
	}
	return result
}
