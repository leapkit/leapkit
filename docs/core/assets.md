---
index: 4
title: "Assets"
---

The assets package provides a way to manage web assets in your web application. It receives a folder where your assets are on the disk and returns an http.Handler that serves these files.

This handler is capable of a few things that are useful for web development:
- Asset fingerprinting (to avoid caching issues)
- Hot code reloading (when the assets change on disk)

## Usage

To set our assets into our Leapkit app, we need to use the `WithAssets` server option. It receives an FS parameter that will help locate template files to be served.

```go

//go:embed templates/**/*.html
var templatesFS embed.FS

s := server.New(
	server.WithAssets(templatesFS),
)
```

## Fingerprinting Helper
The `server.WithAssets` option also setus the `assetsPath` helper that can be used in your templates to use the fingerprinted version of an asset.

```html
<link rel="stylesheet" href="<%= assetPath(`/css/app.css`) %>">
// will output something like
<link rel="stylesheet" href="/css/app-cafe123ff22112eedd.css">
```

## Hotcode Reloading
The assets managers provides a handler function capable of serving the files in the assets filesystem. This handler considers the `GO_ENV` variable to look in for files in disk before looking into the embedded filesystem passed.
