package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
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
	// Register handlers for requests to different paths.
	http.HandleFunc("/time", handleTime)
	http.HandleFunc("/events", handleEvent)
	http.HandleFunc("/", handleDefault)
	// http.Handle("/public", http.FileServer(http.Dir("./public")))

	// Format the server's hostname and port number.
	authority := fmt.Sprintf(":%d", port)

	// Listen for incoming HTTP requests.
	fmt.Printf("Mixing it up at http://localhost:%d\n", port)
	err := http.ListenAndServe(authority, nil)
	if err != nil {
		panic(err)
	}
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	// Handle requests to the homepage "index.html", by responding with
	// a formatted version of the templated "index.html" file.
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		serveIndex(w, r)
		return
	}
	// Create a filesystem rooted at the "/public" subdirectory.
	fsys := os.DirFS("./public")
	file, err := fsys.Open(r.URL.Path[1:])
	// Notify the user if the requested file is not found.
	if errors.Is(err, fs.ErrNotExist) {
		http.NotFound(w, r)
	} else if err != nil {
		// Panic if any other error is encountered.
		panic(err)
	} else {
		// If there is no error, the file was found, so output write it as
		// the response to this request.
		_, err = io.Copy(w, file)
		// Panic if we were unable to do so.
		if err != nil {
			panic(err)
		}
	}
}

// Responds to a request by providing a formatted version of the templated
// "index.html" file.
func serveIndex(w http.ResponseWriter, r *http.Request) {
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
	return
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	// Set the headers necessary to output a server-side event.
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
	// Output a "data" segment every second.
	// TODO: provide more useful server-side events, such as
	//  user renames, users joining the lobby, and users leaving the lobby.
	for {
		time.Sleep(1 * time.Second)
		_, err := fmt.Fprintf(w, "data: logging %v...", time.Now())
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
}

// TODO: Utilize this request on the front-end.
func handleTime(w http.ResponseWriter, r *http.Request) {
	timers := map[int]Clock{
		1: {
			start:    time.Now(),
			duration: time.Minute * 5,
		},
	}

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
