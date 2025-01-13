package database

import (
	"dev11/internal/entities"
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	AddEvent(event entities.Event) error
	UpdateEvent(event entities.Event) error
	DeleteEvent(userID, eventID int) error

	GetEventsForPeriod(userID int, date string, period int) ([]entities.Event, error)

	Close() error
}

type service struct {
	mu     *sync.Mutex
	db     map[int][]entities.Event
	nextID int
}

var (
	dbInstance *service
)

func New() Service {

	if dbInstance != nil {
		return dbInstance
	}

	dbInstance = &service{
		mu:     &sync.Mutex{},
		db:     make(map[int][]entities.Event, 0),
		nextID: 1,
	}
	return dbInstance
}

func (s *service) Close() error {
	log.Printf("Disconnected from database")
	s.db = nil
	return nil
}

func (s *service) AddEvent(event entities.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = s.nextID
	s.nextID++

	userID := event.UserID
	s.db[userID] = append(s.db[userID], event)

	return nil
}

func (s *service) UpdateEvent(event entities.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID := event.UserID
	events, ok := s.db[userID]
	if !ok {
		return entities.ErrUserNotFound
	}

	for i, e := range events {
		if e.ID == event.ID {
			events[i] = event
			return nil
		}
	}

	return entities.ErrEventNotFound
}

func (s *service) DeleteEvent(userID, eventID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	events, ok := s.db[userID]
	if !ok {
		return entities.ErrUserNotFound
	}

	for i, e := range events {
		if e.ID == eventID {
			s.db[userID] = append(events[:i], events[i+1:]...)
			return nil
		}
	}

	return entities.ErrEventNotFound
}

func (s *service) GetEventsForPeriod(userID int, date string, period int) ([]entities.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events, ok := s.db[userID]
	if !ok {
		return nil, entities.ErrUserNotFound
	}

	t, err := time.Parse("2006-01-02", date)
	t2 := t.AddDate(0, 0, period)
	if err != nil {
		return nil, err
	}

	var result []entities.Event
	for _, event := range events {
		if event.Date.After(t) && event.Date.Before(t2) {
			result = append(result, event)
		}
	}

	return result, nil
}
