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
