package service

import "regexp"

func checkUsername(username string) (bool, error) {
	r, err := regexp.Compile("^[a-zA-Z0-9_]{3,30}$")
	if err != nil {
		return false, err
	}

	return r.MatchString(username), nil
}

func checkPassword(password string) (bool, error) {
	r, err := regexp.Compile("^[A-Za-z0-9!@#$%^&*]{8,64}$")
	if err != nil {
		return false, err
	}

	return r.MatchString(password), nil
}
