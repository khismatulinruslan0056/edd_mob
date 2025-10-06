package model

type User struct {
	ID          int     `json:"id,omitempty"`
	Name        string  `json:"name" example:"Dmitriy"`
	Surname     string  `json:"surname" example:"Ivanov"`
	Patronymic  *string `json:"patronymic,omitempty" example:"Sergeevich"`
	Age         *int    `json:"age,omitempty" example:"30"`
	Gender      *string `json:"gender,omitempty" example:"male"`
	Nationality *string `json:"nationality,omitempty" example:"RU"`
}

type UserAge struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type UserGender struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type UserNationality struct {
	Count     int       `json:"count"`
	Name      string    `json:"name"`
	Countries []Country `json:"country"`
}

type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
