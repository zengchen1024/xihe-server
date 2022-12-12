package repository

type WuKong interface {
	ListSamples(string, []int) ([]string, error)
}
