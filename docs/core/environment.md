---
title: Environment Variables
index: 5
---

Leapkit provides a handy tool to load environment variables from a `.env` file. This feature is useful when you want to load your environment variables from a file instead of setting them directly in your system, like when you're in development mode.

### Creating a `.env` file

To use it you can create a `.env` file in the root of your project and add your environment variables in the following format:

```.env
# .env
PORT=8080
```

### Adding Comments
You can also add comments to your `.env` file by starting the line with a `#` character.
```.env
# .env
# This is a comment
PORT=8080
```
### Multi-line values

Leapkit supports multi-line values by using the quoted string syntax:

```.env
GITHUB_SECRET_KEY="-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAqTmwQppL07nBl/0TEQ5sHcqj/Iz9BmuaaEu26jMXYt1QttHn
-----END RSA PRIVATE KEY-----"
```

## Loading the environment variables
To load the environment variables from the .env file into your application, perform an underscore import in your main.go file:

```go
// main.go
import _ "github.com/leapkit/leapkit/core/envload"
```

This will load the environment variables from the `.env` file into your application.
