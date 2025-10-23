package user

type User struct {
	ID    uint   `db:"id"`
	Name  string `db:"name"`
	Phone string `db:"phone"`
	Role  string `db:"role"`
}
