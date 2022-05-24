package models

type Response struct {
	Status  string
	Message interface{}
}

type ErrorMessage struct {
	Code  int
	Error string
}
