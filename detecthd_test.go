package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestWriteHDDInfo2DB(t *testing.T) {
	CreateRedisPool(48)
	detectHDD := NewDataDetect()
	detectHDD.Locpath = "[AAAAAA]"
	detectHDD.Type = "disk"
	detectHDD.Manufacture = "WD"
	detectHDD.Model = "adb"
	detectHDD.Version = "mmmm"
	detectHDD.LinuxName = "/dev/sdc"
	detectHDD.Size = "8 TB"
	detectHDD.SGLibName = "/dev/sg5"
	//use smartctl get info
	//public String model replace using this
	detectHDD.Serialno = "sSerialno"
	detectHDD.LuwwndevID = "luwwndevId"
	detectHDD.Calibration = "sCalibration"
	detectHDD.UILabel = "label2"
	detectHDD.Label = 2
	//public String firware version replace by sVersion;
	detectHDD.Otherinfo["HD-health"] = "OK"
	WriteHDDInfo2DB(detectHDD)
}

func TestHardDiskMap(t *testing.T) {
	file, err := os.Open("disklist.txt")

	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var hddinfo []string
	//var hddchanged bool
	for scanner.Scan() {
		ss := scanner.Text()
		fmt.Println(ss)
		hddinfo = append(hddinfo, ss)
		if !DetectData.MatchKey(ss) {
			//hddchanged = true
			fmt.Println("Print log")
		}
		if !DetectData.ContainsKey(ss) {
			//hddchanged = true
			//\s Matches any white-space character.
			r := regexp.MustCompile(`^([\s\S]{13})(disk[\s\S]{4})([\s\S]{9})([\s\S]{17})([\s\S]{6})([\s\S]{11})([\s\S]{11})([\s\S]+)$`)
			diskinfos := r.FindStringSubmatch(ss)
			if len(diskinfos) == 9 {
				var dddect = NewSyncDataDetect()
				dddect.detectHDD.Locpath = strings.Trim(diskinfos[1], " ")
				dddect.detectHDD.Type = strings.Trim(diskinfos[2], " ")
				dddect.detectHDD.Manufacture = strings.Trim(diskinfos[3], " ")
				dddect.detectHDD.Model = strings.Trim(diskinfos[4], " ")
				dddect.detectHDD.Version = strings.Trim(diskinfos[5], " ")
				dddect.detectHDD.LinuxName = strings.Trim(diskinfos[6], " ")
				dddect.detectHDD.SGLibName = strings.Trim(diskinfos[7], " ")
				dddect.detectHDD.Size = strings.Trim(diskinfos[8], " ")

				if !strings.Contains(dddect.detectHDD.LinuxName, `/dev/`) {
					continue
				}
				DetectData.AddValue(ss, dddect)
			}
		}
	}

	// add info to
	diskinfo, err := os.Open("disksinfo.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer diskinfo.Close()

	scannerInfo := bufio.NewScanner(diskinfo)
	scannerInfo.Split(bufio.ScanLines)
	for scannerInfo.Scan() {
		sLine := scannerInfo.Text()
		//sudo smartctl /dev/sg0 -i -H
		r1 := regexp.MustCompile(`^sudo smartctl (.*?) -i -H`)
		disnames := r1.FindStringSubmatch(sLine)
		if len(disnames) == 2 {

		} else {

		}

	}

}
