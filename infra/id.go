package infra

import "strconv"

type IdGenerator struct {
	id int
}

// Only use with mutex in controller
func (g *IdGenerator) Generate() string {
	result := g.id
	g.id++
	return strconv.FormatInt(int64(result), 10)
}
