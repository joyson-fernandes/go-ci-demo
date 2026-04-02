package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type WeatherResponse struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity"`
	Wind        float64 `json:"wind"`
	Timestamp   string  `json:"timestamp"`
}

var cities = map[string]WeatherResponse{
	"london":    {Location: "London, UK", Temperature: 14.2, Condition: "Cloudy", Humidity: 72, Wind: 12.5},
	"newyork":   {Location: "New York, US", Temperature: 18.5, Condition: "Sunny", Humidity: 55, Wind: 8.3},
	"tokyo":     {Location: "Tokyo, JP", Temperature: 22.1, Condition: "Partly Cloudy", Humidity: 65, Wind: 6.1},
	"dubai":     {Location: "Dubai, AE", Temperature: 35.8, Condition: "Clear", Humidity: 40, Wind: 15.2},
	"sydney":    {Location: "Sydney, AU", Temperature: 20.3, Condition: "Rainy", Humidity: 80, Wind: 18.7},
	"mumbai":    {Location: "Mumbai, IN", Temperature: 32.4, Condition: "Humid", Humidity: 85, Wind: 9.4},
	"paris":     {Location: "Paris, FR", Temperature: 16.7, Condition: "Overcast", Humidity: 68, Wind: 11.0},
	"singapore": {Location: "Singapore, SG", Temperature: 30.1, Condition: "Thunderstorm", Humidity: 90, Wind: 5.5},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/weather/", weatherHandler)
	http.HandleFunc("/api/cities", citiesHandler)
	http.HandleFunc("/healthz", healthHandler)

	log.Printf("Weather API starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

const htmlPage = `<!DOCTYPE html>
<html>
<head>
    <title>Weather App</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, sans-serif; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 100vh; display: flex; flex-direction: column; align-items: center; padding: 40px 20px; color: white; }
        h1 { font-size: 2.5em; margin-bottom: 10px; }
        p.sub { opacity: 0.8; margin-bottom: 30px; }
        .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 20px; max-width: 1200px; width: 100%; }
        .card { background: rgba(255,255,255,0.15); backdrop-filter: blur(10px); border-radius: 16px; padding: 24px; transition: transform 0.2s; }
        .card:hover { transform: translateY(-4px); }
        .city { font-size: 1.3em; font-weight: 600; }
        .temp { font-size: 3em; font-weight: 700; margin: 10px 0; }
        .details { display: flex; justify-content: space-between; opacity: 0.8; font-size: 0.9em; margin-top: 12px; }
        .condition { background: rgba(255,255,255,0.2); display: inline-block; padding: 4px 12px; border-radius: 20px; font-size: 0.85em; margin-top: 8px; }
        footer { margin-top: 40px; opacity: 0.6; font-size: 0.85em; }
    </style>
</head>
<body>
    <h1>Weather App</h1>
    <p class="sub">Deployed on Kubernetes via GitHub Actions</p>
    <div class="grid" id="cards"></div>
    <footer>Built by Joyson Fernandes | Go + GitHub Actions + K8s</footer>
    <script>
        fetch('/api/cities').then(function(r){return r.json()}).then(function(cities){
            cities.forEach(function(c){
                fetch('/api/weather/'+c).then(function(r){return r.json()}).then(function(w){
                    var card = '<div class="card">';
                    card += '<div class="city">'+w.location+'</div>';
                    card += '<div class="temp">'+w.temperature+'\u00B0C</div>';
                    card += '<span class="condition">'+w.condition+'</span>';
                    card += '<div class="details">';
                    card += '<span>Humidity: '+w.humidity+'%</span>';
                    card += '<span>Wind: '+w.wind+' km/h</span>';
                    card += '</div></div>';
                    document.getElementById('cards').innerHTML += card;
                });
            });
        });
    </script>
</body>
</html>`

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlPage)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Path[len("/api/weather/"):]
	weather, ok := cities[city]
	if !ok {
		http.Error(w, `{"error":"city not found"}`, http.StatusNotFound)
		return
	}
	weather.Timestamp = time.Now().UTC().Format(time.RFC3339)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func citiesHandler(w http.ResponseWriter, r *http.Request) {
	names := make([]string, 0, len(cities))
	for k := range cities {
		names = append(names, k)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
