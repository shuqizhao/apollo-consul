package main

import (
	"fmt"
	"os"
	"log"
	"strconv"
	consulapi "github.com/hashicorp/consul/api"
	"xcfg"
	"text/template"
	"io"
	"os/exec"
	"strings"
)

func Register(apolloEntity *Apollo){
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
			//log.Println("get check.TCP:", check)

			err = client.Agent().ServiceRegister(registration)

			if err != nil {
				log.Println("register server error : ", err)
			}
		}
	}
}

func Check(apolloEntity *Apollo)  *Apollo{
	client, err := consulapi.NewClient(consulapi.DefaultConfig())

	if err != nil {
		fmt.Println(err)
		return nil
	}

	newApolloEntity := &Apollo{}
	newApolloEntity.AfterBuild=apolloEntity.AfterBuild
	newApolloEntity.ConsulUrl=apolloEntity.ConsulUrl
	newApolloEntity.BuildPath=apolloEntity.BuildPath
	newApolloEntity.FixPage=apolloEntity.FixPage
	newApolloEntity.ServiceGroups=[]ServiceGroup{}
	for _,serviceGroup:=range apolloEntity.ServiceGroups {
		serviceGroupT:=ServiceGroup{}
		serviceGroupT.Name=serviceGroup.Name
		serviceGroupT.Online=serviceGroup.Online
		serviceGroupT.Services=[]ServiceItem{}
		if serviceGroup.Online{
			services,_, err := client.Health().Service(serviceGroup.Name,"",true,nil)
			for _, v := range services {
				if v.Service.Service == serviceGroup.Name && IsOnline(v.Service.Tags[0],serviceGroup.Name,apolloEntity){
					serviceGroupT.Services = append(serviceGroupT.Services,ServiceItem{Id: v.Service.Tags[0],Url: v.Service.Address,Port: strconv.Itoa(v.Service.Port),Online:true})
				}
			}
			newApolloEntity.ServiceGroups = append(newApolloEntity.ServiceGroups, serviceGroupT)
			count := len(serviceGroupT.Services)
			if count == 0 {
				//serviceGroupT.Services = append(serviceGroupT.Services,ServiceItem{Id: "FixPage",Url: apolloEntity.FixPage,Online:true})
			}
			if err != nil {
				fmt.Println(err)
			}
		}else{
			//serviceGroupT.Services = append(serviceGroupT.Services,ServiceItem{Id: "FixPage",Url: apolloEntity.FixPage,Online:true})
		}
	}
	return newApolloEntity
}

func Build(apolloEntity *Apollo) {
	defer func(){
		if err:=recover();err!=nil{
			fmt.Println(err)
		}
	}()
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

	tmpFile := apolloEntity.BuildPath + "." + xcfg.GetGuid()
	CopyFile(tmpFile,build_cfg_path)
	if xcfg.Exist(apolloEntity.BuildPath){
		os.Remove(apolloEntity.BuildPath)
	}
	os.Rename(tmpFile, apolloEntity.BuildPath)
	os.Remove(tmpFile)

    cmds := strings.Split(apolloEntity.AfterBuild," ")
    if len(cmds)>1{
		f, err := exec.Command(cmds[0], cmds[1:]...).Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(f))
	}else{
		f, err := exec.Command(cmds[0], "").Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(f))
	}
	
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

func IsChange(apolloEntity *Apollo,newApolloEntity *Apollo) bool {
	if apolloEntity == nil{
		return true
	}
	for _,v := range apolloEntity.ServiceGroups{
		serGroup:=GetServiceGroup(v.Name,newApolloEntity)
		if serGroup==nil{
			return true
		}else if len(serGroup.Services)!=len(v.Services){
			return true
		}
	}
	return false
}

func GetServiceGroup(name string,apollo *Apollo) *ServiceGroup{
	for _,v := range apollo.ServiceGroups{
		if v.Name == name{
			return &v
		}
	}
	return nil
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
