package taxonomy

import (
	gbif "darco/proto/models/taxonomy/GBIF"

	log "github.com/sirupsen/logrus"
)

type GBIF_ID int

// State represents the current state of import processes
type State map[GBIF_ID]gbif.ImportProcess

var state State = make(State)

// EventServer manages the communication between clients and the import process.
type EventServer struct {
	Message       chan gbif.ImportProcess
	NewClients    chan chan State
	ClosedClients chan chan State
	TotalClients  map[chan State]bool
	Running       bool
}

func (stream *EventServer) AddClient() ClientChan {
	clientChan := make(ClientChan)
	stream.NewClients <- clientChan
	return clientChan
}

type ClientChan chan State

func (stream *EventServer) listen() {
	stream.Running = true
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			client <- state
			stream.TotalClients[client] = true
			log.Debugf(
				"Client added. %d registered clients",
				len(stream.TotalClients),
			)

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			log.Debugf(
				"Removed client. %d registered clients",
				len(stream.TotalClients),
			)
			close(client)

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			log.Debugf("MSG : %v", eventMsg)
			GBIF_ID := GBIF_ID(eventMsg.GBIF_ID)
			if !state[GBIF_ID].Done {
				state[GBIF_ID] = eventMsg
				for clientMessageChan := range stream.TotalClients {
					clientMessageChan <- state
				}
			}
		}
	}
}

// monitor sends an import process event to the server.
func (stream *EventServer) monitor(p *gbif.ImportProcess) {
	stream.Message <- *p
}

func NewServer() (event *EventServer) {
	event = &EventServer{
		Message:       make(chan gbif.ImportProcess),
		NewClients:    make(chan chan State),
		ClosedClients: make(chan chan State),
		TotalClients:  make(map[chan State]bool),
		Running:       false,
	}
	return
}
