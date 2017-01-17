package backend

import (
	"net/http/httputil"
	"net/http"
	"encoding/json"
	"bufio"

	"github.com/mijara/statspout/log"
)

type Event struct {
	Type   string `json:"Type"`
	Action string `json:"Action"`
	Actor struct {
		Attributes struct {
			Name string `json:"name"`
			OldName string `json:"oldName,omitempty"`
		} `json:"Attributes"`
	} `json:"Actor"`
}

type EventsMonitor struct {
	client *httputil.ClientConn
	quit   chan bool
}

func NewEventsMonitor(http bool, address string) (*EventsMonitor, error) {
	conn, err := createConn(http, address)
	if err != nil {
		return nil, err
	}

	return &EventsMonitor{
		client: httputil.NewClientConn(conn, nil),
	}, nil
}

func (em *EventsMonitor) monitor(containers map[string]bool) {
	em.quit = make(chan bool, 1)
	go em.loop(containers)
}

func (em *EventsMonitor) Close() {
	em.quit <- true
}

func (em *EventsMonitor) loop(containers map[string]bool) {
	req, err := http.NewRequest("GET", "/events", nil)
	if err != nil {
		log.Error.Printf("Could not monitor events: %s", err.Error())
	}

	res, err := em.client.Do(req)
	if err != nil {
		log.Error.Printf("Events request failed: %s", err.Error())
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)

	for {
		select {
		case <-em.quit:
			return
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				log.Error.Printf("Events response error: %s", err.Error())
			}

			event := Event{}
			err = json.Unmarshal(line, &event)
			if err != nil {
				log.Error.Printf("Events response error: %s", err.Error())
			}

			if event.Type == "container" {
				switch event.Action {
				case "stop":
					log.Info.Printf("Container %s stopped.", event.Actor.Attributes.Name)

					delete(containers, event.Actor.Attributes.Name)
				case "start":
					log.Info.Printf("Container %s started.", event.Actor.Attributes.Name)

					containers[event.Actor.Attributes.Name] = true
				case "rename":
					oldName := event.Actor.Attributes.OldName[1:]
					log.Info.Printf("Container %s renamed to %s.", oldName, event.Actor.Attributes.Name)

					delete(containers, oldName)
					containers[event.Actor.Attributes.Name] = true
				}
			}
		}
	}
}
