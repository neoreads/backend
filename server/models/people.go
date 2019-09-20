package models

type Person struct {
	ID         string `json:"id,omitempty"`
	FullName   string `json:"fullname,omitempty"`
	Intro      string `json:"intro,omitempty"`
	OtherNames string `json:"othernames,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
}
