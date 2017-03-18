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
