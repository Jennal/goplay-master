package master

import (
	"github.com/jennal/goplay/aop"
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

func (self *MasterClient) Add(pack *ServicePack) (stat pkg.Status, err *pkg.ErrorMessage) {
	aop.Parallel(func(c chan bool) {
		e := self.Request("master.services.add", pack, func(st pkg.Status) {
			stat = st
			err = nil
			c <- true
		}, func(e *pkg.ErrorMessage) {
			stat = pkg.STAT_ERR
			err = e
			c <- true
		})

		if e != nil {
			stat = pkg.STAT_ERR
			err = pkg.NewErrorMessage(pkg.STAT_ERR, e.Error())
		}
	})

	return
}

func (self *MasterClient) Update(pack *ServicePack, succCB func(pkg.Status), failCB func(*pkg.ErrorMessage)) error {
	return self.Request("master.services.update", pack, succCB, failCB)
}

func (self *MasterClient) GetListByName(name string, succCB func([]ServicePack), failCB func(*pkg.ErrorMessage)) error {
	return self.Request("master.services.getlistbyname", name, succCB, failCB)
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
