package entities

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

func ParseEvent(r *http.Request) (Event, error) {
	var event Event

	err := r.ParseForm()
	if err != nil {
		return event, err
	}

	event.UserID, err = strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return event, err
	}
	event.Title = r.FormValue("title")
	event.Description = r.FormValue("description")
	event.Date, err = time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		return event, err
	}
	event.Date = event.Date.Add(7 * time.Hour)
	event.Duration, err = time.ParseDuration(r.FormValue("duration"))
	if err != nil {
		return event, err
	}
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	return event, nil
}

func ParseEventRef(r *http.Request) (Event, error) {
	event := Event{}

	err := r.ParseForm()
	if err != nil {
		return event, err
	}

	err = parseFormIntoStruct(r.Form, &event)
	if err != nil {
		return event, err
	}

	return event, nil
}

func parseFormIntoStruct(form url.Values, obj interface{}) error {
	const formTag = "form"
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := reflect.New(field.Type)

		switch field.Type.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			valueStr := form.Get(field.Tag.Get(formTag))
			if valueStr == "" {
				continue
			}
			valueInt, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return err
			}
			value.SetInt(valueInt)
		case reflect.String:
			valueStr := form.Get(field.Tag.Get(formTag))
			if valueStr == "" {
				continue
			}
			value.SetString(valueStr)
		case reflect.Struct:
			switch field.Type.Name() {
			case "Time":
				valueStr := form.Get(field.Tag.Get(formTag))
				if valueStr == "" {
					continue
				}
				valueTime, err := time.Parse("2006-01-02 15:04:05", valueStr)
				if err != nil {
					return err
				}
				value.Set(reflect.ValueOf(valueTime))
			case "Duration":
				valueStr := form.Get(field.Tag.Get(formTag))
				if valueStr == "" {
					continue
				}
				valueDuration, err := time.ParseDuration(valueStr)
				if err != nil {
					return err
				}
				value.Set(reflect.ValueOf(valueDuration))
			}
		}

		v.Field(i).Set(value.Elem())
	}

	return nil
}
