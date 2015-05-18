default: image

image: gox
	docker build -t bazooka/parser-php .

gox:
	gox -osarch="linux/amd64" -output="main"
