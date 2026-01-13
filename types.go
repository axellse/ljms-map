package main

type Dream struct {
	Id string //url name
	ImageCacheType string
	Name string
	Image string
	ImageHqLink string
	ImageHqLocalLink string
	Center bool
}

type Connection struct {
	FromID string
	ToID string
}