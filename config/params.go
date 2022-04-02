package config

import (
	"errors"
	"net"
	"os"
	"strings"

	"github.com/magiconair/properties"
)

const Mode = "debug"

var Params = map[string]interface{}{}
var MyIP string

func init() {

	MyIP = GetIP()
	if _, err := os.Stat("./config.properties"); errors.Is(err, os.ErrNotExist) {
		createInitFile()
	}

	p := properties.MustLoadFile("./config.properties", properties.UTF8)
	Params["bind"] = p.GetString("bind", ":4000")
}

func createInitFile() {
	f, err := os.Create("./config.properties")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := properties.NewProperties()
	p.SetValue("bind", ":4000")
	p.Write(f, properties.UTF8)
}

func GetIP() string {
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		panic(err)
	}

	for _, i := range ifaces {
		if strings.Index(i.Name, "en") == 0 || strings.Index(i.Name, "et") == 0 {
			addrs, err := i.Addrs()
			if err != nil {
				panic(err)
			}
			// handle err

			for _, addr := range addrs {
				if idx := strings.IndexByte(addr.String(), '/'); idx != -1 &&
					!strings.Contains(addr.String(), ":") {
					return addr.String()[:idx]
				}
				// process IP address
			}
		}
	}
	return ""
}

func Set(key, value string) {

	var p *properties.Properties
	if _, err := os.Stat("./config.properties"); errors.Is(err, os.ErrNotExist) {
		p = properties.NewProperties()
	} else {
		p = properties.MustLoadFile("./config.properties", properties.UTF8)
		os.Remove("./config.properties")
	}
	f, err := os.Create("./config.properties")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p.SetValue(key, value)
	p.Write(f, properties.UTF8)
	Params[key] = value
}
