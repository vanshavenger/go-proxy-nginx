package utils

// Header represents a key-value pair of HTTP headers
type Header struct {
	Key   string `yaml:"key" json:"key" validate:"required"`
	Value string `yaml:"value" json:"value" validate:"required"`
}

// Upstream represents an upstream server
type Upstream struct {
	ID       string `yaml:"id" json:"id" validate:"required"`
	URL      string `yaml:"url" json:"url" validate:"required"`
	Protocol string `yaml:"protocol" json:"protocol" validate:"omitempty,oneof=http https"`
}

// Rule represents a routing rule
type Rule struct {
	Path      string   `yaml:"path" json:"path" validate:"required"`
	Upstreams []string `yaml:"upstreams" json:"upstreams" validate:"required,min=1"`
}

// Server represents the server configuration
type Server struct {
	Listen    int        `yaml:"listen" json:"listen" validate:"required"`
	Workers   *int       `yaml:"workers" json:"workers" validate:"omitempty,min=1"`
	Upstreams []Upstream `yaml:"upstreams" json:"upstreams" validate:"required,min=1,dive"`
	Headers   []Header   `yaml:"headers" json:"headers" validate:"omitempty,dive"`
	Rules     []Rule     `yaml:"rules" json:"rules" validate:"required,min=1,dive"`
}

// RootConfig represents the root configuration
type RootConfig struct {
	Server Server `yaml:"server" json:"server" validate:"required"`
}
