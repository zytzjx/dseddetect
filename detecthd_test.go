package main

import "testing"

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
