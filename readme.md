### Hard Drive Detect

#### Database 0:
  This is common data. for settings.
  |key|default|detail|
  |:---|-----:|:---------|
  |labelcnt|0|max support Harddisk, it basic on cards|

  Transaction:
  |key|default|detail|
  |:---|-----:|:---------|
  |site||from cmc config|
  |company||from cmc config|
  |productid||from cmc config|
  |solutionid||from cmc config|
  |operator||from login|
  |workstationName||get pc name "hostname"|
  |windowsversion||get Linux version "uname -srio"|   
       
       


#### Databse 1~n:
  These are label1~n. for Detect, Run, Trascation etc.  
  Detect:  
   if disconnect == status then remove the table field except status field
  |key|default|detail|
  |:---|-----:|:---------|
  |status|connected/disconnected/running|port device status|
  |size||hd size|
  |make||manufacture|
  |model||model|
  |serialnumber|||
  |linuxname||linux device name, /dev/sda|
  |sglibName||sg device name, /dev/sg2|
  |badsectors|||
  |HD-health||PASSED,OK is GOOD, other unhealthy|
  |trasaction||it is HSET, all field will upload to cmc server|  
  ||||

 
  Runing:    
   if disconnect == status then remove all fields
  |key|default|detail|
  |:----|-----:|:---------|
  |speed|0 MB/s|write hd speed|
  |optime|00:01/00:01|write data time and ...|
  |progress|0.00%|write data progress|

  




