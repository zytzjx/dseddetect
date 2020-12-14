package main

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type item struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

type config struct {
	CardList string `xml:"value,attr"`
	PortMap  []item `xml:"add"`
}

type configuration struct {
	Conf config `xml:"mapconf"`
}

// Split rune
func Split(r rune) bool {
	return r == ',' || r == ';' || r == ':'
}

func (cc config) GetCardListIndex() ([]int, error) {
	var myslice []int
	if len(cc.CardList) > 0 {
		patts := strings.FieldsFunc(cc.CardList, Split)
		myslice := make([]int, len(patts))
		var i int = 0

		for _, patt := range patts {
			index, err := strconv.Atoi(patt)
			if err == nil {
				myslice[i] = index
				i++
			}
		}
		return myslice[:i], nil
	}
	return myslice, errors.New("card list is empty")
}

// GetPortMap get map
func (cc config) GetPortMap() map[string]string {
	labmap := make(map[string]string)
	for _, vv := range cc.PortMap {
		labmap[vv.Key] = vv.Value
	}
	return labmap
}

var configxmldata *configuration
var calibrationports int = 48

// LoadConfigXML load xml
func LoadConfigXML() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
		return
	}

	configxmlpath := path.Join(dir, "appconf.xml")

	xmlFile, err := os.Open(configxmlpath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		//xmlFile = []byte(xmlstr)
		return
	}
	defer xmlFile.Close()

	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	configxmldata = &configuration{}
	err = xml.Unmarshal(b, configxmldata)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configxmldata)
	calibrationports = len(configxmldata.Conf.PortMap)
	RunCardList()
	SetLabelCount(calibrationports)
}

// RunCardList   get card count
func RunCardList() {
	cmd := exec.Command("./sas2ircu", "LIST")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	timer := time.AfterFunc(600*time.Second, func() {
		cmd.Process.Kill()
	})

	cnt := 0
	scanner := bufio.NewScanner(stdout)
	var re = regexp.MustCompile(`(?m)(\d)\s+\S+\s+[0-9a-f]{4}h\s+[0-9a-f]{2}h\s+\S+\s+[0-9a-f]{4}h\s+[0-9a-f]{4}h`)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			ss := scanner.Text()
			fmt.Println(ss)
			if len(re.FindAllString(ss, -1)) > 0 {
				cnt++
			}
		}
		done <- true
	}()
	<-done
	if cnt > 0 {
		if calibrationports > cnt*16 {
			calibrationports = cnt * 16
		}
	}
	if cmd.Wait() != nil {
		fmt.Println("run sas2ircu failed")
	}
	timer.Stop()

}
