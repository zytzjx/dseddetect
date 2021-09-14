package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	Log "github.com/zytzjx/anthenacmc/loggersys"
)

// DataDetect Data Detect
type DataDetect struct {
	Locpath     string `json:"locpath"`
	Type        string `json:"sType"`
	Manufacture string `json:"sManufacture"`
	Model       string `json:"sModel"`
	Version     string `json:"sVersion"`
	LinuxName   string `json:"sLinuxName"`
	Size        string `json:"sSize"`
	SGLibName   string `json:"sSGLibName"`
	//use smartctl get info
	//public String model replace using this
	Serialno    string `json:"sSerialno"`
	LuwwndevID  string `json:"luwwndevId"`
	Calibration string `json:"sCalibration"`
	UILabel     string `json:"sUILabel"`
	Label       int    `json:"-"`
	//public String firware version replace by sVersion;
	Otherinfo map[string]string `json:"otherinfo"`
}

// NewDataDetect create Detect Data
func NewDataDetect() *DataDetect {
	return &DataDetect{Otherinfo: make(map[string]string)}
}

// AddMap2Other  add map
func (dd *DataDetect) AddMap2Other(ddother map[string]string) {
	if len(ddother) == 0 {
		return
	}

	for kk, vv := range ddother {
		dd.Otherinfo[kk] = vv
	}
}

// SyncDataDetect sync data
type SyncDataDetect struct {
	lock      *sync.RWMutex
	IsRunning bool
	detectHDD *DataDetect
}

// NewSyncDataDetect create data
func NewSyncDataDetect() *SyncDataDetect {
	return &SyncDataDetect{lock: new(sync.RWMutex), IsRunning: false, detectHDD: NewDataDetect()}
}

// SetRunning running
func (sdd *SyncDataDetect) SetRunning() {
	sdd.IsRunning = true
}

// CleanRunning Clean Running
func (sdd *SyncDataDetect) CleanRunning() {
	sdd.IsRunning = false
}

func (sdd *SyncDataDetect) String() string {
	sdd.lock.Lock()
	defer sdd.lock.Unlock()

	jsonString, err := json.Marshal(sdd.detectHDD)
	if err != nil {
		return "{}"
	}
	return string(jsonString)
}

// SyncMap struct
type SyncMap struct {
	lock     *sync.RWMutex
	dddetect map[string]*SyncDataDetect
}

// NewSyncMap  new SyncMap struct
func NewSyncMap() *SyncMap {
	return &SyncMap{lock: new(sync.RWMutex), dddetect: make(map[string]*SyncDataDetect)}
}

// MatchKey find key
func (sm *SyncMap) MatchKey(key string) bool {
	var okret bool
	if len(key) == 0 {
		return okret
	}
	re := regexp.MustCompile(`^([\s\S]{13})(disk[\s\S]{4})([\s\S]{9})([\s\S]{17})([\s\S]{6})([\s\S]{11})([\s\S]{11})([\s\S]+)$`)
	if !re.MatchString(key) {
		return true
	}
	var sik = re.ReplaceAllString(key, `$1$2$3$4$5$7$8`)
	sm.lock.Lock()
	defer sm.lock.Unlock()

	for kk := range sm.dddetect {
		sikk := re.ReplaceAllString(kk, `$1$2$3$4$5$7$8`)
		if strings.Compare(sik, sikk) == 0 {
			okret = true
			//fmt.Print(key)
			break
		}
	}

	return okret
}

// ContainsKey check key
func (sm *SyncMap) ContainsKey(key string) bool {
	if len(key) == 0 {
		return false
	}
	sm.lock.Lock()
	defer sm.lock.Unlock()
	_, ok := sm.dddetect[key]
	return ok
}

// AddValue add value
func (sm *SyncMap) AddValue(key string, dd *SyncDataDetect) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if dd == nil {
		delete(sm.dddetect, key)
	} else {
		sm.dddetect[key] = dd
	}

}

// RemoveOld remove
func (sm *SyncMap) RemoveOld(newkeylist []string) int {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	icount := 0 //remove count
	for kk := range sm.dddetect {
		if !stringInSlice(kk, newkeylist) {
			//sm.dddetect[kk].detectHDD.Label
			//TODO: need remove data from database
			ResetDB(sm.dddetect[kk].detectHDD.Label)
			Set(sm.dddetect[kk].detectHDD.Label, "status", "disconnected", 0)
			Publish(sm.dddetect[kk].detectHDD.Label, "status", "disconnected")
			delete(sm.dddetect, kk)
			icount++
		}
	}
	return icount
}

// Get get value
func (sm *SyncMap) Get(key string) (*SyncDataDetect, bool) {
	if len(key) == 0 {
		return nil, false
	}
	sm.lock.Lock()
	defer sm.lock.Unlock()
	vv, ok := sm.dddetect[key]
	return vv, ok
}

// Add add
func (sm *SyncMap) Add(key string) *SyncDataDetect {
	if len(key) == 0 {
		return nil
	}
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if vv, ok := sm.dddetect[key]; ok {
		return vv
	}
	sbc := NewSyncDataDetect()
	sm.dddetect[key] = sbc
	return sbc
}

// Remove remove
func (sm *SyncMap) Remove(key string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	delete(sm.dddetect, key)
}

// ItemToString get single Detect by Label
func (sm *SyncMap) ItemToString(id string) string {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	var smddinfo *DataDetect
	for _, v := range sm.dddetect {
		if v.detectHDD.UILabel == fmt.Sprintf("label%s", id) {
			smddinfo = v.detectHDD
			break
		}
	}

	if smddinfo == nil {
		return "{}"
	}

	jsonString, err := json.Marshal(smddinfo)
	if err != nil {
		return "{}"
	}
	return string(jsonString)
}

// GetLabels get label connect
func (sm *SyncMap) GetLabels() []int {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	labels := make([]int, 0)
	for _, v := range sm.dddetect {
		if v.detectHDD.Label > 0 {
			labels = append(labels, v.detectHDD.Label)
		}
	}
	return labels
}

// IndexString get label connect
func (sm *SyncMap) IndexString() string {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	labels := make([]int, 0)
	for _, v := range sm.dddetect {
		if v.detectHDD.Label > 0 {
			labels = append(labels, v.detectHDD.Label)
		}
	}
	jsonString, _ := json.Marshal(labels)
	return string(jsonString)
}

func (sm *SyncMap) String() string {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	smddinfo := make(map[string]*DataDetect)
	for k, v := range sm.dddetect {
		smddinfo[k] = v.detectHDD
	}

	jsonString, err := json.Marshal(smddinfo)
	if err != nil {
		return "{}"
	}
	return string(jsonString)
}

// DetectData for detect
var DetectData = NewSyncMap()

// SASHDDinfo sas hd info
var SASHDDinfo = NewSyncSASHDDMap()

//var detectIndexes []int

// WriteHDDInfo2DB write detect information to labelDB
func WriteHDDInfo2DB(detectHDD *DataDetect) {
	// type DataDetect struct {
	// 	Locpath     string `json:"locpath"`
	// 	Type        string `json:"sType"`
	// 	Manufacture string `json:"sManufacture"`
	// 	Model       string `json:"sModel"`
	// 	Version     string `json:"sVersion"`
	// 	LinuxName   string `json:"sLinuxName"`
	// 	Size        string `json:"sSize"`
	// 	SGLibName   string `json:"sSGLibName"`
	// 	Serialno    string `json:"sSerialno"`
	// 	LuwwndevID  string `json:"luwwndevId"`
	// 	Calibration string `json:"sCalibration"`
	// 	UILabel     string `json:"sUILabel"`
	// 	Label       int    `json:"-"`
	// 	//public String firware version replace by sVersion;
	// 	Otherinfo map[string]string `json:"otherinfo"`
	// }
	Log.Log.Info("WriteHDDInfo2DB++")
	if detectHDD.Label == 0 {
		return
	}
	//ResetDB(detectHDD.Label)
	Set(detectHDD.Label, "status", "connected", 0)
	Publish(detectHDD.Label, "status", "connected")
	Set(detectHDD.Label, "size", detectHDD.Size, 0)
	Set(detectHDD.Label, "make", detectHDD.Manufacture, 0)
	Set(detectHDD.Label, "model", detectHDD.Model, 0)
	Set(detectHDD.Label, "serialnumber", detectHDD.Serialno, 0)
	Set(detectHDD.Label, "linuxname", detectHDD.LinuxName, 0)
	Set(detectHDD.Label, "sglibName", detectHDD.SGLibName, 0)

	SetTransaction(detectHDD.Label, "portNumber", detectHDD.Label, "sourceModel", detectHDD.Model, "sourceMake", detectHDD.Manufacture, "serialnumber", detectHDD.Serialno, "esnNumber", detectHDD.Serialno)

	// string sSize = deviceNode["size"] != null ? deviceNode["size"].InnerText : "";
	// if (sSize.IndexOf("BadSectors(")>0) sSize = sSize.Substring(0,sSize.IndexOf("BadSectors(")-1).Trim();
	// writer.WriteElementString("size", sSize);
	SetTransaction(detectHDD.Label, "size", detectHDD.Size)

	badsectors := detectHDD.Otherinfo["badsectors"]
	SetTransaction(detectHDD.Label, "badsectors", badsectors)
	Set(detectHDD.Label, "badsectors", badsectors, 0)
	jsonString, err := json.Marshal(detectHDD.Otherinfo)
	if err == nil {
		SetTransaction(detectHDD.Label, "otherinfo", jsonString)
	}
	healthy, ok := detectHDD.Otherinfo["HD-health"]
	if !ok {
		healthy = "Failed.FD"
	}
	Set(detectHDD.Label, "HD-health", healthy, 0)
	Set(detectHDD.Label, "calibration", detectHDD.Calibration, 0)
	Set(detectHDD.Label, "uilabel", detectHDD.UILabel, 0)
	SetTransaction(detectHDD.Label, "Healthy", healthy)

}

func contains(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func SetDisconnected() {
	labels := DetectData.GetLabels()
	for i := 1; i <= 48; i++ {
		if contains(labels, i) {
			continue
		}
		Set(i, "status", "disconnected", 0)
	}
}

// MergeCalibration merge
func MergeCalibration() {
	Log.Log.Info("MergeCalibration ++")
	SASHDDinfo.lock.Lock()
	defer SASHDDinfo.lock.Unlock()
	DetectData.lock.Lock()
	defer DetectData.lock.Unlock()

	//Log.Log.Info(DetectData.String())
	//detectIndexes = []int{}
	for index, sasmap := range SASHDDinfo.SASHDDMapData {
		if mm, ok := SASHDDinfo.ReadStatus[index]; ok {
			if !mm {
				continue
			}
		}
		SASHDDinfo.ReadStatus[index] = false
		for _, card := range sasmap {
			Serial, oks := card["Serial"]
			if Serial == "" {
				oks = false
			}
			GUID, okg := card["GUID"]
			if GUID == "" {
				okg = false
			}
			for k, v := range DetectData.dddetect {
				if len(v.detectHDD.Serialno) == 0 {
					continue
				}
				sserno := v.detectHDD.Serialno
				sserno = strings.Replace(sserno, "-", "", -1)
				var sLuID string
				var ok bool
				if sLuID, ok = v.detectHDD.Otherinfo["LogicalUnitID"]; !ok {
					sLuID = ""
				}
				if (oks && ((sserno != "" && Serial != "") && strings.HasPrefix(sserno, Serial) || strings.HasPrefix(Serial, sserno))) ||
					(okg && ((GUID != "" && v.detectHDD.LuwwndevID != "" && sLuID != "") && strings.EqualFold(GUID, v.detectHDD.LuwwndevID) || strings.EqualFold(GUID, sLuID))) {
					if slot, ok := card["Slot"]; ok {
						v.detectHDD.Calibration = fmt.Sprintf("%d_%s", index, slot)
					}
					if len(v.detectHDD.Calibration) > 0 {
						v.detectHDD.UILabel, ok = configxmldata.Conf.GetPortMap()[v.detectHDD.Calibration]
						if ok {
							reg, _ := regexp.Compile(`label(\d+)`)
							ss := reg.FindStringSubmatch(v.detectHDD.UILabel)
							if len(ss) == 2 {
								v.detectHDD.Label, _ = strconv.Atoi(ss[1])
							}
						}
					}

					for kkk, vvv := range card {
						v.detectHDD.Otherinfo[kkk] = vvv
					}
				}
				DetectData.dddetect[k] = v
				WriteHDDInfo2DB(v.detectHDD)
			}

		}
	}

}

func main() {
	Log.NewLogger("dseddetect")
	verinfo := "version:21.9.13.0; author:Jeffery zhang"
	Log.Log.Info(verinfo)
	fmt.Println(verinfo)
	fmt.Println("http://localhost:12000/print")
	fmt.Println("http://localhost:12000/labels")
	fmt.Println("http://localhost:12000/print/{[0-9]+}")
	nDelay := flag.Int("interval", 10, "interval run check disk.")
	flag.Parse()

	LoadConfigXML()
	CreateRedisPool(calibrationports)
	go func() {
		for {
			RunListDisk()

			MergeCalibration()
			//fmt.Println("interval:", *nDelay)
			SetDisconnected()

			time.Sleep(time.Duration(*nDelay) * time.Second)
		}
	}()

	time.Sleep(20 * time.Second)

	r := mux.NewRouter()
	// Add your routes as needed
	r.HandleFunc("/print", PrintHandler).Methods("GET")
	r.HandleFunc("/labels", LabelsHandler).Methods("GET")
	r.HandleFunc("/print/{id:[0-9]+}", LabelHandler).Methods("GET")

	srv := &http.Server{
		Addr: "127.0.0.1:12000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			Log.Log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	Log.Log.Info("shutting down")
	os.Exit(0)
}
