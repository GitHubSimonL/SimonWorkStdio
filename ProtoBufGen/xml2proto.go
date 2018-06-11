package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var cc Cap

type Cap struct {
	Cmdset Commands `xml:"Commands"`
	Bodys  []Body   `xml:"body"`
}

type Commands struct {
	Name       string    `xml:"name,attr"`
	Prefix     string    `xml:"prefix,attr"`
	StartValue int       `xml:"startvalue,attr"`
	Cmds       []Command `xml:"command"`
}
type Command struct {
	Name  string `xml:"name,attr"`
	Value int    `xml:"value,attr"`
	Desc  string `xml:"desc,attr"`
}

type Body struct {
	Name    string `xml:"name,attr"`
	Prefix  string `xml:"prefix,attr"`
	Postfix string `xml:"postfix,attr"`
	Items   []Item `xml:"item"`
}
type Item struct {
	Name    string `xml:"name,attr"`
	PreType string `xml:"pretype,attr"`
	Type    string `xml:"type,attr"`
	Desc    string `xml:"desc,attr"`
}

func parseCmd() string {
	com := fmt.Sprintf("syntax = \"proto3\";\nimport \"protocol/common/common.proto\";\npackage protocol;\nenum %s {\n\t_Blank = 0;\n", cc.Cmdset.Name)

	for k, v := range cc.Cmdset.Cmds {

		val := cc.Cmdset.StartValue + k
		comline := fmt.Sprintf("\t%s%s = %d;\t\t//%s\n", cc.Cmdset.Prefix, v.Name, val, v.Desc)
		comline = strings.ToUpper(comline)
		com += comline
	}
	com += "}\n"
	return com
}
func parseBody(param string) string {
	body := ""
	for _, v := range cc.Bodys {
		msg := fmt.Sprintf("message %s {\n", v.Name+v.Postfix)
		for k, item := range v.Items {
			dataType := item.Type
			if strings.EqualFold(strings.ToLower(item.Type), "boolean") {
				dataType = "bool"
			} else if strings.EqualFold(strings.ToLower(item.Type), "integer") {
				dataType = "int32"
			} else if strings.EqualFold(strings.ToLower(item.Type), "short") {
				dataType = "int32"
			} else if strings.EqualFold(strings.ToLower(item.Type), "ushort") {
				dataType = "int32"
			} else if strings.EqualFold(strings.ToLower(item.Type), "byte") {
				if strings.EqualFold(strings.ToLower(item.PreType), "repeated") {
					item.PreType = ""
					dataType = "bytes"
				}
			}
			PreType := item.PreType
			// if PreType != "repeated" {
			// 	PreType = "required"
			// }
			prefix := v.Prefix
			if param == "lua" {
				prefix = ""
			}
			vi := fmt.Sprintf("\t%s %s %s = %d;\n", PreType, dataType, prefix+item.Name, k+1)
			msg += vi
		}
		msg += "}\n"
		body += msg
	}
	return body
}
func loadProto(fileName string) {

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(content, &cc)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(cc)
}

func main() {
	type commands struct {
		cmds []Command
	}
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Println("usage: exename lua proto.xml")
		return
	}
	fileName := flag.Arg(1)
	fmt.Println(flag.Arg(1))
	loadProto(fileName)

	s := strings.Split(fileName, ".")
	fmt.Println(s[0])
	outFileName := s[0] + ".proto"
	of, err := os.Create(outFileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	com := parseCmd()
	h := []byte(com)
	of.Write(h)
	body := parseBody(flag.Arg(0))
	b := []byte(body)
	of.Write(b)
	defer of.Close()
}
