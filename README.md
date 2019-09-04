# Web-crawler

What this app does
------------------

Application written in Go which navigates through the links of a given domain and creates a file with the domain sitemap. It excludes external links.

I.e. 

For the domain http://www.mydomain.com the application parses the initial page and navigates through all the found links to generate a list of all the links belonging to the given domain. It will exclude external links like www.facebook.com or foo.mydomain.com

## Instalation

###Requirements

- This application has been built and tested with Golang 1.12

Clone the repository and build the application using go. This applications uses Go Moduiles and therefore it can be checked out to any location.

## How to run it

Once the applicaiton has been built execute the binary with the minimum required arguments (-url)

`./web-crawler -url https://www.domain.com`

This will create a new file named `sitemap.txt` in the current locaiton. 

Optionally you can specify the path where the output file must be created:

`./web-crawler -url https://www.domain.com -output-file path_to_sitemap`

## Using Docker

###Requirements

- Docker
  
#### Build the Docker image

This Dockerfile builds and executes all the tests of the application. Then it copies the generated binary to a new image with the minimum required elements.

Build the image with Docker, i.e.

`docker build -t web-crawler .`

#### Run Webcrawler with Docker 

For execution with Docker a volume needs to be created with the desired location where the output will be created. `-output-file` parameter is required and must match the mapped volume, i.e.

`docker run -v $(pwd):/data web-crawler -url https://www.domain.com -output-file /data/my-sitemap.txt`

