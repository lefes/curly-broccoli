package domain

var Roles = []Role{
	{
		ID:              1,
		Name:            "Писька",
		RespectRequired: 0,
	},
	{
		ID:              10,
		Name:            "Крыса",
		RespectRequired: 100,
	},
	{
		ID:              20,
		Name:            "Голова",
		RespectRequired: 300,
	},
	{
		ID:              30,
		Name:            "Нормис",
		RespectRequired: 600,
	},
	{
		ID:              40,
		Name:            "Казах",
		RespectRequired: 1200,
	},
	{
		ID:              50,
		Name:            "Дед",
		RespectRequired: 2500,
	},
	{
		ID:              1000,
		Name:            "Ебанатор3000",
		RespectRequired: -1,
	},
}

type Role struct {
	ID              int
	Name            string
	RespectRequired int
	Privileges      map[string]interface{}
}
