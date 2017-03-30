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
	"github.com/jennal/goplay/filter"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/service"
	"github.com/jennal/goplay/session"
)

type BackendFilter struct {
	service *service.Service
}

func NewBackendFilter(serv *service.Service) filter.IFilter {
	return &BackendFilter{
		service: serv,
	}
}

func (self *BackendFilter) OnNewClient(sess *session.Session) bool /* return false to ignore */ {
	/* Do Nothing */
	return true
}

func (self *BackendFilter) OnRecv(sess *session.Session, header *pkg.Header, body []byte) bool /* return false to ignore */ {
	if header.Type != pkg.PKG_RPC_NOTIFY || header.Route != ON_CONNECTOR_GOT_NET_CLIENT {
		return true
	}

	self.service.HandlerOnNewClient(sess)
	return false
}
