package main

import (
	"log"
	"net"
	"sync"
	"time"
)

var lightsAddress string
var lightsMutex sync.Mutex

func sendLightFunc(run runFunc) error {
	conn, err := net.Dial("udp", lightsAddress)
	if err != nil {
		log.Printf("light-conn-error: %v\n", err)
		return err
	}
	if err := run(conn); err != nil {
		return err
	}
	conn.Close()
	return nil
}

func sendLightCmd(conn net.Conn, cmd byte) error {
	_, err := conn.Write([]byte{cmd, 00, 85})
	if err != nil {
		log.Printf("light-send-error: %v\n", err)
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return nil
}

func setLightsToFullOn() error {
	lightsMutex.Lock()
	defer lightsMutex.Unlock()
	return sendLightFunc(func(conn net.Conn) error {
		if err := sendLightCmd(conn, 53); err != nil { // all on
			return err
		}
		if err := sendLightCmd(conn, 181); err != nil { // full brightness
			return err
		}
		return nil
	})
}

func setLightsToFullOff() error {
	lightsMutex.Lock()
	defer lightsMutex.Unlock()
	return sendLightFunc(func(conn net.Conn) error {
		err := sendLightCmd(conn, 57) // all off
		return err
	})
}

func setLightsToDiningMode() error {
	lightsMutex.Lock()
	defer lightsMutex.Unlock()
	return sendLightFunc(func(conn net.Conn) error {

		if err := sendLightCmd(conn, 53); err != nil { // all on
			return err
		}

		// NOTE: there are 10 light-levels, however, we want them at the dimmest state and this is UDP, so not always reliable
		for i := 0; i < 14; i++ {
			if err := sendLightCmd(conn, 52); err != nil { //dim
				return err
			}
		}

		if err := sendLightCmd(conn, 59); err != nil { // zone 1 off
			return err
		}
		if err := sendLightCmd(conn, 54); err != nil { // zone 4 off
			return err
		}

		return nil
	})
}
