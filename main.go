package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	port = 1337
)

// Player represents the structure of a Spyfall player.
type Player struct {
	Name string
	Role string
}

// Location represents the current setting of the current game of Spyfall.
type Location struct {
	Name string
	Role []string
}

// Clock represents the current status of the in-game timer. When there is no
// game in process, "start" will be set to the zero-value of time.Time. For now,
// the value of duration is set to the default value of 5 minutes.
type Clock struct {
	start    time.Time
	duration time.Duration
}

// Lobby represents a Spyfall lobby, in which there is a host, a list of
// players, and potentially, an active game taking place.
type Lobby struct {
	location Location
	clock    Clock
	player   []Player
}

func main() {
	// Read in the list of locations and roles from a file.
	data, err := os.ReadFile("data/location.yaml")
	if err != nil {
		panic(err)
	}
	var locations []Location
	err = yaml.Unmarshal(data, &locations)
	if err != nil {
		panic(err)
	}
	// Choose a location for this game from the list of available locations.
	l := locations[rand.Intn(len(locations))]

	// TODO: support custom duration of time specified by the admin of the lobby

	timers := map[int]Clock{
		1: {
			start:    time.Now(),
			duration: time.Minute * 5,
		},
	}

	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		id, e := strconv.Atoi(r.FormValue("id"))
		if e != nil {
			panic(e)
		}
		t, ok := timers[id]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Invalid timer id %d", id)
			return
		}
		remaining := time.Until(t.start.Add(t.duration))
		_, e = fmt.Fprintf(w, "%v\n", remaining)
		if e != nil {
			panic(e)
		}
	})

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Connection", "keep-alive")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(
				w,
				"Server does not support Flusher!",
				http.StatusInternalServerError,
			)
			return
		}
		for {
			time.Sleep(1 * time.Second)
			_, err = fmt.Fprintf(w, "data: logging %v...", time.Now())
			if err != nil {
				// TODO: Handle the "broken pipe" error that is thrown when a
				//  user prematurely leaves the lobby (exits the page).
				panic(err)
			}
			_, err = fmt.Fprintf(w, "\n\n")
			if err != nil {
				panic(err)
			}
			flusher.Flush()
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the template file.
		t, e := template.ParseFiles("tmpl/index.html")
		if e != nil {
			panic(err)
		}
		// Select a random role among the roles available in this location.
		role := l.Role[rand.Intn(len(l.Role))]

		// Extract the name from the HTTP request.
		name := r.FormValue("name")

		// Create a player with the name provided in the HTTP request,
		// and assign the player a random role from the now-shuffled list.
		p := Player{
			Name: name,
			Role: role,
		}

		// Output the templated HTML file in response to this HTTP request.
		e = t.Execute(w, p)
		if e != nil {
			panic(err)
		}
	})

	// Format the hostname and port number to use for the HTTP server.
	authority := fmt.Sprintf(":%d", port)

	// Listen for incoming HTTP requests.
	err = http.ListenAndServe(authority, nil)
	if err != nil {
		panic(err)
	}
}

func init() {
	// Attempt to read a custom port number, if one is provided in the form
	// of an environment variable.
	val := os.Getenv("SPYFALL_PORT")
	if val != "" {
		// Only alter the value of the port number if it is a valid integer.
		num, err := strconv.Atoi(val)
		if err == nil {
			port = num
		} else {
			fmt.Fprintf(
				os.Stderr,
				"Error parsing port from environment variable 'SPYFALL_PORT'... Listening on port %d instead\n",
				port,
			)
		}
	}
}
