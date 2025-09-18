// --- INITIALIZE MAP ---
const map = L.map('map').setView([49.4431, 1.0993], 13); // Rouen coords
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '&copy; OpenStreetMap contributors'
}).addTo(map);

// --- WEATHER ---
async function fetchWeather(city = "Rouen") {
  try {
    const res = await fetch(`/api/weather?city=${encodeURIComponent(city)}`);
    if (!res.ok) throw new Error(`Server error: ${res.status}`);
    const data = await res.json();

    if (data.error) {
      document.getElementById("weatherInfo").innerText = "Weather error: " + data.error;
      return;
    }

    // Extract fields with defaults
    const cityName   = data.city || "Unknown";
    const temperature = data.temperature ?? "N/A";
    const feelsLike  = data.feels_like ?? "N/A";
    const tempMin    = data.temp_min ?? "N/A";
    const tempMax    = data.temp_max ?? "N/A";
    const pressure   = data.pressure ?? "N/A";
    const humidity   = data.humidity ?? "N/A";
    const conditions = (data.conditions || []).join(", ") || "N/A";
    const windSpeed  = data.wind_speed ?? "N/A";
    const windDeg    = data.wind_deg ?? "N/A";
    const visibility = data.visibility ?? "N/A";
    const rain1h     = data.rain_1h ?? 0;
    const clouds     = data.clouds ?? "N/A";
    const sunrise    = data.sunrise ? new Date(data.sunrise * 1000).toLocaleTimeString() : "N/A";
    const sunset     = data.sunset ? new Date(data.sunset * 1000).toLocaleTimeString() : "N/A";

    // Update DOM
    document.getElementById("weatherInfo").innerHTML = `
      <strong>City:</strong> ${cityName}<br>
      <strong>Temperature:</strong> ${temperature}°C (Feels like ${feelsLike}°C)<br>
      <strong>Min/Max:</strong> ${tempMin}°C / ${tempMax}°C<br>
      <strong>Conditions:</strong> ${conditions}<br>
      <strong>Humidity:</strong> ${humidity}%<br>
      <strong>Pressure:</strong> ${pressure} hPa<br>
      <strong>Wind:</strong> ${windSpeed} m/s at ${windDeg}°<br>
      <strong>Visibility:</strong> ${visibility} m<br>
      <strong>Rain (last 1h):</strong> ${rain1h} mm<br>
      <strong>Cloud cover:</strong> ${clouds}%<br>
      <strong>Sunrise:</strong> ${sunrise}<br>
      <strong>Sunset:</strong> ${sunset}
    `;
  } catch (err) {
    console.error("Weather fetch error:", err);
    document.getElementById("weatherInfo").innerText = "Failed to load weather.";
  }
}

// Auto-refresh every 5 minutes
setInterval(() => fetchWeather(currentCity), 300000);

// --- CRIME ---
async function fetchCrime() {
  try {
    const res = await fetch('/api/crime');
    const data = await res.json();
    const ul = document.getElementById('crime');
    ul.innerHTML = '';
    data.forEach(item => {
      const li = document.createElement('li');
      li.innerText = `${item.type} at ${item.location} (Severity: ${item.severity})`;
      ul.appendChild(li);

      L.marker([49.44 + Math.random()*0.01, 1.09 + Math.random()*0.01])
        .addTo(map)
        .bindPopup(`${item.type} at ${item.location}`);
    });
  } catch (err) {
    console.error("Crime fetch error:", err);
  }
}

// --- TRANSPORT ---
let transportMarkers = [];
let lastSearchedLabel = ""; // persist the last searched bus label

async function fetchTransport(label = "") {
  if (!label) return; // DO NOTHING if no label

  try {
    const params = new URLSearchParams();
    params.append("label", label);
    const res = await fetch("/api/transport?" + params.toString());
    const data = await res.json();

    const ul = document.getElementById("transport");
    ul.innerHTML = "";

    // Update input box to current label
    document.getElementById("labelInput").value = lastSearchedLabel;

    // Remove old markers
    transportMarkers.forEach(m => map.removeLayer(m));
    transportMarkers = [];

    data.forEach(v => {
      const li = document.createElement("li");
      li.innerHTML = `<strong>Vehicle ${v.label}</strong> (Route ${v.route_id}, Direction ${v.direction_id}) - Status: ${v.current_status}, Occupancy: ${v.occupancy}`;

      // Next stops + ETA if available
      if (v.next_stops && v.next_stops.length > 0) {
        const stops = v.next_stops.map(s => `${s.stop_id} at ${s.eta}`).join(", ");
        li.innerHTML += `<br>Next stops: ${stops}`;
      }

      ul.appendChild(li);

      // Map marker
      const marker = L.marker([v.lat, v.lon], { rotationAngle: v.bearing })
        .addTo(map)
        .bindPopup(`Vehicle ${v.label} - Route ${v.route_id}`);
      transportMarkers.push(marker);
    });

  } catch(err) {
    console.error("Transport fetch error:", err);
    document.getElementById("transport").innerText = "Failed to load transport data.";
  }
}

// Search button
document.getElementById("searchTransport").addEventListener("click", () => {
  lastSearchedLabel = document.getElementById("labelInput").value.trim();
  fetchTransport(lastSearchedLabel);
});

// Auto-refresh every 1 min using last searched label, only if label is set
setInterval(() => {
  if (lastSearchedLabel) fetchTransport(lastSearchedLabel);
}, 100000);


// --- EVENTS ---
async function fetchEvents() {
  try {
    const res = await fetch('/api/events');
    const data = await res.json();
    const ul = document.getElementById('events');
    ul.innerHTML = '';
    data.forEach(item => {
      const li = document.createElement('li');
      li.innerText = `${item.name} at ${item.location} on ${item.date}`;
      ul.appendChild(li);

      L.marker([49.44 + Math.random()*0.01, 1.09 + Math.random()*0.01])
        .addTo(map)
        .bindPopup(`${item.name} at ${item.location}`);
    });
  } catch (err) {
    console.error("Events fetch error:", err);
  }
}

// --- INITIAL LOAD ---
const cityInput = document.getElementById("cityInput");
let currentCity = cityInput?.value || "Rouen";

fetchWeather(currentCity);
fetchCrime();
fetchTransport();
fetchEvents();

// --- MANUAL REFRESH WEATHER ---
const refreshBtn = document.getElementById("refreshWeather");
refreshBtn?.addEventListener("click", () => {
  currentCity = cityInput?.value || "Rouen";
  fetchWeather(currentCity);
});

// --- AUTO-REFRESH WEATHER EVERY 5 MINUTES ---
setInterval(() => fetchWeather(currentCity), 300000);

// --- AUTO-REFRESH TRANSPORT EVERY 30 SECONDS ---
setInterval(fetchTransport, 30000);
