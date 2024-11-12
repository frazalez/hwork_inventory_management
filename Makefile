build:
	env GOOS=windows GOARCH=386 go build -o srv.exe src/main.go
	env GOOS=linux GOARCH=386 go build -o srv src/main.go

run:
	./templ generate
	go run src/main.go

tailwind:
	./tailwindcss-linux-x64 -i static/input.css -o static/tailwind.css -w
