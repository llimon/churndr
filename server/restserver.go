package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/llimon/churndr/common"
)

func SetRoutes(router *mux.Router) {
	router.HandleFunc("/", HelloServer)
	router.HandleFunc("/users", StatusServer)
	router.HandleFunc("/pod/log/container/{namespace}/{pod}/{container}/{restart}", GetPodContainerLogs)

}

func RESTServer() {

	var bindAddress string

	if common.Config.Development {
		bindAddress = "localhost"
	} else {
		bindAddress = ""
	}

	router := mux.NewRouter().StrictSlash(true)
	SetRoutes(router)

	http.ListenAndServe(bindAddress+":"+strconv.Itoa(common.Config.Port), router)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go Away!, %s!", r.URL.Path[1:])
}

func GetPodContainerLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace, okNamespace := vars["namespace"]
	pod, okPod := vars["pod"]
	container, okContainer := vars["container"]
	restartCount, okRestartCount := vars["restart"]

	if !okNamespace || !okPod || !okContainer || !okRestartCount {
		common.Sugar.Infow("Missing parameter to GetPodContainerLogs",
			"pod", pod,
			"namespace", namespace,
			"container", container,
			"restartCount", restartCount)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		http.NotFound(w, r)

		// RETURN HTTP ERROR
	} else {
		common.Sugar.Infow("Debug",
			"pod", pod,
			"namespace", namespace,
			"container", container,
			"restartCount", restartCount)

		log, ok := common.PodLogs[pod+"/"+container+"/"+restartCount]
		if !ok {
			http.NotFound(w, r)
		} else {
			json, err := json.MarshalIndent(log, "", " ")
			if err != nil {
				common.Sugar.Infof("Failed to Marshal log, Now what?")
			}
			w.Header().Set("Content-Type", "text/json; charset=utf-8") // normal header
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, string(json))
		}

	}

}

func StatusServer(w http.ResponseWriter, r *http.Request) {

	// Convert PodCache map into a array for external consumpsion.
	statusPods := common.Status{Name: "status"}
	statusPods.SetLink("self", "/users", "")
	var out = []common.PodDB{}

	for _, currPod := range common.PodCache {
		out = append(out, currPod)
	}
	statusPods.SetEmbedded("users", jsonhal.Embedded(out))
	json, err := json.MarshalIndent(statusPods, "", " ")
	if err != nil {
		common.Sugar.Infof("Failed to Marshal PodCache, Now what?")
	}
	w.Header().Set("Content-Type", "text/json; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, string(json))

}
