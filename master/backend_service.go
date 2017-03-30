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
