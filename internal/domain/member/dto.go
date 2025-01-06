package member

type MemberDTO struct {
	MemberID uint64 `db:"member_id"`
	Username string `db:"username"`
	Gender   string `db:"gender"`
	Role     string `db:"role"`
}
