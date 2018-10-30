package traefik

type Backend struct {
	Name   string
	Node   string
	URL    string
	Weight uint
}
