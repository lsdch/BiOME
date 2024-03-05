package taxonomy

import (
	gbif "darco/proto/models/taxonomy/GBIF"
	"io"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
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

// Controller represents an API controller with two intertwined endpoints/
type Controller struct {
	// Endpoint handles requests to import a clade from GBIF.
	Endpoint func(*gin.Context, *edgedb.Client)
	// ProgressTracker monitors the progress of an import process, using Server-Sent Events.
	ProgressTracker func(*gin.Context)
}

// @Summary Import GBIF clade
// @Description Imports a clade from the GBIF taxonomy, using a its GBIF ID
// @id ImportGBIF
// @tags Taxonomy
// @Accept json
// @Produce json
// @Success 202
// @Failure 403
// @Failure 400
// @Router /taxonomy/import [put]
// @Param data body gbif.ImportRequestGBIF true "Import parameters"
func ImportCladeGBIF() Controller {
	var stream = NewServer()

	endpoint := func(c *gin.Context, db *edgedb.Client) {
		target := gbif.ImportRequestGBIF{}
		err := c.BindJSON(&target)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Failed to parse import parameters from URL query parameters",
			})
			return
		}

		log.Infof("Received GBIF import request : %v", target)
		if !stream.Running {
			go stream.listen()
		}

		log.Infof("Started import of taxon : [GBIF %d] with children: %v", target.Key, target.Children)
		go gbif.ImportTaxon(db, target, stream.monitor)
		c.Status(http.StatusAccepted)
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
