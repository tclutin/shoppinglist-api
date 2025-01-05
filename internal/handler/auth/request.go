package auth

type SignUpRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	Gender   string `json:"gender" binding:"required,oneof=MALE FEMALE NONE"`
}

type LogInRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}
