package user

type User struct {
	Id uint64
	UserName string
	// other stats
}

func Load() (User, error) {
	// TODO:
	// Look into data file for credentials and send a get request to the server.
	// If the credentials are wrong or there are not correct, return an empty user.
	// Otherwise parse the response and load the data accordingly
	return User{}, nil
}
