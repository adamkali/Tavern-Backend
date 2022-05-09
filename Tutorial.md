### Getting started

First run:

```
$ go mod init <project name>
```

This will creat the go project and create a `go.mod` file. When you import different packages using `$ go get <package>` they will be created stored here.

Now let us create a `main.go` file. It will be very simple:
```go

package main

import (
	"fmt"
	"log"
	"net/http"
)

// Make a handler to write a response to the client.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! I am calling out to you from: %s!", r.URL.Path[1:])
}

// Make the server proper.
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

```

All this is really doing is just creating a http server that when the `/` path is pinged it will write whatever path you passed next to it. For exampe if you put in postman/thunderclient/curl `GET http://locahost:8000/ping/` it will responde with: `Hello, World! I am calling out to you from ping!`.

### Using Data

The first thing that I want to do is make 5 endpoints to get The user for my app:

* [] GET 	api/user 		--- get all the users.
* [] GET 	api/user/:id 		--- get the user with ID :id.
* [] POST 	api/user		--- post a user/users using a body.
* [] PUT	api/user/:id		--- update the user with ID :id using a body.
* [] Delete	api/user/:id		--- delete the user with ID :id.

To do this I want to make first a User struct containing what I need. Use my struct as a guiding point for you to make your own data structures as well.

```go

package models

type User struct {
	ID              string   `json:"id"`
	Username        string   `json:"username"`
	Bio             string   `json:"bio"`
	Tags            []string `json:"tags"`
	PlayerPrefrence string   `json:"player_prefrence"`

	Plots struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Plot   string `json:"plot"`
		UserFk string `json:"user_fk"`
	} `json:"plots,omitempty"`

	Character struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Backstory    string `json:"backstory"`
		Bio          string `json:"bio"`
		Strength     int    `json:"strength"`
		Dexterity    int    `json:"dexterity"`
		Constitution int    `json:"constitution"`
		Intellegence int    `json:"intellegence"`
		Wisdom       int    `json:"wisdom"`
		Charisma     int    `json:"charisma"`
		UserFk       string `json:"user_fk"`
	} `json:"character,omitempty"`

}

type Users []User

type UserDetailedResponse struct {
	Data       User   `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type UsersDetailedResponse struct {
	Data       Users  `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

```

Now I like to wrap my responses in a Detailed response class in C# and Python, so this is really a force of habit for me; however, This allows you to extend error messages and successes with your response so I urge you to do it. This will help in large scale projects and testing in your future projects.

Also, I like to segment my code so that the models/types are in its own folder and so are the controller functions. So if you have a project tree like the followig:

```
	project
	|
	|__> types
	|   |
	|   |__> types.go
	|
	|__> main.go
```

You can bring in your the strucures from `types.go` into `main.go` by importing them:

```go

import (
	"<your project name>/<name of the package specified in types.go>"
)

```

I suggest that you keep the project name to be the same as the directory name, for your own sanity 😊.


