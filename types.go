package main

type Dream struct {
	Id string //url name
	Name string
	Image string
	ImageHqLink string
	ImageViewLink string
	Center bool
}

type Connection struct {
	FromID string
	ToID string
}