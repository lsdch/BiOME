package taxonomy

import (
	"darco/proto/models/taxonomy"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type State map[int]taxonomy.ImportProcess

var state State = make(State)

type EventServer struct {
	Message       chan taxonomy.ImportProcess
	NewClients    chan chan State
	ClosedClients chan chan State
	TotalClients  map[chan State]bool
	Running       bool
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
			log.Debugf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			log.Debugf("Removed client. %d registered clients", len(stream.TotalClients))
			close(client)

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			log.Debugf("MSG : %v", eventMsg)
			if !state[eventMsg.GBIF_ID].Done {
				state[eventMsg.GBIF_ID] = eventMsg
				for clientMessageChan := range stream.TotalClients {
					clientMessageChan <- state
				}
			}
			// if eventMsg.Done {
			// 	delete(state, eventMsg.GBIF_ID)
			// }
		}
	}
}

func (stream *EventServer) monitor(p *taxonomy.ImportProcess) {
	stream.Message <- *p
}

func NewServer() (event *EventServer) {
	event = &EventServer{
		Message:       make(chan taxonomy.ImportProcess),
		NewClients:    make(chan chan State),
		ClosedClients: make(chan chan State),
		TotalClients:  make(map[chan State]bool),
		Running:       false,
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
		var taxon taxonomy.TaxonGBIF
		if err := c.ShouldBindJSON(&taxon); err != nil {
			log.Errorf("Invalid taxon definition to import as anchor \n%s", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid taxon definition", "error": err})
			return
		}
		if !stream.Running {
			go stream.listen()
		}

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
				c.SSEvent("progress", msg)
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
