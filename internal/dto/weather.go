package dto

// OpenWeatherMapResponse represents the full response from the OpenWeatherMap API
type OpenWeatherMapResponse struct {
	ID         uint      `json:"id"`
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility uint      `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         uint64    `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
	Rain       *Rain     `json:"rain,omitempty"`
	Snow       *Snow     `json:"snow,omitempty"`
}

// Coord represents the coordinates
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// Weather represents a weather condition
type Weather struct {
	ID          uint   `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Main contains temperature and atmospheric details
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  int     `json:"pressure"`
	Humidity  uint    `json:"humidity"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

// Wind contains wind speed and direction
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

// Clouds contains cloud percentage
type Clouds struct {
	All uint `json:"all"`
}

// Sys contains system-related information
type Sys struct {
	Type    int     `json:"type"`
	ID      uint    `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise uint64  `json:"sunrise"`
	Sunset  uint64  `json:"sunset"`
}

// Rain contains rain volume for the last hour
type Rain struct {
	OneHour float64 `json:"1h,omitempty"`
}

// Snow contains snow volume for the last hour
type Snow struct {
	OneHour float64 `json:"1h,omitempty"`
}
