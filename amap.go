package puppy

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
)

const weather_url = "https://restapi.amap.com/v3/weather/weatherInfo"

type WeatherResult struct {
	Date         string
	DayTemp      string
	NightTemp    string
	DayWeather   string
	NightWeather string
}

type AmapWeather struct {
	key    string
	client *req.Client
}

func NewAmapWeather(key string, client *req.Client) *AmapWeather {
	return &AmapWeather{
		key:    key,
		client: client,
	}
}

func (w *AmapWeather) GetWeather(ctx context.Context, code string) (*WeatherResult, error) {
	r := w.client.R().SetContext(ctx).
		AddQueryParam("key", w.key).
		AddQueryParam("city", code).
		AddQueryParam("extensions", "all")

	resp, err := r.Get(weather_url)
	if err != nil {
		return nil, err
	}

	x := gjson.GetBytes(resp.Bytes(), "forecasts.0.casts.0")
	return &WeatherResult{
		Date:         x.Get("date").String(),
		DayTemp:      x.Get("daytemp").String(),
		NightTemp:    x.Get("nighttemp").String(),
		DayWeather:   x.Get("dayweather").String(),
		NightWeather: x.Get("nightweather").String(),
	}, nil
}
