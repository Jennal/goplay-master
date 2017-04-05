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

import "github.com/jennal/goplay/data"
import "fmt"

type ServiceType uint8

const (
	ST_FRONTEND  ServiceType = 0x01
	ST_BACKEND   ServiceType = 0x02
	ST_MASTER    ServiceType = 0x10 | ST_BACKEND
	ST_CONNECTOR ServiceType = 0x20 | ST_FRONTEND
)

func (st ServiceType) Is(t ServiceType) bool {
	return (st & t) == t
}

func (st ServiceType) IsFrontend() bool {
	return st.Is(ST_FRONTEND)
}

func (st ServiceType) IsBackend() bool {
	return st.Is(ST_BACKEND)
}

type ServicePack struct {
	data.TagContainerImpl

	Type        ServiceType
	Name        string
	IP          string
	Port        int
	ClientCount int
}

func NewServicePack(t ServiceType, name string, port int) ServicePack {
	return ServicePack{
		TagContainerImpl: data.TagContainerImpl{
			Tags: make(map[string]bool),
		},
		Type:        t,
		Name:        name,
		Port:        port,
		ClientCount: 0,
	}
}

func (sp ServicePack) Addr() string {
	return fmt.Sprintf("%v:%v", sp.IP, sp.Port)
}
