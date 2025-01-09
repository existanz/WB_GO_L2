package database

import (
	"dev11/internal/entities"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Close() error
}

type service struct {
	db []entities.Event
}

var (
	dbInstance *service
)

func New() Service {

	if dbInstance != nil {
		return dbInstance
	}
	db := make([]entities.Event, 0)
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Close() error {
	log.Printf("Disconnected from database")
	return nil
}

func (s *service) AddEvent(event entities.Event) error {
	s.db = append(s.db, event)
	return nil
}

func (s *service) UpdateEvent(event entities.Event) error {
	// Update event logic
	return nil
}

func (s *service) DeleteEvent(id int) error {
	// Delete event logic
	return nil
}

func (s *service) GetEventsForDay(date string) ([]entities.Event, error) {
	// Get events for a specific day
	return []entities.Event{}, nil
}

func (s *service) GetEventsForWeek(date string) ([]entities.Event, error) {
	// Get events for a specific week
	return []entities.Event{}, nil
}

func (s *service) GetEventsForMonth(date string) ([]entities.Event, error) {
	// Get events for a specific month
	return []entities.Event{}, nil
}
