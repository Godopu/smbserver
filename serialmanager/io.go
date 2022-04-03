package serialmanager

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

func recv(port io.Reader, h func(e Event)) error {
	reader := bufio.NewReader(port)

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println("USB is disconnected")
				
				return err
			}
		}
		recvObj := map[string]interface{}{}
		err = json.Unmarshal(b, &recvObj)

		if err == nil && h != nil {

			h(NewEvent(recvObj, "recv"))
		}
	}
}
