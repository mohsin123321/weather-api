# Weather App
A simple weather application that fetches real-time weather data based on the user's selected city using the OpenWeatherMap API. The weather details are displayed in a user-friendly, well-structured card format.

## Technologies
- **Golang 1.22.2**
- **OpenMapWeatherApi**
- **HTML, CSS, and JavaScript**

## Configuration
This application requires certain environment variables to be set. These values can be provided via a `.env` file that should be present at the root directory or set manually in the system environment, otherwise default value will be assigned to them except `WEATHER_API_KEY`.

### **Environment Variables**
| Variable Name               | Description                                | Default Value |
|-----------------------------|--------------------------------------------|--------------|
| `PORT`                      | The port on which the server runs          | `8080`       |
| `WEATHER_API_KEY`           | OpenWeatherMap API key (required)          | _empty_      |
| `CACHE_EXPIRY`              | Cache expiry time in minutes               | `5`          |
| `OPEN_WEATHER_MAP_BASE_URL` | Base URL for OpenWeatherMap API requests   | `https://api.openweathermap.org/data/2.5` |

## Start the Server
To start the server, run the following command:

```bash
go run ./cmd/api/main.go