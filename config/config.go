package config

type ServerProperties struct {
	Bind      string `cfg:"bind"`
	Port      int    `cfg:"port"`
	Databases int    `cfg:"databases"`
}

var Properties *ServerProperties
