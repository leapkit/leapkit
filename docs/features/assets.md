---
index: 2
title: "Assets"
---

The assets package provides a way to manage web assets in your web application. It receives a folder where your assets are on the disk and returns an http.Handler that serves these files.

This handler is capable of a few things that are useful for web development:
- Asset fingerprinting (to avoid caching issues)
- Hot code reloading (when the assets change on disk)

## Usage

The usage of the assets package is centered around the Assets manager instance.

```go

// Assets is the manager for the public assets
// it allows to watch for changes and reload the assets
// when changes are made.
Assets = assets.NewManager(public.Files)
...
// Register the assets handler
// This handler will serve the files from the public folder
...
	r.HandleFunc(Assets.HandlerPattern(), Assets.HandlerFn)
}
```

## Fingerprinting Helper
The assets manager provides a PathFor helper that can be used in your templates to use the fingerprinted version of an asset.

```html
<link rel="stylesheet" href="<%= assets.PathFor("/css/app.css") %>">
// will output something like
<link rel="stylesheet" href="/css/app-cafe123ff22112eedd.css">
```

## Hotcode Reloading
The assets managers provides a handler function capable of serving the files in the assets filesystem. This handler considers the `GO_ENV` variable to look in for files in disk before looking into the embedded filesystem passed.
