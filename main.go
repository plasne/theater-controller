package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// PROJECTOR=xxx.xxx.xxx.xxx:4998
// RECEIVER=xxx.xxx.xxx.xxx:23
// ROKU=xxx.xxx.xxx.xxx:8060
// LIGHTS=xxx.xxx.xxx.xxx:8899

var port int

type runFunc func(net.Conn) error

func first(items ...error) error {
	for _, item := range items {
		if item != nil {
			return item
		}
	}
	return nil
}

func listen(conn net.Conn, input chan string) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\r')
		if err != nil {
			// NOTE: we only listen for a period of time and then close the connection, so this is expected
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("listen-err: %v\n", err)
			}
			return
		}
		input <- str
	}
}

type VolumeResponse struct {
	Volume float64 `json:"volume"`
}

func init() {

	// configure
	_ = godotenv.Load()
	if potential, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64); err == nil {
		port = int(potential)
	} else {
		port = 8080
	}
	projectorAddress = os.Getenv("PROJECTOR")
	receiverAddress = os.Getenv("RECEIVER")
	rokuAddress = os.Getenv("ROKU")
	lightsAddress = os.Getenv("LIGHTS")
	lightsGroup = os.Getenv("LIGHTS_GROUP")
	lightsKey = os.Getenv("LIGHTS_KEY")
	if projectorAddress == "" || receiverAddress == "" || rokuAddress == "" || lightsAddress == "" || lightsGroup == "" || lightsKey == "" {
		log.Fatal("you must supply the following env vars: PROJECTOR, RECEIVER, ROKU, LIGHTS, LIGHTS_GROUP, LIGHTS_KEY.")
	}
	if potential, err := strconv.ParseInt(os.Getenv("PROJECTOR_IR_PORT"), 10, 64); err == nil {
		projectorIrPort = int(potential)
	} else {
		log.Fatal("you must supply a value for PROJECTOR_IR_PORT.")
	}

}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "./default.html", http.StatusMovedPermanently)
	})

	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// returning a 200 is good enough
	})

	r.HandleFunc("/api/lights/on", func(w http.ResponseWriter, r *http.Request) {
		if err := setLightsToFullOn(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/lights/off", func(w http.ResponseWriter, r *http.Request) {
		if err := setLightsToFullOff(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/lights/dining", func(w http.ResponseWriter, r *http.Request) {
		if err := setLightsToDiningMode(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/projector/on", func(w http.ResponseWriter, r *http.Request) {
		if err := turnProjectorOn(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/projector/off", func(w http.ResponseWriter, r *http.Request) {
		if err := turnProjectorOff(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/projector/toggle", func(w http.ResponseWriter, r *http.Request) {
		if err := toggleProjectorPower(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/receiver/on", func(w http.ResponseWriter, r *http.Request) {
		if err := turnReceiverOn(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/receiver/off", func(w http.ResponseWriter, r *http.Request) {
		if err := turnReceiverOff(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/receiver/volume", func(w http.ResponseWriter, r *http.Request) {
		if volume, err := getReceiverVolume(); err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(&VolumeResponse{Volume: volume})
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		}
	})

	r.HandleFunc("/api/receiver/up", func(w http.ResponseWriter, r *http.Request) {
		if volume, err := increaseReceiverVolume(); err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(&VolumeResponse{Volume: volume})
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		}
	})

	r.HandleFunc("/api/receiver/down", func(w http.ResponseWriter, r *http.Request) {
		if volume, err := decreaseReceiverVolume(); err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(&VolumeResponse{Volume: volume})
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		}
	})

	r.HandleFunc("/api/receiver/mute", func(w http.ResponseWriter, r *http.Request) {
		if err := muteReceiver(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/receiver/input/{input}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if err := setReceiverInput(vars["input"]); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/roku/press/{key}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if err := pressKeyOnRoku(vars["key"]); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/roku/launch/{appId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if err := launchOnRoku(vars["appId"]); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/roku/text", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "/api/roku/text only supports POST", 405)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		err = sendTextToRoku(string(body))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/roku/clear", func(w http.ResponseWriter, r *http.Request) {
		if err := clearTextOnRoku(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/on", func(w http.ResponseWriter, r *http.Request) {
		if err := turnTheaterOn(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	r.HandleFunc("/api/off", func(w http.ResponseWriter, r *http.Request) {
		if err := turnTheaterOff(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	log.Printf("sharing file system: %s", http.Dir("./www"))
	fs := http.FileServer(http.Dir("./www"))
	r.PathPrefix("/").Handler(fs)

	http.Handle("/", r)
	log.Printf("listening on port %v....", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))

}
