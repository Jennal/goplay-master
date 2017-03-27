package master

import (
	"sort"
	"testing"

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

	s, err := client.Add(&sp)
	assert.Equal(t, pkg.STAT_OK, s)
	assert.Nil(t, err)
	assert.Equal(t, map[uint32]ServicePack{
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

	newSp := sp
	newSp.Tags = map[string]bool{
		"Hello": true,
		"World": true,
	}
	s, err = client1.Add(&newSp)
	assert.Equal(t, pkg.STAT_OK, s)
	assert.Nil(t, err)

	assert.Equal(t, map[uint32]ServicePack{
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

	sp.ClientCount = 10
	s, err = client.Update(&sp)
	assert.Equal(t, pkg.STAT_OK, s)
	assert.Nil(t, err)

	assert.Equal(t, map[uint32]ServicePack{
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

	list, err := client.GetListByName("test")
	assert.Nil(t, list)
	assert.Equal(t, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results"), err)

	list, err = client.GetListByName(NAME)
	assert.Nil(t, err)

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

	list, err = client.GetListByTags([]string{
		"Hello",
	})
	assert.Nil(t, err)

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

	list, err = client.GetListByTags([]string{
		"Hello", "World",
	})
	assert.Nil(t, err)

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

	list, err = client.GetListByTags([]string{
		"Hello", "world",
	})
	assert.Nil(t, list)
	assert.Equal(t, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results"), err)

	result, err := client.GetByName(NAME)
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
	}, result)
	assert.Nil(t, err)

	result, err = client.GetByName("none")
	assert.Equal(t, ServicePack{}, result)
	assert.Equal(t, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results"), err)

	result, err = client.GetByTags([]string{
		"Hello",
	})
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
	}, result)
	assert.Nil(t, err)

	result, err = client.GetByTags([]string{
		"Hello", "World",
	})
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
	}, result)
	assert.Nil(t, err)

	result, err = client.GetByTags([]string{
		"Hello", "not exists tag",
	})
	assert.Equal(t, ServicePack{}, result)
	assert.Equal(t, pkg.NewErrorMessage(pkg.STAT_ERR_EMPTY_RESULT, "no results"), err)
}
