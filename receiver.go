package main

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var receiverAddress string
var receiverMutex sync.Mutex

func sendReceiverFunc(run runFunc) error {
	conn, err := net.Dial("tcp", receiverAddress)
	if err != nil {
		log.Printf("receiver-conn-error: %v\n", err)
		return err
	}
	defer conn.Close()
	if err := run(conn); err != nil {
		return err
	}
	return nil
}

func sendToReceiver(conn net.Conn, cmd string) error {
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		log.Printf("receiver-send-error: %v", err)
		return err
	}
	return nil
}

func askReceiver(conn net.Conn, cmd []byte) (*string, error) {
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		log.Printf("receiver-send-error: %v", err)
		return nil, err
	}

	messages := make(chan string, 1)
	go listen(conn, messages)
	select {
	case message := <-messages:
		return &message, nil
	case <-time.After(250 * time.Millisecond):
		return nil, errors.New("timeout without response")
	}
}

func turnReceiverOn() error {
	receiverMutex.Lock()
	defer receiverMutex.Unlock()
	if err := sendReceiverFunc(func(conn net.Conn) error {
		err := sendToReceiver(conn, "PWON\r")
		return err
	}); err != nil {
		return err
	}
	if err := sendReceiverFunc(func(conn net.Conn) error {
		err := sendToReceiver(conn, "SISAT/CBL\r")
		return err
	}); err != nil {
		return err
	}
	return nil
}

func turnReceiverOff() error {
	receiverMutex.Lock()
	defer receiverMutex.Unlock()
	return sendReceiverFunc(func(conn net.Conn) error {
		err := sendToReceiver(conn, "PWSTANDBY\r")
		return err
	})
}

func parseVolume(val *string) float64 {
	if val == nil {
		return 0.0
	}
	str := *val
	if !strings.HasPrefix(str, "MV") {
		return 0.0
	}
	short := str[2 : len(str)-1]
	v64, err := strconv.ParseFloat(short, 64)
	if err != nil {
		return 0.0
	}
	if len(short) > 2 {
		return v64 / 10
	}
	return v64
}

func runReceiverVolumeCommand(cmd string) (float64, error) {
	var volume float64
	err := sendReceiverFunc(func(conn net.Conn) error {
		val, err := askReceiver(conn, []byte(cmd))
		if err == nil {
			volume = parseVolume(val)
		}
		return err
	})
	return volume, err
}

func getReceiverVolume() (float64, error) {
	receiverMutex.Lock()
	defer receiverMutex.Unlock()
	return runReceiverVolumeCommand("MV?\r")
}

func increaseReceiverVolume() (float64, error) {
	receiverMutex.Lock()
	defer receiverMutex.Unlock()
	return runReceiverVolumeCommand("MVUP\r")
}

func decreaseReceiverVolume() (float64, error) {
	receiverMutex.Lock()
	defer receiverMutex.Unlock()
	return runReceiverVolumeCommand("MVDOWN\r")
}

func muteReceiver() error {
	receiverMutex.Lock()
	defer receiverMutex.Unlock()
	return sendReceiverFunc(func(conn net.Conn) error {
		err := sendToReceiver(conn, "MUON\r")
		return err
	})
}

func setReceiverInput(input string) error {
	receiverMutex.Lock()
	defer receiverMutex.Unlock()
	switch strings.ToLower(input) {
	case "roku":
		return sendReceiverFunc(func(conn net.Conn) error {
			err := sendToReceiver(conn, "SISAT/CBL\r")
			return err
		})
	case "apple", "appletv":
		return sendReceiverFunc(func(conn net.Conn) error {
			err := sendToReceiver(conn, "SIBD\r")
			return err
		})
	case "dvd", "bluray", "disc":
		return sendReceiverFunc(func(conn net.Conn) error {
			err := sendToReceiver(conn, "SIDVD\r")
			return err
		})
	}
	return errors.New("specified input is not valid")
}
