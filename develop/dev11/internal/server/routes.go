package server

import (
	"dev11/internal/entities"
	"dev11/internal/mw"
	"dev11/internal/service"
	"dev11/pkg/util"
	"fmt"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", s.CreateEvent)
	mux.HandleFunc("/update_event", s.UpdateEvent)
	mux.HandleFunc("/delete_event", s.DeleteEvent)
	mux.HandleFunc("/events_for_day", s.EventsForDay)
	mux.HandleFunc("/events_for_week", s.EventsForWeek)
	mux.HandleFunc("/events_for_month", s.EventsForMonth)

	loggingMux := mw.Logging(mux)

	return loggingMux
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	event, err := entities.ParseEvent(r)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := entities.ValidateEvent(event); err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Create event: ", event)

	if err := service.CreateEvent(event); err != nil {
		util.WriteError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	util.WriteResult(w, "Event created successfully")
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	// Similar implementation as CreateEvent
}

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	// Implementation
}

func (s *Server) EventsForDay(w http.ResponseWriter, r *http.Request) {
	// Implementation
}

func (s *Server) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	// Implementation
}

func (s *Server) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	// Implementation
}
