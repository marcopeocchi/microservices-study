package pkg

type DirectortyList = string

// Generic response for directory listing operations
type Response[T any] struct {
	List []T `json:"list"`
}

// Struct used for serdes operation on login requests
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
