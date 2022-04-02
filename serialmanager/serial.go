package serialmanager

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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
		if strings.Contains(f.Name(), "ttyACM") {
			return filepath.Join("/dev", f.Name()), nil
		}
	}

	return "", errors.New("not found device")
}

func WatchNewDevice(ctx context.Context, ch_discover chan<- notify.EventInfo) error {
	defer close(ch_discover)

	filter := make(chan notify.EventInfo, 1)
	if err := notify.Watch("/dev", filter, notify.Create); err != nil {
		return err
	}
	defer notify.Stop(filter)

	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-filter:
			if strings.Contains(e.Path(), "/dev/ttyACM") {
				fmt.Println(e.Path())
				ch_discover <- e
			}
		}
	}
}
