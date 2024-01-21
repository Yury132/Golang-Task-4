package models

type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Nation     string `json:"nation"`
}

// Структура получаемого возраста от внешнего api
type AgeApi struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

// Структура получаемого пола от внешнего api
type GenderApi struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

// Структура получаемой национальности от внешнего api
type NationApi struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []Country `json:"country"`
}

// Вспомогательная структура - страна с вероятностью
type Country struct {
	Country_id  string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
