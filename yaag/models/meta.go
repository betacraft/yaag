package models

//go:generate easytags $GOFILE

type MetaData struct {
	Path        string
	Description string
	Custom      map[string]string
}
