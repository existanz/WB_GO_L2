package server

import (
	"dev11/internal/entities"
	"dev11/internal/mw"
	"dev11/pkg/util"
	"net/http"
	"strconv"
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

	if err := s.db.AddEvent(event); err != nil {
		util.WriteError(w, http.StatusServiceUnavailable, entities.ErrServiceUnavailable.Error())
		return
	}

	util.WriteResult(w, "Event created successfully")
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
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

	if err := s.db.UpdateEvent(event); err != nil {
		if err == entities.ErrUserNotFound || err == entities.ErrEventNotFound {
			util.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		util.WriteError(w, http.StatusServiceUnavailable, entities.ErrServiceUnavailable.Error())
		return
	}

	util.WriteResult(w, "Event updated successfully")
}

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	eventID, err := strconv.Atoi(r.FormValue("event_id"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "event_id should be a valid positive integer")
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "user_id should be a valid positive integer")
		return
	}

	if err := s.db.DeleteEvent(userID, eventID); err != nil {
		if err == entities.ErrUserNotFound || err == entities.ErrEventNotFound {
			util.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		util.WriteError(w, http.StatusServiceUnavailable, entities.ErrServiceUnavailable.Error())
		return
	}

	util.WriteResult(w, "Event deleted successfully")
}

func (s *Server) EventsForDay(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	if date == "" {
		util.WriteError(w, http.StatusBadRequest, "date should be provided")
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "user_id should be a valid positive integer")
		return
	}

	events, err := s.db.GetEventsForPeriod(userID, date, 1)
	if err != nil {
		if err == entities.ErrUserNotFound || err == entities.ErrEventNotFound {
			util.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		util.WriteError(w, http.StatusServiceUnavailable, entities.ErrServiceUnavailable.Error())
		return
	}

	util.WriteJSON(w, events)
}

func (s *Server) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	if date == "" {
		util.WriteError(w, http.StatusBadRequest, "date should be provided")
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "user_id should be a valid positive integer")
		return
	}

	events, err := s.db.GetEventsForPeriod(userID, date, 7)
	if err != nil {
		if err == entities.ErrUserNotFound || err == entities.ErrEventNotFound {
			util.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		util.WriteError(w, http.StatusServiceUnavailable, entities.ErrServiceUnavailable.Error())
		return
	}

	util.WriteJSON(w, events)
}

func (s *Server) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	if date == "" {
		util.WriteError(w, http.StatusBadRequest, "date should be provided")
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "user_id should be a valid positive integer")
		return
	}

	events, err := s.db.GetEventsForPeriod(userID, date, 30)
	if err != nil {
		if err == entities.ErrUserNotFound || err == entities.ErrEventNotFound {
			util.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		util.WriteError(w, http.StatusServiceUnavailable, entities.ErrServiceUnavailable.Error())
		return
	}

	util.WriteJSON(w, events)
}
