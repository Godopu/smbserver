package serialmanager

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

func discoverDevice() (string, error) {
	fs, err := ioutil.ReadDir("/dev")
	if err != nil {
		return "", err
	}

	for _, f := range fs {
		if strings.Contains(f.Name(), "ttyACM") || strings.Contains(f.Name(), "ttyUSB") {
			return filepath.Join("/dev", f.Name()), nil
		}
	}

	return "", errors.New("not found device")
}

func WatchNewDevice(ctx context.Context) (string, error) {
	log.Println("Watching device")
	filter := make(chan notify.EventInfo, 1)
	if err := notify.Watch("/dev", filter, notify.Create); err != nil {
		return "", err
	}
	defer notify.Stop(filter)

	for {
		select {
		case <-ctx.Done():
			return "", nil
		case e := <-filter:
			if strings.Contains(e.Path(), "/dev/ttyACM") || strings.Contains(e.Path(), "/dev/ttyUSB") {
				return e.Path(), nil

			}
		}
	}
}
