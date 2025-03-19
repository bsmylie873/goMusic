package models

type Band struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Nationality     string `json:"nationality"`
	NumberOfMembers int    `json:"number_of_members"`
	DateFormed      string `json:"date_formed"`
	Age             int    `json:"age"`
	Active          bool   `json:"active"`
}
