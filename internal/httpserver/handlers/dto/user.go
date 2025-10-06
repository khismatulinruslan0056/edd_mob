package dto

type UserRequest struct {
	Name       string  `json:"name" example:"Dmitriy"`
	Surname    string  `json:"surname" example:"Ivanov"`
	Patronymic *string `json:"patronymic,omitempty" example:"Sergeevich"`
}

type UserResponse struct {
	ID          int     `json:"id" example:"1"`
	Name        string  `json:"name" example:"Dmitriy"`
	Surname     string  `json:"surname" example:"Ivanov"`
	Patronymic  *string `json:"patronymic,omitempty" example:"Sergeevich"`
	Age         *int    `json:"age,omitempty" example:"30"`
	Gender      *string `json:"gender,omitempty" example:"male"`
	Nationality *string `json:"nationality,omitempty" example:"RU"`
}
