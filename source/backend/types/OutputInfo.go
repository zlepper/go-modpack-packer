package types

type OutputInfo struct {
	File             string          `json:"file"`
	Name             string          `json:"name"`
	// Slug
	Id               string          `json:"id"`
	Version          string          `json:"version"`
	MinecraftVersion string          `json:"minecraftVersion"`
	Description      string          `json:"description"`
	Author           string          `json:"author"`
	Url              string          `json:"url"`
	ProgressKey      string          `json:"progressKey"`
	IsOnSolder       bool            `json:"isOnSolder"`
	Permissions      *UserPermission `json:"permissions"`
}

func (o *OutputInfo) GenerateOnlineVersion() string {
	return o.MinecraftVersion + "-" + o.Version
}
