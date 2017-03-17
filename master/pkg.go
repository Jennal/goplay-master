package master

import "github.com/jennal/goplay/data"

type ServiceType uint8

const (
	ST_FRONTEND  ServiceType = 0x01
	ST_BACKEND   ServiceType = 0x02
	ST_MASTER    ServiceType = 0x10 | ST_BACKEND
	ST_CONNECTOR ServiceType = 0x20 | ST_FRONTEND
)

func (st ServiceType) IsFrontend() bool {
	return (st & ST_FRONTEND) == ST_FRONTEND
}

func (st ServiceType) IsBackend() bool {
	return (st & ST_BACKEND) == ST_BACKEND
}

type ServicePack struct {
	data.TagContainerImpl

	Type        ServiceType
	Name        string
	IP          string
	Port        int
	ClientCount int
}
