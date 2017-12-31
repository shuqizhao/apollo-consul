package main

import (
	"fmt"
	"os"
	"log"
	"strconv"
	consulapi "github.com/hashicorp/consul/api"
	"xcfg"
	"text/template"
)

func Register(){

	apolloEntity:=NewApollo()

	os.Setenv("CONSUL_HTTP_ADDR", apolloEntity.ConsulUrl)

	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
		return
	}
	for _,serviceGroup:=range apolloEntity.ServiceGroups{
		for _,service:=range serviceGroup.Services {
			registration := new(consulapi.AgentServiceRegistration)
			registration.ID = serviceGroup.Name+service.Id
			registration.Name = serviceGroup.Name
			a, _ := strconv.Atoi(service.Port)
			registration.Port = a
			registration.Tags = []string{service.Id}
			registration.Address = service.Url

			//增加check。
			check := new(consulapi.AgentServiceCheck)
			check.TCP =service.Url+":"+service.Port
			//设置超时 5s。
			check.Timeout = "5s"
			//设置间隔 5s。
			check.Interval = "5s"
			//注册check服务。
			registration.Check = check
			log.Println("get check.TCP:", check)

			err = client.Agent().ServiceRegister(registration)

			if err != nil {
				log.Println("register server error : ", err)
			}
		}
	}
}

func Check()  *Apollo{
	client, err := consulapi.NewClient(consulapi.DefaultConfig())

	if err != nil {
		fmt.Println(err)
		return nil
	}

	apolloEntity:=NewApollo()
	newApolloEntity := &Apollo{}
	newApolloEntity.AfterBuild=apolloEntity.AfterBuild
	newApolloEntity.ConsulUrl=apolloEntity.ConsulUrl
	newApolloEntity.BuildPath=apolloEntity.BuildPath
	newApolloEntity.FixPage=apolloEntity.FixPage
	newApolloEntity.ServiceGroups=[]ServiceGroup{}
	for _,serviceGroup:=range apolloEntity.ServiceGroups {
		if serviceGroup.Online{
			serviceGroupT:=ServiceGroup{}
			serviceGroupT.Name=serviceGroup.Name
			serviceGroupT.Online=serviceGroup.Online
			serviceGroupT.Services=[]ServiceItem{}
			services,_, err := client.Health().Service(serviceGroup.Name,"",true,nil)
			for _, v := range services {
				if v.Service.Service == serviceGroup.Name && IsOnline(v.Service.Tags[0],serviceGroup.Name,apolloEntity){
					serviceGroupT.Services = append(serviceGroupT.Services,ServiceItem{Id: v.Service.Tags[0],Url: v.Service.Address,Port: strconv.Itoa(v.Service.Port),Online:true})
				}
			}
			newApolloEntity.ServiceGroups = append(newApolloEntity.ServiceGroups, serviceGroupT)
			count := len(serviceGroupT.Services)
			if count == 0 {

			}

			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return newApolloEntity
}

func Build(apolloEntity *Apollo) {
	if apolloEntity == nil{
		return
	}
	cfgFolder := xcfg.GetAppCfgFolder()
	template_cfg_path := cfgFolder + "/" + CfgName + ".template"
	b := xcfg.ReadFile(template_cfg_path)

	t := template.New("MyApollo")
	t, _ = t.Parse(string(b))
	build_cfg_path := cfgFolder + "/" + CfgName + ".build"
	file, err := os.Create(build_cfg_path)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(file, apolloEntity)
	file.Close()
}

func IsOnline(id string,name string,apollo *Apollo) bool  {
	for _,v := range apollo.ServiceGroups{
		if v.Name == name{
			for _,v1:=range v.Services{
				if v1.Id==id{
					return v1.Online
				}
			}
		}
	}
	return false
}
