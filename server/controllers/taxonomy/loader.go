package taxonomy

import (
	"darco/proto/models/taxonomy"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type State map[int]taxonomy.ImportProcess

var state State = make(State)

type EventServer struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan taxonomy.ImportProcess

	// New client connections
	NewClients chan chan State

	// Closed client connections
	ClosedClients chan chan State

	// Total client connections
	TotalClients map[chan State]bool

	terminate chan struct{}
}

type ClientChan chan State

func (stream *EventServer) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			log.Printf("Last message : %v", state)
			client <- state
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))
			close(client)

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			log.Printf("MSG : %v", eventMsg)
			if !state[eventMsg.GBIF_ID].Done {
				state[eventMsg.GBIF_ID] = eventMsg
				for clientMessageChan := range stream.TotalClients {
					clientMessageChan <- state
				}
			}
			if eventMsg.Done {
				delete(state, eventMsg.GBIF_ID)
			}
		case <-stream.terminate:
			return
		}
	}
}

func (stream *EventServer) monitor(p *taxonomy.ImportProcess) {
	log.Printf("TRIGGER")
	// p.Imported += c
	stream.Message <- *p
}

func NewServer() (event *EventServer) {
	event = &EventServer{
		Message:       make(chan taxonomy.ImportProcess),
		NewClients:    make(chan chan State),
		ClosedClients: make(chan chan State),
		TotalClients:  make(map[chan State]bool),
		terminate:     make(chan struct{}),
	}
	return
}

type Controller struct {
	Endpoint        func(*gin.Context)
	ProgressTracker func(*gin.Context)
}

func UpdateTaxonomyDB() Controller {
	var stream = NewServer()

	endpoint := func(c *gin.Context) {
		log.Printf("State : %v", state)

		var taxon taxonomy.TaxonGBIF
		if c.ShouldBindJSON(&taxon) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid taxon search terms"})
			return
		}
		log.Printf("TAXON : %v", taxon)
		go stream.listen()

		go taxonomy.ImportTaxon(taxon.Key, stream.monitor)

		c.JSON(http.StatusAccepted, nil)
	}

	tracker := func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		clientChan := make(ClientChan)
		stream.NewClients <- clientChan

		go func() {
			<-c.Request.Context().Done()
			stream.ClosedClients <- clientChan
		}()

		c.Stream(func(writer io.Writer) bool {
			msg, ok := <-clientChan
			if ok {
				c.SSEvent("download", msg)
				return true
			}
			return false
		})
	}

	return Controller{
		Endpoint:        endpoint,
		ProgressTracker: tracker,
	}
}
