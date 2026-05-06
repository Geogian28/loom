package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func setupRoutes(ac appConfig) {
	// 1. Define the FuncMap BEFORE parsing
	funcMap := template.FuncMap{
		"seq": func(start, end int) []int {
			if end < start {
				return []int{}
			}
			s := make([]int, end-start+1)
			for i := range s {
				s[i] = start + i
			}
			return s
		},
		"sub": func(a, b int) int { return a - b },
		"mod": func(a, b int) int { return a % b },
	}

	// 2. Properly chain: New -> Funcs -> ParseFiles
	// The name passed to New() must match the filename in ParseFiles
	var err error
	tmpl, err = template.New("dashboard.html").Funcs(funcMap).ParseFiles(ac.templatesDir + "/dashboard.html")
	// tmpl, err := template.ParseGlob(ac.templatesDir + "/*.html")
	if err != nil {
		log.Fatalf("Loom: Failed to parse template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Fetch live data
		data, err := getDeviceInfo(ac.address, ac.port, ac.username, ac.password)
		if err != nil {
			http.Error(w, "Loom: Failed to poll switch - "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Lookup physical layout
		layout, ok := ChassisTemplates[data.model]
		if !ok {
			layout = ChassisTemplates["DCS-7050S-64"] // Fallback
		}

		// Map ports for easy template lookup
		portMap := make(map[string]portInfo)
		for _, p := range data.port {
			portMap[p.name] = p
		}

		pageData := struct {
			Device deviceInfo
			Layout Layout
			Ports  map[string]portInfo
		}{
			Device: data,
			Layout: layout,
			Ports:  portMap,
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Printf("Loom: Render error: %v", err)
		}
	})
}
