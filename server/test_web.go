package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Authcode struct {
	ReceiverGuid string `json:"receiverguid"`
}

type Join struct {
	SenderGuid string `json:"senderguid"`
}

type ReceiverStandby struct {
	Authcode              string `json:"authcode"`
	Connectguid           string `json:"connectguid"`
	MptSessionDescription struct {
		Sdp  string `json:"sdp"`
		Type string `json:"type"`
	} `json:"mptSessionDescription"`
	Receiverguid string `json:"receiverguid"`
}

type Joined struct {
	RetCode     string `json:"retcode"`
	Message     string `json:"message"`
	ConnectGuid string `json:"connectguid"`
	Signal      struct {
		Address string `json:"address"`
		IP      string `json:"ip"`
		Type    string `json:"type"`
		Port    int    `json:"port"`
	} `json:"signal"`
	MptSessionDescription struct {
		Sdp  string `json:"sdp"`
		Type string `json:"type"`
	} `json:"mptSessionDescription"`
}

var authcode Authcode
var standby ReceiverStandby

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/wasm")
}

func main() {
	http.HandleFunc("/receiver/pt_authcode", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &authcode)
		response := `{"retcode": "200","message": "","authcode": "234567","signal": {"address": "stpush.startsupport.com","ip": "14.63.228.216","type": "SSL","port": 443}}`

		fmt.Fprintf(w, response)
	})

	http.HandleFunc("/receiver/pt_standby", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &standby)
		fmt.Print(standby)
		fmt.Fprintf(w, `{"retcode":"200,"message":""}`)
	})

	http.HandleFunc("/sender/pt_join", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		body, _ := ioutil.ReadAll(r.Body)

		var res Join
		json.Unmarshal(body, &res)

		joined := Joined{}
		joined.RetCode = "200"
		joined.Message = ""
		joined.ConnectGuid = authcode.ReceiverGuid
		joined.Signal.Address = "stpush.startsupport.com"
		joined.Signal.IP = "14.63.228.216"
		joined.Signal.Type = "SSL"
		joined.Signal.Port = 443
		joined.MptSessionDescription.Type = standby.MptSessionDescription.Type
		joined.MptSessionDescription.Sdp = standby.MptSessionDescription.Sdp

		jsonBytes, _ := json.Marshal(joined)
		jsonString := string(jsonBytes)
		fmt.Fprintf(w, jsonString)
	})

	http.Handle("/", http.FileServer(http.Dir("./")))

	fmt.Println("http://127.0.0.1:8081")
	log.Fatal(http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil))
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
