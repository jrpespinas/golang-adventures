// Package weather provides a simple function to forcast the weather.
package weather

// CurrentCondition is a package level variable which represents the weather condition.
var CurrentCondition string

// CurrentLocation is a package level variable which represents the city or location to be forcasted.
var CurrentLocation string

// Forecast is a function which returns weather condition in a given city.
func Forecast(city, condition string) string {
	CurrentLocation, CurrentCondition = city, condition
	return CurrentLocation + " - current weather condition: " + CurrentCondition
}
