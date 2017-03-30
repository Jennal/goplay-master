package master

import (
	"github.com/jennal/goplay/service"
	"github.com/jennal/goplay/transfer"
	"github.com/jennal/goplay/transfer/tcp"
)

type BackendService struct {
	*service.Service
	mc *MasterClient
}

func NewBackendService(name string, serv transfer.IServer) *BackendService {
	return &BackendService{
		Service: service.NewService(name, serv),
		mc:      NewMasterClient(tcp.NewClient()),
	}
}

func (self *BackendService) ConnectMaster(host string, port int) error {
	sp := NewServicePack(ST_BACKEND, self.Name, self.Port())
	return self.mc.Bind(self, &sp, host, port)
}
