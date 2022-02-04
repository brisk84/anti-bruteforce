package app

type LoginInfo struct {
	Login    string
	Password string
	Ip       string
}

type NetworkInfo struct {
	Ip string
}

type Error struct {
	Code int32
	Info string
}
