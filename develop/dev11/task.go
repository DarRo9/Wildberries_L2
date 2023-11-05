/*
HTTP-сервер

Реализовать HTTP-сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP-библиотекой.


В рамках задания необходимо:
Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
Реализовать middleware для логирования запросов


Методы API: 
POST /create_event 
POST /update_event 
POST /delete_event 
GET /events_for_day 
GET /events_for_week 
GET /events_for_month


Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09). В GET методах параметры передаются через queryString, 
в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON-документ содержащий либо {"result": "..."} в случае успешного выполнения метода, 
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
Реализовать все методы.
Бизнес логика НЕ должна зависеть от кода HTTP сервера.
В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен 
возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге 
и выводить в лог каждый обработанный запрос.
*/
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Event struct {
	UserID      int    `json:"user_id"`
	ID          int    `json:"id"`
	Date        Date                        `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ConcreteEvent структура для получения конкретного события
type ConcreteEvent struct {
	UserID int `json:"user_id"`
	ID     int `json:"id"`
}

type EventRepository struct {
	m map[int][]Event
	*sync.RWMutex
}

type Date struct {
	date time.Time
}

type Logger struct {
	*log.Logger
}

// Scope - сервер, логер и хранилище
type Scope struct {
	srv        *http.ServeMux
	logger     Logger
	EventRepository EventRepository
}

func CreateScope() *Scope {
	return &Scope{
		srv: http.NewServeMux(),
		logger: Logger{
			log.New(os.Stdout, "logger: ", log.Lshortfile),
		},
		EventRepository: EventRepository{
			m:       make(map[int][]Event),
			RWMutex: new(sync.RWMutex),
		},
	}
}

// startingServer настраивает роуты и запускает сервер
func (scope *Scope) startingServer() {
	scope.srv.HandleFunc("/create_event", scope.CreateEvent)
	scope.srv.HandleFunc("/update_event", scope.UpdateEvent)
	scope.srv.HandleFunc("/delete_event", scope.RemoveEvent)

	scope.srv.HandleFunc("/events_for_day", scope.DayEvents)
	scope.srv.HandleFunc("/events_for_week", scope.WeekEvents)
	scope.srv.HandleFunc("/events_for_month", scope.MonthEvents)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), scope.srv))
}

// CreateNewEvent создает новое событие и сохраняет его в хранилище
func (scope *Scope) CreateNewEvent(event Event) error {
	scope.EventRepository.RWMutex.Lock()
	defer scope.EventRepository.RWMutex.Unlock()
	if scope.checkEvent(event) {
		return errors.New("duplicate event not allowed")
	} else {
		scope.EventRepository.m[event.UserID] = append(scope.EventRepository.m[event.UserID], event)
	}
	return nil
}

func (scope *Scope) CreateEvent(w http.ResponseWriter, r *http.Request) {
	scope .logger.Println(r.URL)
	if r.Method != http.MethodPost {
		sendErr(w, "Not correct method", http.StatusBadRequest)
		return
	}

	event, err := parseJSON(r)

	if err != nil {
		sendErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := scope.CreateNewEvent(event); err != nil {
		sendErr(w, err.Error(), http.StatusBadRequest)
	}
	sendRes(w, "Success", []Event{event}, http.StatusCreated)
}

func (scope  *Scope) UpdateEventFunc(e Event) error {
	scope .EventRepository.RWMutex.Lock()
	defer scope .EventRepository.RWMutex.Unlock()

	events := scope .EventRepository.m[e.UserID]
	for ind := 0; ind < len(events); ind++ {
		if events[ind].ID == e.ID {
			events[ind].Title = e.Title
			events[ind].Description = e.Description
			events[ind].Date = e.Date
			return nil
		}
	}
	return errors.New("event not found")
}

func (scope *Scope) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	scope.logger.Println(r.URL)
	if r.Method != http.MethodPost {
		sendErr(w, "Not correct method", http.StatusBadRequest)
		return
	}
	event, err := parseJSON(r)
	if err != nil {
		sendErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ValidEvents(event) {
		err = scope.UpdateEventFunc(event)
		if err != nil {
			sendErr(w, "Event not found", http.StatusInternalServerError)
		} else {
			sendRes(w, "Success", []Event{event}, http.StatusOK)
		}
	} else {
		sendErr(w, "Not valid event", http.StatusBadRequest)
	}
}

func (scope *Scope) RemoveEventFunc(cEvent ConcreteEvent) error {
	scope.EventRepository.RWMutex.Lock()
	defer scope.EventRepository.RWMutex.Unlock()

	events := scope.EventRepository.m[cEvent.UserID]
	for ind := 0; ind < len(events); ind++ {
		if events[ind].ID == cEvent.ID {
			scope.EventRepository.m[cEvent.UserID] = append(scope.EventRepository.m[cEvent.UserID][0:ind], scope.EventRepository.m[cEvent.UserID][ind+1:]...)
			return nil
		}
	}
	return errors.New("event not found")
}

func (scope *Scope) RemoveEvent(w http.ResponseWriter, r *http.Request) {
	scope.logger.Println(r.URL)
	if r.Method != http.MethodPost {
		sendErr(w, "Not correct method", http.StatusBadRequest)
		return
	}

	var cEvent ConcreteEvent
	err := json.NewDecoder(r.Body).Decode(&cEvent)
	if err != nil {
		sendErr(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = scope.RemoveEventFunc(cEvent)
	if err != nil {
		sendErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendRes(w, "Success", nil, http.StatusOK)
}

// sendErr отправляет заданную ошибку с заданным кодов статуса
func sendErr(writer http.ResponseWriter, errorStr string, status int) {
	response := struct {
		Error string `json:"error"`
	}{errorStr}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(jsonResponse)
}

// sendRes отправляет результат запроса
func sendRes(writer http.ResponseWriter, resStr string, events []Event, status int) {
	response := struct {
		Result string  `json:"result"`
		Events []Event `json:"events"`
	}{resStr, events}
	responseJson, err := json.Marshal(response)
	if err != nil {
		sendErr(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(responseJson)
}

// checkEvent проверяет содержится ли событие в хранилище
func (scope *Scope) checkEvent(e Event) bool {
	for _, v := range scope.EventRepository.m[e.UserID] {
		if v.ID == e.ID {
			return true
		}
	}
	return false
}

func ValidEvents(event Event) bool {
	if event.ID <= 0 || event.UserID <= 0 || event.Title == "" || event.Description == "" {
		return false
	}
	return true
}

func parseJSON(r *http.Request) (Event, error) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return event, errors.New("cannot decode json")
	}
	return event, nil
}

              
func (d *Date) UnmarshalJSON(input []byte) error {
	var err error
	d.date, err = time.Parse(`"2006-01-02"`, string(input))
	return err
}

// String для типа Date                       
func (d Date) String() string {
	return d.date.String()
}

// MarshalJSON для типа Date                       
func (d *Date) MarshalJSON() ([]byte, error) {
	dateStr := d.date.Format("2006-01-02")
	return json.Marshal(dateStr)
}

func (scope *Scope) DayEventsFunc(userID int, date time.Time) ([]Event, error) {
	scope.EventRepository.RWMutex.RLock()
	scope.EventRepository.RWMutex.RUnlock()

	var result []Event

	var allUserEvents []Event
	allUserEvents = scope.EventRepository.m[userID]
	if allUserEvents == nil {
		return nil, errors.New("unknown user_id")
	}

	for _, event := range allUserEvents {
		if event.Date.date.Year() == date.Year() &&
			event.Date.date.Month() == date.Month() &&
			event.Date.date.Day() == date.Day() {
			result = append(result, event)
		}
	}
	return result, nil
}

func (scope *Scope) DayEvents(w http.ResponseWriter, r *http.Request) {
	scope.logger.Println(r.URL)
	if r.Method != http.MethodGet {
		sendErr(w, "Not correct method", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		sendErr(w, "Incorrect args", http.StatusBadRequest)
		return
	}
	events, err := scope.DayEventsFunc(userID, date)
	if err != nil {
		sendErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendRes(w, "Success", events, http.StatusOK)
}

func (scope *Scope) WeekEventsFunc(userID int, date time.Time) ([]Event, error) {
	scope.EventRepository.RWMutex.RLock()
	scope.EventRepository.RWMutex.RUnlock()

	var result []Event

	var allUserEvents []Event
	allUserEvents = scope.EventRepository.m[userID]
	if allUserEvents == nil {
		return nil, errors.New("unknown user_id")
	}

	for _, event := range allUserEvents {
		difference := date.Sub(event.Date.date)
		if difference < 0 {
			difference = -difference
		}
		if difference <= time.Duration(7*24)*time.Hour {
			result = append(result, event)
		}
	}
	return result, nil
}

func (scope *Scope) WeekEvents(w http.ResponseWriter, r *http.Request) {
	scope.logger.Println(r.URL)
	if r.Method != http.MethodGet {
		sendErr(w, "Not correct method", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		sendErr(w, "Incorrect args", http.StatusBadRequest)
		return
	}
	events, err := scope.WeekEventsFunc(userID, date)
	if err != nil {
		sendErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendRes(w, "Success", events, http.StatusOK)
}

func (scope *Scope) MonthEventsFunc(userID int, date time.Time) ([]Event, error) {
	scope.EventRepository.RWMutex.RLock()
	scope.EventRepository.RWMutex.RUnlock()

	var result []Event

	var allUserEvents []Event
	allUserEvents = scope.EventRepository.m[userID]
	if allUserEvents == nil {
		return nil, errors.New("unknown user_id")
	}

	for _, event := range allUserEvents {
		if event.Date.date.Year() == date.Year() || event.Date.date.Month() == date.Month() {
			result = append(result, event)
		}
	}
	return result, nil
}

func (scope *Scope) MonthEvents(w http.ResponseWriter, r *http.Request) {
	scope.logger.Println(r.URL)
	if r.Method != http.MethodGet {
		sendErr(w, "Not correct method", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		sendErr(w, "Incorrect args", http.StatusBadRequest)
		return
	}
	events, err := scope.MonthEventsFunc(userID, date)
	if err != nil {
		sendErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendRes(w, "Success", events, http.StatusOK)
}

func main() {
	scope := CreateScope()
	scope.startingServer()
}
