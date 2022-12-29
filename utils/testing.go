package utils

type Case struct {
	Check  bool // case fails if check is true
	Format string
	Args   []any
}
