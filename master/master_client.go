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
			c <- true
		}
	})

	return
}

func (self *MasterClient) Update(pack *ServicePack) (stat pkg.Status, err *pkg.ErrorMessage) {
	aop.Parallel(func(c chan bool) {
		e := self.Request("master.services.update", pack, func(s pkg.Status) {
			stat = s
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
			c <- true
		}
	})

	return
}

func (self *MasterClient) GetListByName(name string) (result []ServicePack, err *pkg.ErrorMessage) {
	aop.Parallel(func(c chan bool) {
		e := self.Request("master.services.getlistbyname", name, func(s []ServicePack) {
			result = s
			err = nil
			c <- true
		}, func(e *pkg.ErrorMessage) {
			result = nil
			err = e
			c <- true
		})

		if e != nil {
			result = nil
			err = pkg.NewErrorMessage(pkg.STAT_ERR, e.Error())
			c <- true
		}
	})

	return
}

func (self *MasterClient) GetListByTags(tags []string) (result []ServicePack, err *pkg.ErrorMessage) {
	aop.Parallel(func(c chan bool) {
		e := self.Request("master.services.getlistbytags", tags, func(s []ServicePack) {
			result = s
			err = nil
			c <- true
		}, func(e *pkg.ErrorMessage) {
			result = nil
			err = e
			c <- true
		})

		if e != nil {
			result = nil
			err = pkg.NewErrorMessage(pkg.STAT_ERR, e.Error())
			c <- true
		}
	})

	return
}

func (self *MasterClient) GetByName(name string) (result ServicePack, err *pkg.ErrorMessage) {
	aop.Parallel(func(c chan bool) {
		e := self.Request("master.services.getbyname", name, func(s ServicePack) {
			result = s
			err = nil
			c <- true
		}, func(e *pkg.ErrorMessage) {
			err = e
			c <- true
		})

		if e != nil {
			err = pkg.NewErrorMessage(pkg.STAT_ERR, e.Error())
			c <- true
		}
	})

	return
}

func (self *MasterClient) GetByTags(tags []string) (result ServicePack, err *pkg.ErrorMessage) {
	aop.Parallel(func(c chan bool) {
		e := self.Request("master.services.getbytags", tags, func(s ServicePack) {
			result = s
			err = nil
			c <- true
		}, func(e *pkg.ErrorMessage) {
			err = e
			c <- true
		})

		if e != nil {
			err = pkg.NewErrorMessage(pkg.STAT_ERR, e.Error())
			c <- true
		}
	})

	return
}
