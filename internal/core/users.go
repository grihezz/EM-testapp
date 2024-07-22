package core

type User struct {
	ID             int    `json:"id" db:"id"`
	PassportNumber string `json:"passport_number" db:"passport_number"`
	Surname        string `json:"surname" db:"surname"`
	Name           string `json:"name" db:"name"`
	Patronymic     string `json:"patronymic" db:"patronymic"`
	Address        string `json:"address" db:"address"`
}
type Users []User

type ServiceUser struct {
	PassportNum string `json:"passportHash"`
	Surname     string `json:"surname"`
	Name        string `json:"name"`
	Patronymic  string `json:"patronymic"`
	Address     string `json:"address"`
}

type UserFilter struct {
	PassportNum string
	Surname     string
	Name        string
	Patronymic  string
	Address     string
}
