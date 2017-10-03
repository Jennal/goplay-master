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

import "github.com/jennal/goplay/pkg"

const (
	/* Events */
	ON_BACKEND_UPDATED = "ON_BACKEND_UPDATED"
	ON_NEW_RPC_CLIENT  = "ON_NEW_RPC_CLIENT"

	/* Push */
	ON_CONNECTOR_GOT_NET_CLIENT      = "ON_CONNECTOR_GOT_NET_CLIENT"
	ON_CONNECTOR_CLIENT_DISCONNECTED = "ON_CONNECTOR_CLIENT_DISCONNECTED"
)

func init() {
	pkg.DefaultHandShake().RegistSpecRoute(ON_CONNECTOR_GOT_NET_CLIENT, 0xFFFF)
	pkg.DefaultHandShake().RegistSpecRoute(ON_CONNECTOR_CLIENT_DISCONNECTED, 0xFFFE)
}
