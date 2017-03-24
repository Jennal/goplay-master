package main

import (
	"github.com/jennal/goplay-master/master"
	"github.com/jennal/goplay/cmd"
	"github.com/jennal/goplay/service"
	"github.com/jennal/goplay/transfer/tcp"
)

func main() {
	ser := tcp.NewServer("", master.PORT)
	serv := service.NewService(master.NAME, ser)

	serv.RegistHanlder(master.NewServices())
	cmd.Start(serv)
}
