package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/shuqizhao/xcfg"
)

func Check(apolloEntity *Apollo) *Apollo {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	newApolloEntity := &Apollo{}
	newApolloEntity.AfterBuild = apolloEntity.AfterBuild
	newApolloEntity.ConsulUrl = apolloEntity.ConsulUrl
	newApolloEntity.BuildPath = apolloEntity.BuildPath
	newApolloEntity.FixPage = apolloEntity.FixPage
	newApolloEntity.ServiceGroups = []ServiceGroup{}
	for _, serviceGroup := range apolloEntity.ServiceGroups {
		serviceGroupT := ServiceGroup{}
		serviceGroupT.Name = serviceGroup.Name
		serviceGroupT.Online = serviceGroup.Online
		serviceGroupT.IsEnable = serviceGroup.IsEnable
		serviceGroupT.Services = []ServiceItem{}
		if serviceGroup.Online {
			for _, v := range serviceGroup.Services {
				serviceGroupT.Services = append(serviceGroupT.Services, ServiceItem{Id: v.Id, Address: v.Address, Url: v.Url, Tag: v.Tag, Port: v.Port, Online: v.Online})
			}
		}
		if serviceGroup.IsEnable {
			newApolloEntity.ServiceGroups = append(newApolloEntity.ServiceGroups, serviceGroupT)
		}
	}
	return newApolloEntity
}

func Build(apolloEntity *Apollo) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if apolloEntity == nil {
		return
	}
	cfgFolder := xcfg.GetAppCfgFolder()
	template_cfg_path := cfgFolder + "/" + CfgName + ".template"
	b := xcfg.ReadFile(template_cfg_path)

	t, err := template.New("MyApollo").Parse(string(b))
	if err != nil {
		fmt.Println(err)
		return
	}
	build_cfg_path := cfgFolder + "/" + CfgName + ".build"
	file, err := os.Create(build_cfg_path)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(file, apolloEntity)
	file.Close()

	tmpFile := apolloEntity.BuildPath + "." + xcfg.GetGuid()
	CopyFile(tmpFile, build_cfg_path)
	if xcfg.Exist(apolloEntity.BuildPath) {
		os.Remove(apolloEntity.BuildPath)
	}
	os.Rename(tmpFile, apolloEntity.BuildPath)
	os.Remove(tmpFile)

	cmds := strings.Split(apolloEntity.AfterBuild, " ")
	if len(cmds) > 1 {
		f, err := exec.Command(cmds[0], cmds[1:]...).Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(f))
	} else {
		f, err := exec.Command(cmds[0], "").Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(f))
	}

}

func IsChange(apolloEntity *Apollo, newApolloEntity *Apollo) bool {
	if apolloEntity == nil || newApolloEntity == nil {
		return true
	}
	for _, v := range apolloEntity.ServiceGroups {
		serGroup := GetServiceGroup(v.Name, newApolloEntity)
		if serGroup == nil {
			return true
		} else if len(serGroup.Services) != len(v.Services) {
			return true
		}
	}
	return false
}

func GetServiceGroup(name string, apollo *Apollo) *ServiceGroup {
	for _, v := range apollo.ServiceGroups {
		if v.Name == name {
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
