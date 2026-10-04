package manifest

const (
	SchemaVersion = 1
	DefaultURL    = "https://raw.githubusercontent.com/yuhaiin/kitte/auto-update/manifest.json"
)

type Manifest struct {
	SchemaVersion int    `json:"schemaVersion"`
	Repository    string `json:"repository"`
	Branch        string `json:"branch"`
	BaseURL       string `json:"baseUrl"`
	Files         []File `json:"files"`
}

type File struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Kind       string `json:"kind"`
	Usage      string `json:"usage"`
	ListType   string `json:"listType,omitempty"`
	Format     string `json:"format"`
	Path       string `json:"path"`
	URL        string `json:"url"`
	SourceURL  string `json:"sourceUrl,omitempty"`
	Size       int64  `json:"size"`
	SHA256     string `json:"sha256"`
	Selectable bool   `json:"selectable"`
}
