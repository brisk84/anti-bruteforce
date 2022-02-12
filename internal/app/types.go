package app

type LoginInfo struct {
	Login    string
	Password string
	IP       string
}

type NetworkInfo struct {
	IP string
}

type Error struct {
	Code int32
	Info string
}
