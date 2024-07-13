package models

type Task struct {
	ID      string `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Date    string `json:"date,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}
