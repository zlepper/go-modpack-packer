package types

type OutputInfo struct {
	File             string
	Name             string
	Id               string
	Version          string
	MinecraftVersion string
	Description      string
	Author           string
	Url              string
	ProgressKey      string
}

func (o *OutputInfo) GenerateOnlineVersion() string {
	return o.MinecraftVersion + "-" + o.Version
}
