package main

import (
	"sync"
)

func turnTheaterOn() error {
	var wg sync.WaitGroup
	wg.Add(3)
	var err1, err2, err3 error
	go func() {
		err1 = setLightsToDiningMode()
		wg.Done()
	}()
	go func() {
		//err2 = turnProjectorOn()
		wg.Done()
	}()
	go func() {
		err3 = turnReceiverOn()
		wg.Done()
	}()
	wg.Wait()
	return first(err1, err2, err3)
}

func turnTheaterOff() error {
	var wg sync.WaitGroup
	wg.Add(3)
	var err1, err2, err3 error
	go func() {
		err1 = setLightsToFullOn()
		wg.Done()
	}()
	go func() {
		//err2 = turnProjectorOff()
		wg.Done()
	}()
	go func() {
		err3 = turnReceiverOff()
		wg.Done()
	}()
	wg.Wait()
	return first(err1, err2, err3)
}
