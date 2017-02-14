package backend

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/mijara/statspout/log"
)

type Event struct {
	Type   string `json:"Type"`
	Action string `json:"Action"`
	Actor  struct {
		Attributes struct {
			Name    string `json:"name"`
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

func (em *EventsMonitor) monitor(cli *Client, containers map[string]Container) {
	em.quit = make(chan bool, 1)
	go em.loop(cli, containers)
}

func (em *EventsMonitor) Close() {
	em.quit <- true
}

func (em *EventsMonitor) loop(cli *Client, containers map[string]Container) {
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
					cli.repo.Clear(event.Actor.Attributes.Name)

				case "start":
					log.Info.Printf("Container %s started.", event.Actor.Attributes.Name)

					// retrieve and store new container data.
					container, err := cli.RequestContainer(event.Actor.Attributes.Name)
					if err != nil {
						log.Error.Printf("Cannot retrieve container data for %s. Error: %s",
							event.Actor.Attributes.Name, err.Error())
					}
					containers[container.CanonicalName] = *container

				case "rename":
					oldName := event.Actor.Attributes.OldName[1:]
					log.Info.Printf("Container %s renamed to %s.", oldName, event.Actor.Attributes.Name)

					// delete registered container from map.
					delete(containers, oldName)
					cli.repo.Clear(oldName)

					// retrieve and store new container data.
					container, err := cli.RequestContainer(event.Actor.Attributes.Name)
					if err != nil {
						log.Error.Printf("Cannot retrieve container data for %s. Error: %s",
							event.Actor.Attributes.Name, err.Error())
					}
					fmt.Println(container)
					containers[container.CanonicalName] = *container
				}
			}
		}
	}
}
