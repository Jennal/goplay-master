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

	"github.com/jennal/goplay/channel"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
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

	serviceInfos map[string]ServicePack
	services     map[string]*service.ServiceClient

	servicesMutex sync.Mutex
}

func NewBackendService(name string, serv transfer.IServer) *BackendService {
	result := &BackendService{
		Service:      service.NewService(name, serv),
		mc:           NewMasterClient(tcp.NewClient()),
		serviceInfos: make(map[string]ServicePack),
		services:     make(map[string]*service.ServiceClient),
	}

	result.mc.On(ON_BACKEND_UPDATED, result, func(sp ServicePack) {
		if result.mc.IsSelfServicePack(&sp) {
			return
		}

		if item, ok := result.serviceInfos[sp.Addr()]; ok {
			if item.IP == sp.IP && item.Port == sp.Port {
				result.serviceInfos[sp.Addr()] = sp
			}

			return
		}

		result.connectBackend(sp)
	})

	result.UnregistDelegate(result.Service)
	result.RegistDelegate(result)

	result.RegistFilter(NewBackendFilter(result.Service))
	return result
}

func (self *BackendService) connectBackend(sp ServicePack) {
	if self.mc.IsSelfServicePack(&sp) {
		return
	}

	self.serviceInfos[sp.Addr()] = sp

	cli := tcp.NewClient()
	client := service.NewServiceClient(cli)

	client.Once(transfer.EVENT_CLIENT_CONNECTED, self, func(cli transfer.IClient) {
		self.servicesMutex.Lock()
		defer self.servicesMutex.Unlock()

		self.services[sp.Addr()] = client
	})
	client.Once(transfer.EVENT_CLIENT_DISCONNECTED, self, func(cli transfer.IClient) {
		self.servicesMutex.Lock()
		defer self.servicesMutex.Unlock()

		delete(self.services, sp.Addr())
		delete(self.serviceInfos, sp.Addr())
	})

	err := client.Connect(sp.IP, sp.Port)
	if err != nil {
		log.Error(err)
	}
}

func (self *BackendService) GetAllBackends() map[string]*service.ServiceClient {
	self.servicesMutex.Lock()
	defer self.servicesMutex.Unlock()
	backends := make(map[string]*service.ServiceClient)
	for name, item := range self.services {
		backends[name] = item
	}

	return backends
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
	err := self.mc.Bind(self, &sp, host, port)
	if err != nil {
		return err
	}

	list, pe := self.mc.GetListByName(self.Name)
	if pe != nil {
		return pe
	}

	for _, sp := range list {
		// log.Log("===> ", sp)
		self.connectBackend(sp)
	}

	return nil
}

func (self *BackendService) OnNewClient(client transfer.IClient) {
	// log.Log("**********************************")
	serviceClient := self.RegistNewClient(client)
	serviceClient.On(transfer.EVENT_CLIENT_SENT, self, func(cli transfer.IClient, header *pkg.Header, data []byte) {
		// if header.Type == pkg.PKG_HEARTBEAT || header.Type == pkg.PKG_HEARTBEAT_RESPONSE {
		// 	return
		// }
		// log.Logf("<==== Recv from Other BE:\n\theader => %#v\n\tbody => %#v | %v\n", header, data, string(data))

		if !channel.IsPush(header) {
			return
		}

		backends := self.GetAllBackends()
		for _, be := range backends {
			// log.Log("====> Sending to: ", be.RemoteAddr())
			be.Broadcast(header.Route, data)
		}
	})

	serviceClient.On(transfer.EVENT_CLIENT_RECVED, self, func(cli transfer.IClient, header *pkg.Header, data []byte) {
		// if header.Type == pkg.PKG_HEARTBEAT || header.Type == pkg.PKG_HEARTBEAT_RESPONSE {
		// 	return
		// }
		// log.Logf("<==== Recv from Other BE:\n\theader => %#v\n\tbody => %#v | %v\n", header, data, string(data))

		if !channel.IsBroadcast(header) {
			return
		}

		name := channel.GetChannelName(header.Route)
		ch := channel.GetChannelManager().Get(name)
		// log.Log("\t=>>>>> ", name, "\t", ch)
		if ch != nil {
			ch.BroadcastRaw(name, data)
		}
	})

	// self.HandlerOnNewClient(serviceClient.Session)
	serviceClient.Emit(transfer.EVENT_CLIENT_CONNECTED, client)
	self.Emit(ON_NEW_RPC_CLIENT, serviceClient)
}
