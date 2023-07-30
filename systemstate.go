package main

// This could be a lot better by having per-client reading/writing

import (
	"time"
)

type Service struct { // Is this not a protobuff
	Name           string
	IPAddress      string
	Status         string
	ParentService  string
	LastConnection time.Time
}
