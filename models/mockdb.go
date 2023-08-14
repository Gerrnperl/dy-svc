package models

type MockDB struct{}

func (m *MockDB) GetUserById(id int64) (*User, error) {
	// 返回一个虚假的用户对象
	return &User{
		Id:   id,
		Name: "test",
	}, nil
}
