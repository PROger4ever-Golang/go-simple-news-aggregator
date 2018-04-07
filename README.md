go-simple-news-aggregator
===============================

The app demonstrates:

* Using an SQL (SQLite) database and configuring the Revel DB module.
* Using the third party [GORP](https://github.com/coopernurse/gorp) *ORM-ish* library
* [Interceptors](../manual/interceptors.html) for checking that an user is logged in.
* Using [validation](../manual/validation) and displaying inline errors


# Database Install and Setup
This example used [sqlite](https://www.sqlite.org/), (Alternatively can use mysql, postgres, etc.)

## sqlite Installation

- The article app uses [go-sqlite3](https://github.com/mattn/go-sqlite3) database driver, which depends on the C library

### Install sqlite on OSX:

1. Install [Homebrew](http://mxcl.github.com/homebrew/) if you don't already have it.
2. Install pkg-config and sqlite3:

~~~
$ brew install pkgconfig sqlite3
~~~

### Install sqlite on Ubuntu:
```sh
$ sudo apt-get install sqlite3 libsqlite3-dev
```

Once SQLite is installed, it will be possible to run the article app:
```sh
	$ revel run github.com/PROger4ever/go-simple-news-aggregator
```

## Database / Gorp Plugin

[`app/controllers/gorp.go`](https://github.com/revel/examples/blob/master/article/app/controllers/gorp.go) defines `GorpPlugin`, which is a plugin that does a couple things:

* **`OnAppStart`** -  Uses the DB module to open a SQLite in-memory database, create the `User`, `Article`, and `Source` tables, and insert some test records.
* **BeforeRequest** -  Begins a transaction and stores the Transaction on the Controller
* **AfterRequest** -  Commits the transaction, or [panics](https://github.com/golang/go/wiki/PanicAndRecover) if there was an error.
* **OnException** -  Rolls back the transaction


## Interceptors

[`app/controllers/init.go`](https://github.com/revel/examples/blob/master/article/app/controllers/init.go)
registers the [interceptors](../manual/interceptors.html) that runs before each action:

```go
func init() {
	revel.OnAppStart(Init)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	revel.InterceptMethod(Sources.checkUser, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}
```

As an example, `checkUser` looks up the username in the `session` and `redirect`s
the user to log in if they do not have a `session` cookie.

```go
func (c Sources) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(Application.Index)
	}
	return nil
}
```

## Validation

The article app does quite a bit of validation.

Revel applies the validation and records errors using the name of the
validated variable (unless overridden).

The `field` template helper looks for errors in the validation context, using
the field name as the key.
