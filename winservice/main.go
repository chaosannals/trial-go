package main

import (
	"log"
	"os"

	"github.com/kardianos/service"
)

type ServiceProgram struct{}

//Start 启动
func (p *ServiceProgram) Start(s service.Service) error {
	go p.run()
	return nil
}

//Stop 停止
func (p *ServiceProgram) Stop(s service.Service) error {
	return nil
}

//run 执行
func (p *ServiceProgram) run() {

}

func main() {
	serviceConfig := &service.Config{
		Name:        "trial-go-winservice",
		DisplayName: "Trial Go WinService",
		Description: "Yet a Go Windows Service",
	}
	program := &ServiceProgram{}
	s, err := service.New(program, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			s.Install()
			log.Println("install successful.")
			return
		}
		if os.Args[1] == "uninstall" {
			s.Uninstall()
			log.Println("uninstall successful.")
			return
		}
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
