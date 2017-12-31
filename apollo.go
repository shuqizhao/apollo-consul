package main

import (
	"encoding/xml"
	"xcfg"
)

const CfgName  ="MyApollo"
type Apollo struct {
	XMLName      xml.Name      `xml:"MyApollo"`
	ConsulUrl      string `xml:"ConsulUrl,attr"`
	BuildPath string `xml:"BuildPath,attr"`
	AfterBuild string `xml:"AfterBuild,attr"`
	FixPage string `xml:"FixPage,attr"`
	ServiceGroups      []ServiceGroup `xml:"Services"`
	MajorVersion int           `xml:"majorVersion,attr"`
	MinorVersion int           `xml:"minorVersion,attr"`
}
type ServiceGroup struct {
	Name string `xml:"Name,attr"`
	Online bool `xml:"Online,attr"`
	Services []ServiceItem `xml:"Service"`
}
type ServiceItem struct {
	Id string `xml:"Id,attr"`
	Address string `xml:"Address,attr"`
	Url string `xml:"Url,attr"`
	Port string `xml:"Port,attr"`
	Tag string `xml:"Tag,attr"`
	Online bool `xml:"Online,attr"`
}

func NewApollo() Apollo {
	apollo := Apollo{}
	xcfg.LoadCfg(&apollo)
	return apollo
}