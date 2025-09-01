package local

type MediaSource struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Extension string `json:"ext"`
	Size      int64  `json:"size"`
}

type Media map[string]MediaSource

type IconSource struct {
	Name      string
	Path      string
	Extension string
	Size      int64
}

type Subtitle struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Ext  string `json:"ext"`
}

type Subtitles map[string]Subtitle

type Group struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

type Collection struct {
	Name  string  `json:"name"`
	Group []Group `json:"group"`
}
