// Copyright (C) 2017 Jennal(jennalcn@gmail.com). All rights reserved.
//
// Licensed under the MIT License (the "License"); you may not use this file except
// in compliance with the License. You may obtain a copy of the License at
//
// http://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package master

import (
	"sync"

	"sort"

	"strings"

	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

// var serviceIDGen = session.IDGen

type Services struct {
	mutex          sync.Mutex
	serviceInfos   map[uint32]ServicePack
	sessionManager *session.SessionManager
}

func NewServices() *Services {
	return &Services{
		serviceInfos:   make(map[uint32]ServicePack),
		sessionManager: session.NewSessionManager(),
	}
}

func (self *Services) OnStarted() {
}

func (self *Services) OnStopped() {
}

func (self *Services) OnNewClient(sess *session.Session) {
	// id := serviceIDGen.NextID()
	// sess.Bind(id)
	log.Logf("%p => %d", sess, sess.ID)
	self.sessionManager.Add(sess)
}

func (self *Services) fixIP(sess *session.Session, pack *ServicePack) {
	if pack.IP == "" {
		pack.IP = sess.RemoteAddr().String()
		arr := strings.Split(pack.IP, ":")
		if arr != nil && len(arr) > 1 {
			pack.IP = arr[0]
		}
	}
}

func (self *Services) onServicePackUpdated(sp ServicePack) {
	if !sp.Type.IsBackend() {
		return
	}

	sessions := self.sessionManager.Sessions()
	for _, sess := range sessions {
		sessSp, ok := self.serviceInfos[sess.ID]
		if !ok {
			/* some service down */
			continue
		}

		if sessSp.IP == sp.IP && sessSp.Port == sp.Port {
			continue
		}

		if sessSp.Type.Is(ST_CONNECTOR) ||
			(sessSp.Type.IsBackend() && sessSp.Name == sp.Name) {
			log.Logf("Push backend Update: %v => %#v", sess.RemoteAddr(), sp)
			sess.Push(ON_BACKEND_UPDATED, sp)
		}
	}
}

func (self *Services) Add(sess *session.Session, pack ServicePack) (ServicePack, *pkg.ErrorMessage) {
	sess.Once(transfer.EVENT_CLIENT_DISCONNECTED, self, func(cli transfer.IClient) {
		self.mutex.Lock()
		defer self.mutex.Unlock()
		delete(self.serviceInfos, sess.ID)
	})

	self.fixIP(sess, &pack)

	self.mutex.Lock()
	log.Logf("%p => %d | %v", sess, sess.ID, pack)
	self.serviceInfos[sess.ID] = pack
	self.mutex.Unlock()

	self.onServicePackUpdated(pack)

	return pack, nil
}

func (self *Services) Update(sess *session.Session, pack ServicePack) (ServicePack, *pkg.ErrorMessage) {
	self.fixIP(sess, &pack)

	self.mutex.Lock()
	self.serviceInfos[sess.ID] = pack
	self.mutex.Unlock()

	self.onServicePackUpdated(pack)

	return pack, nil
}

func (self *Services) GetListByName(sess *session.Session, name string) ([]ServicePack, *pkg.ErrorMessage) {
	result := make([]ServicePack, 0)
	self.mutex.Lock()
	for _, sp := range self.serviceInfos {
		log.Log("====> ", sp.Addr(), " | ", sess.RemoteAddr().String())
		if sp.Addr() == sess.RemoteAddr().String() {
			continue
		}

		if sp.Name == name {
			result = append(result, sp)
		}
	}
	self.mutex.Unlock()

	if len(result) == 0 {
		return nil, pkg.NewErrorMessage(pkg.Status_ERR_EMPTY_RESULT, "no results")
	}

	return result, nil
}

func (self *Services) GetListByTags(sess *session.Session, tags []string) ([]ServicePack, *pkg.ErrorMessage) {
	result := make([]ServicePack, 0)
	self.mutex.Lock()
	for _, sp := range self.serviceInfos {
		log.Log("====> ", sp.Addr(), " | ", sess.RemoteAddr().String())
		if sp.Addr() == sess.RemoteAddr().String() {
			continue
		}

		if sp.Contains(tags...) {
			result = append(result, sp)
		}
	}
	self.mutex.Unlock()

	if len(result) == 0 {
		return nil, pkg.NewErrorMessage(pkg.Status_ERR_EMPTY_RESULT, "no results")
	}

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

func (self *Services) GetListByType(sess *session.Session, t ServiceType) ([]ServicePack, *pkg.ErrorMessage) {
	result := []ServicePack{}
	self.mutex.Lock()
	for _, sp := range self.serviceInfos {
		if sp.Type&t == t {
			result = append(result, sp)
		}
	}
	self.mutex.Unlock()

	if len(result) == 0 {
		return nil, pkg.NewErrorMessage(pkg.Status_ERR_EMPTY_RESULT, "no results")
	}

	return result, nil
}

func (self *Services) GetUniqueListByType(sess *session.Session, t ServiceType) ([]ServicePack, *pkg.ErrorMessage) {
	dataMap := make(map[string][]ServicePack)
	self.mutex.Lock()
	for _, sp := range self.serviceInfos {
		if sp.Type&t == t {
			if _, ok := dataMap[sp.Name]; !ok {
				dataMap[sp.Name] = []ServicePack{sp}
			} else {
				dataMap[sp.Name] = append(dataMap[sp.Name], sp)
			}
		}
	}
	self.mutex.Unlock()

	if len(dataMap) == 0 {
		return nil, pkg.NewErrorMessage(pkg.Status_ERR_EMPTY_RESULT, "no results")
	}

	result := []ServicePack{}

	for _, list := range dataMap {
		sort.Slice(list, func(i, j int) bool {
			return list[i].ClientCount < list[j].ClientCount
		})
		log.Log("list: ", list)
		result = append(result, list[0])
	}

	return result, nil
}
