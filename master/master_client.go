package master

import (
	"sync"

	"github.com/jennal/goplay/aop"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/service"
	"github.com/jennal/goplay/transfer"
)

type MasterClient struct {
	*service.ServiceClient
	server transfer.IServer

	data        *ServicePack
	dataMutex   sync.Mutex
	isDataDirty bool
}

func NewMasterClient(cli transfer.IClient) *MasterClient {
	return &MasterClient{
		ServiceClient: service.NewServiceClient(cli),
		data:          nil,
		isDataDirty:   false,
	}
}

func (self *MasterClient) Bind(serv transfer.IServer, sp *ServicePack, host string, port int) error {
	self.server = serv
	self.data = sp

	loopFunc := func() {
		go func() {
			for {
				if serv.IsStarted() == false {
					return
				}

				if self.isDataDirty {
					self.isDataDirty = false
					_, err := self.Update(self.data)
					if err != nil {
						log.Error(err)
						return
					}
				}
			}
		}()
	}

	if serv.IsStarted() {
		loopFunc()
	} else {
		serv.On(transfer.EVENT_SERVER_STARTED, self, loopFunc)
	}

	serv.On(transfer.EVENT_SERVER_NEW_CLIENT, self, func(sess transfer.IClient) {
		sess.On(transfer.EVENT_CLIENT_DISCONNECTED, self, func(cli transfer.IClient) {
			self.dataMutex.Lock()
			defer self.dataMutex.Unlock()
			self.data.ClientCount--
			self.isDataDirty = true
		})

		self.dataMutex.Lock()
		defer self.dataMutex.Unlock()
		self.data.ClientCount++
		self.isDataDirty = true
	})

	err := self.Connect(host, port)
	if err != nil {
		return err
	}

	_, em := self.Add(self.data)
	if err != nil {
		return log.NewError(em.Error())
	}

	return nil
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
