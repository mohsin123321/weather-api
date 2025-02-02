document.getElementById("weather-form").addEventListener("submit", function(event) {
    event.preventDefault();
    let city = document.getElementById("city-input").value;

    // Use Axios to fetch data
    axios.get(`/api/weather?city=${city}`)
        .then(response => {
            let data = response.data;

            if (data.cod !== 200) {
                alert("Error: " + data.message);
                return;
            }

            document.getElementById("weather-result").style.display = "block";
            document.getElementById("city-name").innerText = `Weather in ${data.name}, ${data.sys.country}`;
            document.getElementById("coordinates").innerText = `Coordinates: ${data.coord.lat}, ${data.coord.lon}`;
            document.getElementById("weather-condition").innerText = `${data.weather[0].main} - ${data.weather[0].description}`;
            document.getElementById("temperature").innerText = `${data.main.temp}°C (Feels like ${data.main.feels_like}°C)`;
            document.getElementById("temperature-range").innerText = `Min: ${data.main.temp_min}°C | Max: ${data.main.temp_max}°C`;
            document.getElementById("wind").innerText = `${data.wind.speed} m/s (Gusts: ${data.wind.gust + " m/s" || "N/A"}, Direction: ${data.wind.deg}°)`;
            document.getElementById("clouds").innerText = `${data.clouds.all}%`;
            document.getElementById("visibility").innerText = `${data.visibility} meters`;
            document.getElementById("pressure").innerText = data.main.pressure;
            document.getElementById("humidity").innerText = data.main.humidity;

            let rainElem = document.getElementById("rain");
            let snowElem = document.getElementById("snow");
            if (data.rain && data.rain["1h"]) {
                document.getElementById("rain-value").innerText = data.rain["1h"];
                rainElem.style.display = "block";
            } else {
                rainElem.style.display = "none";
            }
            if (data.snow && data.snow["1h"]) {
                document.getElementById("snow-value").innerText = data.snow["1h"];
                snowElem.style.display = "block";
            } else {
                snowElem.style.display = "none";
            }

            let sunriseTime = new Date((data.sys.sunrise + data.timezone) * 1000).toLocaleTimeString();
            let sunsetTime = new Date((data.sys.sunset + data.timezone) * 1000).toLocaleTimeString();
            document.getElementById("sunrise-time").innerText = sunriseTime;
            document.getElementById("sunset-time").innerText = sunsetTime;
        })
        .catch(error => {
            console.error("Error fetching weather:", error);
            alert("Failed to fetch weather data.");
        });
});
