package master

import (
	"sort"
	"testing"

	"time"

	"github.com/jennal/goplay/data"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/service"
	"github.com/jennal/goplay/transfer/tcp"
	"github.com/stretchr/testify/assert"
)

func newClient() *MasterClient {
	cli := tcp.NewClient()
	client := NewMasterClient(cli)

	client.Connect("", PORT)
	return client
}

func sortList(list []ServicePack) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].ClientCount > list[j].ClientCount
	})
}

func TestMasterClient(t *testing.T) {
	ser := tcp.NewServer("", PORT)
	serv := service.NewService(NAME, ser)

	hdl := NewServices()
	serv.RegistHanlder(hdl)
	serv.Start()

	client := newClient()
	client1 := newClient()

	sp := ServicePack{
		TagContainerImpl: data.TagContainerImpl{
			Tags: map[string]bool{
				"Hello": true,
			},
		},
		Type:        ST_MASTER,
		Name:        NAME,
		IP:          "",
		Port:        PORT,
		ClientCount: 0,
	}

	callIn := make(map[string]bool)

	client.Add(sp, func(s pkg.Status) {
		assert.Equal(t, pkg.STAT_OK, s)

		assert.Equal(t, map[int]ServicePack{
			0: ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 0,
			},
		}, hdl.serviceInfos)
		callIn["add"] = true
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "never come to here")
	})
	time.Sleep(100 * time.Millisecond)
	newSp := sp
	newSp.Tags = map[string]bool{
		"Hello": true,
		"World": true,
	}
	client1.Add(newSp, func(s pkg.Status) {
		assert.Equal(t, pkg.STAT_OK, s)

		assert.Equal(t, map[int]ServicePack{
			0: ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 0,
			},
			1: ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
						"World": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 0,
			},
		}, hdl.serviceInfos)
		callIn["add1"] = true
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "never come to here")
	})
	time.Sleep(100 * time.Millisecond)
	sp.ClientCount = 10
	client.Update(sp, func(s pkg.Status) {
		assert.Equal(t, pkg.STAT_OK, s)

		assert.Equal(t, map[int]ServicePack{
			0: ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 10,
			},
			1: ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
						"World": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 0,
			},
		}, hdl.serviceInfos)
		callIn["update"] = true
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "never come to here")
	})
	time.Sleep(100 * time.Millisecond)
	client.GetListByName("test", func(list []ServicePack) {
		assert.True(t, false, "never come to here")
	}, func(err *pkg.ErrorMessage) {
		assert.Equal(t, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results"), err)
		callIn["GetListByName"] = true
	})

	client.GetListByName(NAME, func(list []ServicePack) {
		sortList(list)

		assert.Equal(t, []ServicePack{
			ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 10,
			},
			ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
						"World": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 0,
			},
		}, list)
		callIn["GetListByName1"] = true
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "never come to here")
	})

	client.GetListByTags([]string{
		"Hello",
	}, func(list []ServicePack) {
		sortList(list)

		assert.Equal(t, []ServicePack{
			ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 10,
			},
			ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
						"World": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 0,
			},
		}, list)
		callIn["GetListByTags"] = true
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "never come to here")
	})

	client.GetListByTags([]string{
		"Hello", "World",
	}, func(list []ServicePack) {
		sortList(list)

		assert.Equal(t, []ServicePack{
			ServicePack{
				TagContainerImpl: data.TagContainerImpl{
					Tags: map[string]bool{
						"Hello": true,
						"World": true,
					},
				},
				Type:        ST_MASTER,
				Name:        NAME,
				IP:          "127.0.0.1",
				Port:        PORT,
				ClientCount: 0,
			},
		}, list)
		callIn["GetListByTags1"] = true
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "never come to here")
	})

	client.GetListByTags([]string{
		"Hello", "world",
	}, func(list []ServicePack) {
		assert.True(t, false, "never come to here")
	}, func(err *pkg.ErrorMessage) {
		assert.Equal(t, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results"), err)
		callIn["GetListByTags2"] = true
	})

	client.GetByName(NAME, func(sp ServicePack) {
		assert.Equal(t, ServicePack{
			TagContainerImpl: data.TagContainerImpl{
				Tags: map[string]bool{
					"Hello": true,
					"World": true,
				},
			},
			Type:        ST_MASTER,
			Name:        NAME,
			IP:          "127.0.0.1",
			Port:        PORT,
			ClientCount: 0,
		}, sp)
		callIn["GetByName"] = true
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "never come to here")
	})

	client.GetByName("none", func(sp ServicePack) {
		assert.True(t, false, "never come to here")
	}, func(err *pkg.ErrorMessage) {
		assert.Equal(t, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results"), err)
		callIn["GetByName1"] = true
	})

	client.GetByTags([]string{
		"Hello",
	}, func(sp ServicePack) {
		assert.Equal(t, ServicePack{
			TagContainerImpl: data.TagContainerImpl{
				Tags: map[string]bool{
					"Hello": true,
					"World": true,
				},
			},
			Type:        ST_MASTER,
			Name:        NAME,
			IP:          "127.0.0.1",
			Port:        PORT,
			ClientCount: 0,
		}, sp)
		callIn["GetByTags"] = true
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "never come to here")
	})

	client.GetByTags([]string{
		"Hello", "World",
	}, func(sp ServicePack) {
		assert.Equal(t, ServicePack{
			TagContainerImpl: data.TagContainerImpl{
				Tags: map[string]bool{
					"Hello": true,
					"World": true,
				},
			},
			Type:        ST_MASTER,
			Name:        NAME,
			IP:          "127.0.0.1",
			Port:        PORT,
			ClientCount: 0,
		}, sp)
		callIn["GetByTags1"] = true
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "never come to here")
	})

	client.GetByTags([]string{
		"Hello", "not exists tag",
	}, func(sp ServicePack) {
		assert.True(t, false, "never come to here")
	}, func(err *pkg.ErrorMessage) {
		assert.Equal(t, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results"), err)
		callIn["GetByTags2"] = true
	})

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, map[string]bool{
		"add":            true,
		"add1":           true,
		"update":         true,
		"GetListByName":  true,
		"GetListByName1": true,
		"GetListByTags":  true,
		"GetListByTags1": true,
		"GetListByTags2": true,
		"GetByName":      true,
		"GetByName1":     true,
		"GetByTags":      true,
		"GetByTags1":     true,
		"GetByTags2":     true,
	}, callIn)
}
