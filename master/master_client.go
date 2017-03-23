package master

import (
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/service"
	"github.com/jennal/goplay/transfer"
)

type MasterClient struct {
	*service.ServiceClient
}

func NewMasterClient(cli transfer.IClient) *MasterClient {
	return &MasterClient{
		ServiceClient: service.NewServiceClient(cli),
	}
}

func (self *MasterClient) Add(pack ServicePack, succCB func(pkg.Status), failCB func(*pkg.ErrorMessage)) error {
	return self.Request("master.services.add", pack, succCB, failCB)
}

func (self *MasterClient) Update(pack ServicePack, succCB func(pkg.Status), failCB func(*pkg.ErrorMessage)) error {
	return self.Request("master.services.update", pack, succCB, failCB)
}

func (self *MasterClient) GetListByName(name string, succCB func([]ServicePack), failCB func(*pkg.ErrorMessage)) error {
	return self.Request("master.services.update", name, succCB, failCB)
}

func (self *MasterClient) GetListByTags(tags []string, succCB func([]ServicePack), failCB func(*pkg.ErrorMessage)) error {
	return self.Request("master.services.getlistbytags", tags, succCB, failCB)
}

func (self *MasterClient) GetByName(name string, succCB func(ServicePack), failCB func(*pkg.ErrorMessage)) error {
	return self.Request("master.services.getbyname", name, succCB, failCB)
}

func (self *MasterClient) GetByTags(tags []string, succCB func(ServicePack), failCB func(*pkg.ErrorMessage)) error {
	return self.Request("master.services.getbytags", tags, succCB, failCB)
}
