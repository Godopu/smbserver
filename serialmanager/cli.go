package serialmanager

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/jacobsa/go-serial/serial"
)

var ctx context.Context
var cancel context.CancelFunc

var discoverHandleFunc func(e Event) = nil
var recvHandleFunc func(e Event) = nil

var onConnected func(e Event) = nil
var onDisconnected func(e Event) = nil

var iface string
var port io.ReadWriteCloser

func AddDiscoverHandleFunc(h func(e Event)) {
	discoverHandleFunc = h
}

func AddRecvHandleFunc(h func(e Event)) {
	recvHandleFunc = h
}

func AddOnConnected(h func(e Event)) {
	onConnected = h
}

func AddOnDisconnected(h func(e Event)) {
	onDisconnected = h
}

func Close() {
	defer port.Close()
	cancel()
}

func init() {
	ctx, cancel = context.WithCancel(context.Background())
}

func Run() error {
	var err error
	for {
		iface, err = discoverDevice()
		if err != nil {
			if err.Error() == "not found device" {
				iface, err = WatchNewDevice(ctx)
			}

			if err != nil {
				return err
			}
		}
		if discoverHandleFunc != nil {
			discoverHandleFunc(NewEvent(map[string]interface{}{"iface": iface}, "discovered"))
		}

		fmt.Println("initDevice")
		err = initDevice()
		if err != nil {
			log.Println(err)
			port.Close()
			continue
		}

		fmt.Println("recv")
		err := recv(port, recvHandleFunc)
		if err.Error() == "EOF" {
			if onDisconnected != nil {
				onDisconnected(NewEvent(map[string]interface{}{"iface": iface}, "disconnected"))
			}
			iface = ""
		} else {
			panic(err)
		}
	}
}

func initDevice() error {
	var err error
	// err = changePermission(iface)
	// if err != nil {
	// 	return err
	// }

	options := serial.OpenOptions{
		PortName:        iface,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 16,
	}

	port, err = serial.Open(options)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(port)
	// encoder := json.NewEncoder(port)

	sndMsg := map[string]interface{}{}
	sndMsg["code"] = 100

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		rcvMsg := map[string]interface{}{}
		err = json.Unmarshal(b, &rcvMsg)
		fmt.Println(string(b))

		if err != nil {
			continue
		}

		code, ok := rcvMsg["code"].(float64)
		if ok && code == 100 {
			if onConnected != nil {
				onConnected(NewEvent(map[string]interface{}{"iface": iface}, "connected"))
			}
			break
		}
		err = Write(sndMsg)
		if err != nil {
			return err
		}
		// encoder.Encode(sndMsg)
	}

	return nil
}

func Write(obj interface{}) error {
	if port != nil {
		return errors.New("device is not connected")
	}
	enc := json.NewEncoder(port)
	err := enc.Encode(obj)
	if err != nil {
		return err
	}
	return nil
}

// func changePermission(iface string) error {
// 	log.Println("changing the mod of file")

// 	cmd := exec.Command("chmod", "a+rw", iface)
// 	_, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
