// Package dto provides ...
package dto

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Readme   bool   `json:"readme"`
}
