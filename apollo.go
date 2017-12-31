package main

import (
	"encoding/xml"
	"xcfg"
)

type Apollo struct {
	XMLName      xml.Name      `xml:"MyApollo"`
	ServiceGroups      []ServiceGroup `xml:"services"`
	MajorVersion int           `xml:"majorVersion,attr"`
	MinorVersion int           `xml:"minorVersion,attr"`
}
type ServiceGroup struct {
	Name string `xml:"name,attr"`
	Online bool `xml:"online,attr"`
	Services []ServiceItem `xml:"service"`
}
type ServiceItem struct {
	Address string `xml:"address,attr"`
	Online bool `xml:"online,attr"`
}

func NewApollo() *Apollo {
	apollo := &Apollo{}
	xcfg.LoadCfg(apollo)
	return apollo
}