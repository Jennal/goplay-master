package master

import (
	"sync"

	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

var serviceIDGen = session.NewIDGen()

type Services struct {
	mutex        sync.Mutex
	serviceInfos map[int]ServicePack
}

func (self *Services) OnStarted() {
}

func (self *Services) OnStopped() {
}

func (self *Services) OnNewClient(sess *session.Session) {
	id := serviceIDGen.NextID()
	sess.Bind(id)
}

func (self *Services) AddService(sess *session.Session, pack ServicePack) (pkg.Status, *pkg.ErrorMessage) {
	sess.On(transfer.EVENT_CLIENT_DISCONNECTED, self, func() {
		self.mutex.Lock()
		defer self.mutex.Unlock()
		delete(self.serviceInfos, sess.ID)
	})

	self.mutex.Lock()
	self.serviceInfos[sess.ID] = pack
	self.mutex.Unlock()

	return pkg.STAT_OK, nil
}

func (self *Services) UpdateService(sess *session.Session, pack ServicePack) (pkg.Status, *pkg.ErrorMessage) {
	self.mutex.Lock()
	self.serviceInfos[sess.ID] = pack
	self.mutex.Unlock()

	return pkg.STAT_OK, nil
}

func (self *Services) GetByName(sess *session.Session, name string) (ServicePack, *pkg.ErrorMessage) {
	return ServicePack{}, nil
}

func (self *Services) GetByTags(sess *session.Session, tags []string) (ServicePack, *pkg.ErrorMessage) {
	return ServicePack{}, nil
}
