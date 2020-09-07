package entity

type Blog struct {
	isDir bool
	name string

	children []Blog
}
