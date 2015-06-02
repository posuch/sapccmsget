// main.go
package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"github.com/fedorpatlin/sapccms"
)

var flags struct {
	mtename *string
	host	*string
	sapsid	*string
	instNr	*string
//	mode    *string // avg1, avg5, avg15
}

func prepareUsage(){
	hostname,_ := os.Hostname()
	flags.sapsid = flag.String("s", "", "SAP system name")
	flags.host = flag.String("h", hostname, "SAP host name")
	flags.instNr = flag.String("n", "00", "Instance number")
//	flags.mode = flag.String("m", "avg1", "Value type of avg1, avg5, avg15" )
	flags.mtename = flag.String("e", "", "CCMS monitoring tree element full name")
}

func usage(){
	
}

func createMteRequest() *sapccms.MsgMtGetTidByNameRequest{
	req := sapccms.NewMsgMtGetTidByNameRequest()
	req.SoapRequest.Item=make([]sapccms.ALMTNAMEL,1)
	if flags.mtename == nil {
		flag.PrintDefaults()
		os.Exit(-1)
	}
	req.SoapRequest.Item[0].SetCompleteName(*flags.mtename)
	return req
}

func createMteResponse() *sapccms.MsgMtGetTidByNameResponse{
	res := sapccms.NewMsgMtGetTidByNameResponse()
	return res
}

func getMteByName(srv sapccms.SAPCCMS, name string)sapccms.ALGTIDLNRC{
		var mteRequest = createMteRequest()
	var mteResponse = createMteResponse()
	
	if err := srv.MtGetTidByName(mteRequest, mteResponse); err != nil {
		log.Fatal(err.Error())
	}
	return mteResponse.GetTidTable().GetItem()[0]
}

func preparePerfRequest(tid sapccms.ALGLOBTID) *sapccms.MsgPerfReadRequest{
	msg := sapccms.NewMsgPerfReadRequest()
	msg.SoapRequest.Item = make([]sapccms.ALGLOBTID,1)
	msg.SoapRequest.Item[0] = tid
	return msg
}

func getPerfByTid(srv sapccms.SAPCCMS, tid sapccms.ALGLOBTID) *sapccms.MsgPerfReadResponse{
	sendMsg := preparePerfRequest(tid)
	rcvMsg := sapccms.NewMsgPerfReadResponse()
	if err := srv.PerfRead(sendMsg, rcvMsg);err != nil{
		log.Fatal(err.Error())
	}
	return rcvMsg
}

func main() {
	prepareUsage()
	flag.Parse()
	srv := sapccms.NewSAPCCMS(fmt.Sprintf("http://%s:5%s13/SAPCCMS.cgi", *flags.host, *flags.instNr))
	tid := getMteByName(srv, *flags.mtename)
	perf := getPerfByTid(srv, tid.Tid)
	fmt.Println(perf.PerfReadResponse.TidTable.Item[0].PerfValue.Avg01PerfValue)
}
