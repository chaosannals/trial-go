package basis

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

// 网页前端客户端
type FrontendClient struct {
	Id              int64
	CompanyId       uint32
	RequestChannel  chan *FrontendRequest
	ResponseChannel chan *FrontendResponse
	CloseChannel    chan bool
}

// 网页前端客户端管理器
type FrontendClientManager struct {
	clients       map[uint32]map[int64]*FrontendClient
	snowflakeNode *snowflake.Node
	rwm           *sync.RWMutex
}

// 创建一个前端客户端管理器
func NewFrontendClientManager() (*FrontendClientManager, error) {
	snode, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	return &FrontendClientManager{
		clients:       make(map[uint32]map[int64]*FrontendClient),
		snowflakeNode: snode,
		rwm:           new(sync.RWMutex),
	}, nil
}

// 创建一个前端客户端
func (i *FrontendClientManager) NewClient(companyId uint32) *FrontendClient {
	i.rwm.Lock()
	defer i.rwm.Unlock()
	c := &FrontendClient{
		Id:              i.snowflakeNode.Generate().Int64(),
		CompanyId:       companyId,
		RequestChannel:  make(chan *FrontendRequest),
		ResponseChannel: make(chan *FrontendResponse),
		CloseChannel:    make(chan bool),
	}
	if i.clients[c.CompanyId] == nil {
		i.clients[c.CompanyId] = make(map[int64]*FrontendClient)
	}
	i.clients[c.CompanyId][c.Id] = c
	return c
}

// 删除一个前端客户端
func (i *FrontendClientManager) DelClient(c *FrontendClient) {
	i.rwm.Lock()
	defer i.rwm.Unlock()
	delete(i.clients[c.CompanyId], c.Id)
}

// 获取指定公司的前端客户端列表
func (i *FrontendClientManager) GetCompanyClients(companyId uint32) []*FrontendClient {
	i.rwm.RLock()
	defer i.rwm.RUnlock()
	m := i.clients[companyId]
	result := make([]*FrontendClient, len(m))
	for k := range m {
		result = append(result, m[k])
	}
	return result
}
