package economy

type Role struct {
	ID              int
	Name            string
	RespectRequired int
	Privileges      map[string]interface{}
}
