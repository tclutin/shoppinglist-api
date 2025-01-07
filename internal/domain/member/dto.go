package member

type MemberDTO struct {
	MemberID uint64 `json:"member_id"  db:"member_id"`
	Username string `json:"username" db:"username"`
	Gender   string `json:"gender" db:"gender"`
	Role     string `json:"role" db:"role"`
}
