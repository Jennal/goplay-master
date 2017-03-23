package master

import (
	"sync"

	"sort"

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

func (self *Services) Add(sess *session.Session, pack ServicePack) (pkg.Status, *pkg.ErrorMessage) {
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

func (self *Services) Update(sess *session.Session, pack ServicePack) (pkg.Status, *pkg.ErrorMessage) {
	self.mutex.Lock()
	self.serviceInfos[sess.ID] = pack
	self.mutex.Unlock()

	return pkg.STAT_OK, nil
}

func (self *Services) GetListByName(sess *session.Session, name string) ([]ServicePack, *pkg.ErrorMessage) {
	result := make([]ServicePack, 0)
	self.mutex.Lock()
	for _, sp := range self.serviceInfos {
		if sp.Name == name {
			result = append(result, sp)
		}
	}
	self.mutex.Unlock()

	return result, nil
}

func (self *Services) GetListByTags(sess *session.Session, tags []string) ([]ServicePack, *pkg.ErrorMessage) {
	result := make([]ServicePack, 0)
	self.mutex.Lock()
	for _, sp := range self.serviceInfos {
		if sp.Contains(tags...) {
			result = append(result, sp)
		}
	}
	self.mutex.Unlock()

	return result, nil
}

func (self *Services) GetByName(sess *session.Session, name string) (ServicePack, *pkg.ErrorMessage) {
	result, err := self.GetListByName(sess, name)
	if err != nil || len(result) <= 0 {
		return ServicePack{}, err
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ClientCount < result[j].ClientCount
	})

	return result[0], nil
}

func (self *Services) GetByTags(sess *session.Session, tags []string) (ServicePack, *pkg.ErrorMessage) {
	result, err := self.GetListByTags(sess, tags)
	if err != nil || len(result) <= 0 {
		return ServicePack{}, err
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ClientCount < result[j].ClientCount
	})

	return result[0], nil
}
