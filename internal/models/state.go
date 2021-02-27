package models

import "time"

// State stores the run time state for execution of
// the  song finder
type State struct {
	Completed       time.Time
	Screenshots     map[string]*Screenshot
	SoftwareVersion string
}
