package database

type Catalog struct {
	Id       int    `json:"id"`
	RegNum   string `json:"regNum"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Year     string `json:"year"`
	FullName string `json:"FullName"`
}

type PageCatalog struct {
	LastInd int       `json:"lastInd"`
	Catalog []Catalog `json:"catalog"`
}

type UpdateRegNum struct {
	RegNum string `json:"regNum"`
}

type UpdateBrandModel struct {
	Brand string `json:"brand"`
	Model string `json:"model"`
}

type UpdateYear struct {
	Year string `json:"year"`
}

type UpdateHolder struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type ItemsCatalog struct {
	RegNum string `json:"regNum"`
	Brand  string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year,omitempty"`
	Owner  People `json:"owner"`
}

type Brand struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Model struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
