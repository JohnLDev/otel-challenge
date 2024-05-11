package config

type config struct {
	ServiceName string
}

var Conf = config{
	ServiceName: "GetCepApp",
}
