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

type IBackendServiceDelegate interface {
	transfer.IServerDelegate

	OnNewRpcClient(*service.ServiceClient)
}

type BackendService struct {
	*service.Service
	mc *MasterClient
}

func NewBackendService(name string, serv transfer.IServer) *BackendService {
	result := &BackendService{
		Service: service.NewService(name, serv),
		mc:      NewMasterClient(tcp.NewClient()),
	}

	result.UnregistDelegate(result.Service)
	result.RegistDelegate(result)

	result.RegistFilter(NewBackendFilter(result.Service))
	return result
}

func (serv *BackendService) RegistBackendDelegate(delegate IBackendServiceDelegate) {
	serv.Service.RegistDelegate(delegate)
	serv.On(ON_NEW_RPC_CLIENT, delegate, delegate.OnNewRpcClient)
}

func (serv *BackendService) UnregistBackendDelegate(delegate IBackendServiceDelegate) {
	serv.Service.UnregistDelegate(delegate)
	serv.Off(ON_NEW_RPC_CLIENT, delegate)
}

func (self *BackendService) ConnectMaster(host string, port int) error {
	sp := NewServicePack(ST_BACKEND, self.Name, self.Port())
	return self.mc.Bind(self, &sp, host, port)
}

func (self *BackendService) OnNewClient(client transfer.IClient) {
	// log.Log("**********************************")
	serviceClient := self.RegistNewClient(client)
	// self.HandlerOnNewClient(serviceClient.Session)
	serviceClient.Emit(transfer.EVENT_CLIENT_CONNECTED, client)
	self.Emit(ON_NEW_RPC_CLIENT, serviceClient)
}
