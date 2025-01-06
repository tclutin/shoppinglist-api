package group

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required,max=255"`
}

type JoinToGroupRequest struct {
	Code string `json:"code" binding:"required"`
}
