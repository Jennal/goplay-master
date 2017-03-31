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
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/service"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

type BackendFilter struct {
	service        *service.Service
	sessionManager *session.SessionManager
}

func NewBackendFilter(serv *service.Service) filter.IFilter {
	return &BackendFilter{
		service:        serv,
		sessionManager: session.NewSessionManager(),
	}
}

func (self *BackendFilter) OnNewClient(sess *session.Session) bool /* return false to ignore */ {
	/* Do Nothing */
	return true
}

func (self *BackendFilter) OnRecv(sess *session.Session, header *pkg.Header, body []byte) bool /* return false to ignore */ {
	if !(header.Type == pkg.PKG_RPC_NOTIFY &&
		(header.Route == ON_CONNECTOR_GOT_NET_CLIENT ||
			header.Route == ON_CONNECTOR_CLIENT_DISCONNECTED)) {
		return true
	}

	s := self.sessionManager.GetSessionByID(sess.ID, sess.ClientID)
	if s == nil {
		s = session.NewSession(sess.IClient)
		s.Bind(sess.ID)
		s.BindClientID(sess.ClientID)
		s.SetEncoding(sess.Encoding)

		self.sessionManager.Add(s)
	}

	switch header.Route {
	case ON_CONNECTOR_GOT_NET_CLIENT:
		self.service.HandlerOnNewClient(s)
	case ON_CONNECTOR_CLIENT_DISCONNECTED:
		log.Log(">>>>>>>> BEFORE")
		s.Emit(transfer.EVENT_CLIENT_DISCONNECTED, s.IClient)
		log.Log("<<<<<<<< AFTER")
	}

	return false
}
