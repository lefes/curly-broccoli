package domain

type WeatherResponse struct {
	ResolvedAddress   string        `json:"resolvedAddress"`
	CurrentConditions CurrentData   `json:"currentConditions"`
	Days              []ForecastDay `json:"days"`
}

type CurrentData struct {
	Temp      float64 `json:"temp"`
	Condition string  `json:"conditions"`
}

type ForecastDay struct {
	Datetime  string  `json:"datetime"`
	TempMax   float64 `json:"tempmax"`
	TempMin   float64 `json:"tempmin"`
	Condition string  `json:"conditions"`
	SunRise   string  `json:"sunrise"`
	SunSet    string  `json:"sunset"`
}

// CURRENT
// "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Sofia/today?unitGroup=metric&include=current&key=TOKEN&contentType=json&lang=ru"

// FORECAST
// "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Sofia/2025-01-05/2025-01-07?unitGroup=metric&include=days&key=TOKEN&contentType=json&lang=ru"
