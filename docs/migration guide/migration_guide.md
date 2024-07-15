---
title: "Migration Guide"
---

## Removing `gofrs/uuid/v5` package.

In the recent Leapkit updates, we have removed the `gofrs/uuid/v5` package to avoid promoting the use of a specific UUID package and version (v4). We want you to select the **UUID** package and version that you prefer to use in your app.

### What Has Changed?

This change impacts the followig packages.

- `leapkit/core/form`: The package no longer supports decoding of **UUID** and **[]UUID** types.
- `leapkit/core/form/validate`: The *ValidUUID* ValidatorFn was removed.
- `leapkit/core/server/session`:, the **UUID** type is not being registered in the session store.

### How to Migrate

You can continue using **UUIDs** throughout the different Leapkit packages by using the packages functions, which allow you to process custom types.

#### Decoding UUID types in the `leapkit/core/form` package

You can register the following functions to allow to decode **UUID** and **[]UUID** types using the package:

```go
package app

import (
	"github.com/gofrs/uuid/v5"
	"github.com/leapkit/leapkit/core/form"
)

func init() {
	form.RegisterCustomTypeFunc(decodeUUID, uuid.UUID{})
	form.RegisterCustomTypeFunc(decodeUUIDSlice, []uuid.UUID{})
}

func decodeUUID(vals []string) (any, error) {
	id, err := uuid.FromString(vals[0])
	if err != nil {
		err = fmt.Errorf("error parsing uuid: %w", err)
	}

	return id, err
}

func decodeUUIDSlice(vals []string) (any, error) {
	var ids []uuid.UUID

	for _, val := range vals {
		id, err := uuid.FromString(val)
		if err != nil {
			err = fmt.Errorf("error parsing uuid: %w", err)
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

```

#### Validating the UUID type with `leapkit/core/form/validations`

You can create a custom validation in which you can handle your validations for :

```go
package app

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/leapkit/leapkit/core/form"
	"github.com/leapkit/leapkit/core/form/validate"
)

func validUUID() validate.ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			if !uuid.FromStringOrNil(val).IsNil() {
				continue
			}

			return fmt.Errorf("'%s' is not a valid uuid", val)
		}

		return nil
	}
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	rules := validate.Fields(
		validate.Field("user_id", validUUID()),
        // ...
	)

	verrs := form.Validate(r, rules)
	if len(verrs) > 0 {
		// handle the error
	}
	// ...
}
```

#### Registering UUID type in `leapkit/core/server/session`

We can use the new `RegisterSessionTypes` function to register your UUID type to serialize/deserialize values into your session.

```go
package app

import (
	"github.com/gofrs/uuid/v5"
	"github.com/leapkit/leapkit/core/server/session"
)


func init() {
	session.RegisterSessionTypes(uuid.UUID{})
}
```