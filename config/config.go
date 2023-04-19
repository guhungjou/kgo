package config

func Init(_prefix string) error {
	return ViperInit(_prefix)
}

func Unmarshal(k string, v interface{}) error {
	return ViperUnmarshal(k, v)
}
