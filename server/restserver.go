package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/llimon/page-on-pod-restarts/common"
)

func RESTServer() {
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/status", StatusServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go Away!, %s!", r.URL.Path[1:])
}

func StatusServer(w http.ResponseWriter, r *http.Request) {

	// Convert PodCache map into a array for external consumtion.
	var out = []common.PodDB{}

	for _, currPod := range common.PodCache {
		out = append(out, currPod)
	}
	json, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		common.Sugar.Infof("Failed to Marshal PodCache, Now what?")
	} else {
		fmt.Println(string(json))
	}
	w.Header().Set("Content-Type", "text/json; charset=utf-8") // normal header
	io.WriteString(w, string(json))

}
