package pkg

type Preview struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Loved     bool   `json:"loved"`
	Thumbnail string `json:"thumbnail"`
}

type DirectortyList = string

// Generic response for directory listing operations
type Response[T Preview | DirectortyList] struct {
	List []T `json:"list"`
}

type PaginatedResponse[T Preview | DirectortyList] struct {
	List          []T   `json:"list"`
	Page          int   `json:"page"`
	Pages         int   `json:"pages"`
	PageSize      int   `json:"pageSize"`
	TotalElements int64 `json:"totalElements"`
}

type Config struct {
	ServerSecret      string `yaml:"serverSecret"`
	Masterpass        string `yaml:"masterpass"`
	CacheDir          string `yaml:"cacheDir"`
	WorkingDir        string `yaml:"workingDir"`
	ForceRegeneration bool   `yaml:"regenerateThumbnailsOnBoot"`
	ThumbnailHeight   int    `yaml:"thumbnailHeight"`
	ThumbnailQuality  int    `yaml:"thumbnailQuality"`
	Port              int    `yaml:"port"`
}

// Struct used for serdes operation on login requests
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
