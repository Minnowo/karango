
tailwindcss  -i ./assets//tailwind.css -o ./assets/static/c/main.css
templ generate

set LOG_LEVEL=debug

go run main.go
