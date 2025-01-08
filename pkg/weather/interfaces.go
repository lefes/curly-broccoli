package weather

type WeatherApi interface {
	CurrentWeather(city string) (*WeatherResponse, error)
	ForecastWeather(city string, days int) (*WeatherResponse, error)
}
