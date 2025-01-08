package group

type CreateGroupDTO struct {
	OwnerID     uint64
	Name        string
	Description string
}

type JoinToGroupDTO struct {
	UserID uint64
	Code   string
}

type GroupUserDTO struct {
	GroupID uint64
	UserID  uint64
}

type KickMemberDTO struct {
	GroupID  uint64
	UserID   uint64
	MemberID uint64
}

type GroupDTO struct {
	GroupID     uint64 `json:"groupID" db:"group_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Code        string `json:"code" db:"code"`
}

type CreateProductDTO struct {
	UserID        uint64
	GroupID       uint64
	ProductNameID uint64
	Quantity      int
	AddedBy       uint64
}

type RemoveProductDTO struct {
	ProductID uint64
	GroupID   uint64
	UserID    uint64
}
