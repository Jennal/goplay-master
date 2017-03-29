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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceType(t *testing.T) {
	st := ST_FRONTEND
	assert.Equal(t, true, st.IsFrontend())
	assert.Equal(t, false, st.IsBackend())
	st = ST_CONNECTOR
	assert.Equal(t, true, st.IsFrontend())
	assert.Equal(t, false, st.IsBackend())

	st = ST_BACKEND
	assert.Equal(t, false, st.IsFrontend())
	assert.Equal(t, true, st.IsBackend())
	st = ST_MASTER
	assert.Equal(t, false, st.IsFrontend())
	assert.Equal(t, true, st.IsBackend())

	// sp := ServicePack{
	// 	TagContainerImpl: data.TagContainerImpl{
	// 		Tags: map[string]bool{
	// 			"Hello": true,
	// 		},
	// 	},
	// 	Type:        ST_MASTER,
	// 	Name:        "",
	// 	IP:          "",
	// 	Port:        10,
	// 	ClientCount: 0,
	// }

	// encoder := encode.GetEncodeDecoder(pkg.ENCODING_JSON)
	// json, err := encoder.Marshal(sp)
	// t.Log(string(json), err)
}
