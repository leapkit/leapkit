## LeapKit Template

<img width="300" alt="logo" src="https://github.com/leapkit/leapkit/template/assets/645522/d5bcb8ed-c763-4b39-8cfb-aed694b87646">
<br><br>

This is the  LeapKit template for building web applications with Go, HTMX and Tailwind CSS. It integrates useful features such as hot code reload and css recompiling.

### Getting started

Install dependencies:

```sh
go mod download
go run ./cmd/setup
```

### Building the application

```sh
# Building TailwindCSS with tailo
> go run github.com/paganotoni/tailo/cmd/build@v1.0.8

# Building the app
> go build -tags osusergo,netgo -buildvcs=false -o bin/app ./cmd/app
```

### Running the application

To run the application in development mode execute:

```sh
> kit dev
```

And open `http://localhost:3000` in your browser.

### Generating a migration

```sh
> kit g migration [name]
```
