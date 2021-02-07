package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var projectorAddress string
var projectorIrPort int
var projectorMutex sync.Mutex

func sendProjectorFunc(run runFunc) error {
	conn, err := net.Dial("tcp", projectorAddress)
	if err != nil {
		log.Printf("projector-conn-error: %v\n", err)
		return err
	}
	if err := run(conn); err != nil {
		return err
	}
	conn.Close()
	return nil
}

func sendProjectorCmd(conn net.Conn, cmd string) error {

	_, err := conn.Write([]byte(cmd))
	if err != nil {
		log.Printf("projector-send-error: %v", err)
		return err
	}

	message, err := bufio.NewReader(conn).ReadString('\r')
	if err != nil {
		log.Printf("projector-receive-error: %v", err)
		return err
	}

	if strings.Split(message, ",")[0] != "completeir" {
		log.Printf("projector-response-error: %v", message)
		return errors.New(message)
	}

	return nil
}

func turnProjectorOn() error {
	projectorMutex.Lock()
	defer projectorMutex.Unlock()
	return sendProjectorFunc(func(conn net.Conn) error {
		err := sendProjectorCmd(conn, fmt.Sprintf("sendir,%v:1,1,38000,1,1,343,171,21,64,21,64,21,21,21,21,21,21,21,21,21,21,21,64,21,64,21,21,21,64,21,21,21,64,21,21,21,64,21,21,21,21,21,21,21,21,21,21,21,64,21,21,21,21,21,64,21,64,21,64,21,64,21,64,21,21,21,64,21,64,21,21,21,1546\r", projectorIrPort))
		return err
	})
}

func turnProjectorOff() error {
	projectorMutex.Lock()
	defer projectorMutex.Unlock()
	return sendProjectorFunc(func(conn net.Conn) error {
		err := sendProjectorCmd(conn, fmt.Sprintf("sendir,%v:1,1,38000,1,1,343,171,21,64,21,64,21,21,21,21,21,21,21,21,21,21,21,64,21,64,21,21,21,64,21,21,21,64,21,21,21,64,21,21,21,64,21,21,21,21,21,21,21,64,21,21,21,21,21,64,21,21,21,64,21,64,21,64,21,21,21,64,21,64,21,21,21,1546\r", projectorIrPort))
		return err
	})
}

func toggleProjectorPower() error {
	projectorMutex.Lock()
	defer projectorMutex.Unlock()
	return sendProjectorFunc(func(conn net.Conn) error {
		err := sendProjectorCmd(conn, fmt.Sprintf("sendir,%v:1,1,38000,1,1,347,173,22,65,22,65,22,22,22,22,22,22,22,22,22,22,22,65,22,65,22,22,22,65,22,22,22,65,22,22,22,65,22,22,22,22,22,22,22,22,22,22,22,65,22,22,22,22,22,65,22,65,22,65,22,65,22,65,22,22,22,65,22,65,22,22,22,1570\r", projectorIrPort))
		return err
	})
}
