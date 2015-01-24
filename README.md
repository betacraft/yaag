# YAAG : Yet Another API doc Generator

Golang is awesome for developing web apps. And people have created a bunch of awesome Web-Frameworks, Web hepler libraries. If we consider the entire ecosystem of web apps in Golang everything except API documentation seems to be in place. So we have created the first API doc generator for Golang based web apps and calling it Yet another.

## Why ?

Most of the web services expose their APIs to the mobile or third party developers. Documenting them is somewhat pain in the ass. We are trying to reduce the pain, atleast for in house projects where you don't have to expose your documentation to the world. YAAG generates simple bootstrap based API documentation without writing a single line of comments.

## How it works ?

YAAG is a middleware. You have to add YAAG handler in your routes and you are done. Just go on calling your APIs using POSTMAN, Curl or from any client, and YAAG will keep on updating the API Doc html. 


### Config parameters 

ReadMode :  If true then YAAG middleware will function and start recording the API calls and updating the API doc
DocPath  :  Path where the API doc will be saved
DocTitle :  API Doc title
BaseUrl  :  Base URL of the Endpoints

## Support

It's a middleware supporting http.Handler interface. So it will support all the frameworks that support Handler like martini, Gorilla Mux etc. 
YAAG also supports revel framework.
YAAG also supports gin framework.

## How to use with Revel

1. Add yaag.record = true in conf/app.conf (before starting to record the api calls)
2. import github.com/gophergala/yaag/filters in app/init.go
3. add 'filters.FilterForApiDoc' in revel.Filters
4. Start recording Api calls


## Team

Aniket Awati (aniket@rainingclouds.com)
Akshay Deo (akshay@rainingclouds.com)
Kaustubh Deshmukh (kaustubh@rainingclouds.com)

This project is initiated by RainingClouds Inc during GopherGala 2015.
