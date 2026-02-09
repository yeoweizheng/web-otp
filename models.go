package main

type User struct {
	Id       int64
	Username string
	Password string
}

type Account struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type AccountOTP struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	OTP  string `json:"otp"`
}
