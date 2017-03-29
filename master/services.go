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

var serviceIDGen = session.IDGen

type Services struct {
	mutex        sync.Mutex
	serviceInfos map[uint32]ServicePack
}

func NewServices() *Services {
	return &Services{
		serviceInfos: make(map[uint32]ServicePack),
	}
}

func (self *Services) OnStarted() {
}

func (self *Services) OnStopped() {
}

func (self *Services) OnNewClient(sess *session.Session) {
	id := serviceIDGen.NextID()
	sess.Bind(id)
	log.Logf("%p => %d", sess, sess.ID)
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

func (self *Services) Add(sess *session.Session, pack ServicePack) (pkg.Status, *pkg.ErrorMessage) {
	sess.On(transfer.EVENT_CLIENT_DISCONNECTED, self, func(cli transfer.IClient) {
		self.mutex.Lock()
		defer self.mutex.Unlock()
		delete(self.serviceInfos, sess.ID)
	})

	self.fixIP(sess, &pack)

	self.mutex.Lock()
	log.Logf("%p => %d", sess, sess.ID)
	self.serviceInfos[sess.ID] = pack
	self.mutex.Unlock()

	return pkg.STAT_OK, nil
}

func (self *Services) Update(sess *session.Session, pack ServicePack) (pkg.Status, *pkg.ErrorMessage) {
	self.fixIP(sess, &pack)

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

	if len(result) == 0 {
		return nil, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results")
	}

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

	if len(result) == 0 {
		return nil, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results")
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
