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

// Player represents the structure of a Spyfall player.
type Player struct {
	Name     string
	Role     string
	Location string
}

type Location struct {
	Name string
	Role []string
}

type Hourglass struct {
	start    time.Time
	duration time.Duration
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
	// Randomize the order of the locations in the list.
	rand.Shuffle(len(locations), func(i, j int) {
		locations[i], locations[j] = locations[j], locations[i]
	})

	// Choose a location for this game, the first index of the now-shuffled elements.
	location := locations[0]

	roles := location.Role

	// TODO: support custom duration of time specified by the admin of the lobby

	timers := map[int]Hourglass{
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

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		t, e := template.ParseFiles("tmpl/example.html")
		// Check for errors.
		if e != nil {
			panic(err)
		}
		// Randomize the order of the roles in the location.
		rand.Shuffle(len(roles), func(i, j int) {
			roles[i], roles[j] = roles[j], roles[i]
		})
		role := roles[0]

		// Create a player with the name provided in the HTTP request,
		// and assign the player a random role from the now-shuffled list.
		p := Player{
			Name:     r.FormValue("name"),
			Role:     role,
			Location: location.Name,
		}
		e = t.Execute(w, p)
		// e = t.Execute(os.Stdout, p)
		if e != nil {
			panic(err)
		}
	})
	err = http.ListenAndServe("[::1]:1337", nil)
	if err != nil {
		panic(err)
	}
}
