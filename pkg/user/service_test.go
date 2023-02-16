package user

// func TestLogin(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	db, err := gorm.Open(sqlite.Open(filepath.Join(".", ".cache", "fuu.db")), &gorm.Config{})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	repo := Repository{db: db}
// 	service := Service{repo: repo}

// 	token, err := service.Login(ctx, "username_not", "password")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if len(*token) > 0 && len(*token) <= 2048 {
// 		fmt.Println(token)
// 	}
// }
