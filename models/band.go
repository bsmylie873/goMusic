package models

type Band struct {
	Id              int    `json:"id"`
	Name            string `json:"name" validate:"required,min=1,max=100"`
	Nationality     string `json:"nationality" validate:"required,min=1,max=100"`
	NumberOfMembers int    `json:"number_of_members" validate:"required,min=1,max=500"`
	DateFormed      string `json:"date_formed" validate:"required,min=1,max=100"`
	Age             int    `json:"age" validate:"required,min=0,max=150"`
	Active          bool   `json:"active"`
}
