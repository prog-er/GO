// Example unit test for server logic
func TestUserService_AddUser(t *testing.T) {
	server := &userServiceServer{}
	user := &user.User{Name: "Test User", Email: "test@example.com"}

	resp, err := server.AddUser(context.Background(), user)
	if err != nil {
		t.Errorf("AddUser failed: %v", err)
	}

	if resp.Id == 0 {
		t.Errorf("AddUser returned an invalid user ID")
	}
}
