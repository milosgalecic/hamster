package main

import (
	"hamster/internal/models"
	"html/template"
	"path/filepath"
)

type JobDisplay struct {
	Job               *models.Job
	DriverName        string
	TruckLicensePlate string
	FormattedDate     string
	FormattedArrival  string
	IsEdit            bool
	DriverIdMap       map[int]string
	TruckIdMap        map[int]string
}

// Nemam vise Job uvek vracam slice
type templateData struct {
	CurrentYear     int
	DriverIdMap     map[int]string
	TruckIdMap      map[int]string
	Job             *models.Job
	Jobs            []*models.Job
	Driver          *models.Driver
	Drivers         []*models.Driver
	Truck           *models.Truck
	Trucks          []*models.Truck
	JobDisplays     []*JobDisplay
	Form            any
	Flash           string
	IsAuthenticated bool
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		// Parse the base template file into a template set.
		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}
		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Add the template set to the map as normal...
		cache[name] = ts
	}
	return cache, nil
}
