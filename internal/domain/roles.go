package domain

var RolesIDs = map[string]int{
	"Писька":       1,
	"Крыса":        10,
	"Голова":       20,
	"Нормис":       30,
	"Казах":        40,
	"Дед":          50,
	"Ебанатор3000": 1000,
}

type Role struct {
	ID              int
	Name            string
	RespectRequired int
	Privileges      map[string]interface{}
}
