package main

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func TestLoadConfigXML(t *testing.T) {
	LoadConfigXML()
}

func TestLabel(t *testing.T) {
	reg, _ := regexp.Compile("label(\\d+)")
	s := reg.FindStringSubmatch("label5")
	if len(s) == 2 {
		d, _ := strconv.Atoi(s[1])
		fmt.Println(d)
	}
}
