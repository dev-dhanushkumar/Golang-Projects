# Golang Masterclass - Basic Concepts

### prerequisites

Go installation & Configuration

```sh
Please check this link https://go.dev/doc/install

go version

```

##### Create go project 

```bash
mkdir go-project

cd go-project

go mod init github.com/newlinedeveloper/go-project

touch main.go

```

###### Hello world program

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}

```



### Go Basic Topics


#### Values

Go has various value types including strings, integers, floats, booleans


```go
package main
import "fmt"
func main() {

    fmt.Println("go" + "lang")

    fmt.Println("1+1 =", 1+1)
    fmt.Println("7.0/3.0 =", 7.0/3.0)

    fmt.Println(true && false)
    fmt.Println(true || false)
    fmt.Println(!true)
}

```

#### Variables

Variables are explicitly declared 


```go
package main
import "fmt"
func main() {

    var a = "initial"
    fmt.Println(a)

    var b, c int = 1, 2
    fmt.Println(b, c)

    var d = true
    fmt.Println(d)
    
    var e int
    fmt.Println(e)
   

    f := "apple"
    fmt.Println(f)
}

```



#### Constants : 

Go supports constants of character, string, boolean, and numeric values.

```go
package main

import (
    "fmt"
    "math"
)

const s string = "constant"

func main() {
    fmt.Println(s)

    const n = 500000000

    const d = 3e20 / n
    fmt.Println(d)

    fmt.Println(int64(d))

    fmt.Println(math.Sin(n))
}

```


#### For loop :

```go
package main
import "fmt"
func main() {

    i := 1
    for i <= 3 {
        fmt.Println(i)
        i = i + 1
    }

    for j := 7; j <= 9; j++ {
        fmt.Println(j)
    }

    for {
        fmt.Println("loop")
        break
    }

    for n := 0; n <= 5; n++ {
        if n%2 == 0 {
            continue
        }
        fmt.Println(n)
    }
}

```

#### If/Else :


```go
package main
import "fmt"
func main() {

    if 7%2 == 0 {
        fmt.Println("7 is even")
    } else {
        fmt.Println("7 is odd")
    }

    if 8%4 == 0 {
        fmt.Println("8 is divisible by 4")
    }

    if num := 9; num < 0 {
        fmt.Println(num, "is negative")
    } else if num < 10 {
        fmt.Println(num, "has 1 digit")
    } else {
        fmt.Println(num, "has multiple digits")
    }
}

```

#### Switch :

```go
package main
import (
    "fmt"
    "time"
)
func main() {

    i := 2
    fmt.Print("Write ", i, " as ")
    switch i {
    case 1:
        fmt.Println("one")
    case 2:
        fmt.Println("two")
    case 3:
        fmt.Println("three")
    }

    switch time.Now().Weekday() {
    case time.Saturday, time.Sunday:
        fmt.Println("It's the weekend")
    default:
        fmt.Println("It's a weekday")
    }

    t := time.Now()
    switch {
    case t.Hour() < 12:
        fmt.Println("It's before noon")
    default:
        fmt.Println("It's after noon")
    }

    whatAmI := func(i interface{}) {
        switch t := i.(type) {
        case bool:
            fmt.Println("I'm a bool")
        case int:
            fmt.Println("I'm an int")
        default:
            fmt.Printf("Don't know type %T\n", t)
        }
    }
    whatAmI(true)
    whatAmI(1)
    whatAmI("hey")
}

```

#### Arrays & Slices :

```go
package main
import "fmt"
func main() {

    var a [5]int
    fmt.Println("emp:", a)
    
    a[4] = 100
    fmt.Println("set:", a)
    fmt.Println("get:", a[4])

    fmt.Println("len:", len(a))

    b := [5]int{1, 2, 3, 4, 5}
    fmt.Println("dcl:", b)

    var twoD [2][3]int
    for i := 0; i < 2; i++ {
        for j := 0; j < 3; j++ {
            twoD[i][j] = i + j
        }
    }
    fmt.Println("2d: ", twoD)
    
    
    # slices
     l := s[2:5]
    fmt.Println("sl1:", l)
    

    l = s[:5]
    fmt.Println("sl2:", l)


    l = s[2:]
    fmt.Println("sl3:", l)
}

```

#### Maps : 

Maps are Go’s built-in associative data type (sometimes called hashes or dicts in other languages)

```go
package main

import "fmt"

func main() {

    m := make(map[string]int)

    m["k1"] = 7
    m["k2"] = 13

    fmt.Println("map:", m)

    v1 := m["k1"]
    fmt.Println("v1:", v1)

    v3 := m["k3"]
    fmt.Println("v3:", v3)

    fmt.Println("len:", len(m))

    delete(m, "k2")
    fmt.Println("map:", m)

    _, prs := m["k2"]
    fmt.Println("prs:", prs)

    n := map[string]int{"foo": 1, "bar": 2}
    fmt.Println("map:", n)
}

```

#### Functions:

```go
package main

import "fmt"

func plus(a int, b int) int {

    return a + b
}

func plusPlus(a, b, c int) int {
    return a + b + c
}

func vals() (int, int) {
    return 3, 7
}

func main() {

    res := plus(1, 2)
    fmt.Println("1+2 =", res)

    res = plusPlus(1, 2, 3)
    fmt.Println("1+2+3 =", res)
    
    a, b := vals()
    fmt.Println(a)
    fmt.Println(b)

    _, c := vals()
    fmt.Println(c)
    
}

```

#### Pointers

```go
package main

import "fmt"

func zeroval(ival int) {
    ival = 0
}

func zeroptr(iptr *int) {
    *iptr = 0
}

func main() {
    i := 1
    fmt.Println("initial:", i)

    zeroval(i)
    fmt.Println("zeroval:", i)

    zeroptr(&i)
    fmt.Println("zeroptr:", i)

    fmt.Println("pointer:", &i)
}

```

#### Structs

Go’s structs are typed collections of fields. They’re useful for grouping data together to form records.

```go
package main

import "fmt"

type person struct {
    name string
    age  int
}

func newPerson(name string) *person {

    p := person{name: name}
    p.age = 42
    return &p
}

func main() {

    fmt.Println(person{"Bob", 20})

    fmt.Println(person{name: "Alice", age: 30})

    fmt.Println(person{name: "Fred"})

    fmt.Println(&person{name: "Ann", age: 40})

    fmt.Println(newPerson("Jon"))

    s := person{name: "Sean", age: 50}
    fmt.Println(s.name)

    sp := &s
    fmt.Println(sp.age)

    sp.age = 51
    fmt.Println(sp.age)
}

```

#### Error Handling

In Go, it's common to handle errors by returning them as a separate value from the function. This is different from languages like Java and Ruby, which use exceptions, and C, which sometimes uses a single result/error value. Go's approach makes it clear which functions can return errors and allows you to use the same code to handle errors as you would for other tasks.


```go
package main

import (
    "errors"
    "fmt"
)

func f1(arg int) (int, error) {
    if arg == 42 {
        return -1, errors.New("can't work with 42")
    }
    return arg + 3, nil
}

func f2(arg int) (int, error) {
    if arg == 42 {
        return -1, fmt.Errorf("%d - can't work with it", arg)
    }
    return arg + 3, nil
}

func main() {
    for _, i := range []int{7, 42} {
        if r, err := f1(i); err != nil {
            fmt.Println("f1 failed:", err)
        } else {
            fmt.Println("f1 worked:", r)
        }
    }

    for _, i := range []int{7, 42} {
        if r, err := f2(i); err != nil {
            fmt.Println("f2 failed:", err)
        } else {
            fmt.Println("f2 worked:", r)
        }
    }

    if _, err := f2(42); err != nil {
        fmt.Println(err)
    }
}


```



#### Cross platform compilation


```sh
go build

env GOOS=target-OS GOARCH=target-architecture go build .

env GOOS=windows GOARCH=amd64 go build .

```



# Golang Masterclass - Advanced Concepts

### prerequisites


```sh
Please check this link https://go.dev/doc/install

go version


# MongoDB installation

Please check this link https://www.mongodb.com/docs/manual/installation/


```

##### Create go-api project 

```sh
mkdir go-api

cd go-api

go mod init github.com/newlinedeveloper/go-api


```

###### Install Required Packages

```sh
go get -u github.com/gorilla/mux go.mongodb.org/mongo-driver/mongo github.com/joho/godotenv github.com/go-playground/validator/v10

```

```go
github.com/gorilla/mux

go.mongodb.org/mongo-driver/mongo

github.com/joho/godotenv

github.com/go-playground/validator/v10

```


### Go simple Web server



```go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getMessage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": "Golang project setup test"})

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", getMessage).Methods("GET")
	fmt.Print("Server is running on port 8000 !!!")
	log.Fatal(http.ListenAndServe(":8000", router))

}

```

To run the application

```sh
go run main.go

Server is running on port 8000 !!!

```

#### Project structure



```
.
└── go-api/
    ├── Routes/
    │   └── member_routes.go
    ├── Controllers/
    │   └── member_controllers.go
    ├── Models/
    │   └── member_models.go
    ├── Configs/
    │   ├── env.go
    │   └── connection.go
    ├── Responses/
    │   └── member_responses.go
    ├── main.go
    ├── go.mod
    ├── go.sum
    ├── .env
    └── .env.example


```


#### MongoDB Connection setup

Create .env file and Add Mongo DB connection uri

```sh
MONGOURI=mongodb://localhost:27017
```

Create `env.go` file in `Configs` folder


```go
package configs

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

func EnvMongoURI() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    return os.Getenv("MONGOURI")
}

```


Create `connection.go` file in `Configs` folder


```go
package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB connection function
func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB Connected successfully !!!! ")
	return client
}

// // MongoDB Client instance
// var DB *mongo.Client = ConnectDB()

//Getting database collection
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("golang-masterclass").Collection(collectionName)
	return collection
}


```

in `main.go` file


```go

package main

import (
	"fmt"
	"log"
	"net/http"

	configs "github.com/newlinedeveloper/go-api/Configs"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// MongoDB Connection
	configs.ConnectDB()

	fmt.Print("Server is running on port 8000 !!!!")
	log.Fatal(http.ListenAndServe(":8000", router))
}



```

#### Create Member Model

Create `member_models.go` file in `Models` folder



```go

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Member struct {
	Id    primitive.ObjectID `json:"id,omitempty"`
	Name  string             `json:"name,omitempty" validate:"required"`
	Email string             `json:"email,omitempty" validate:"required"`
	City  string             `json:"city,omitempty" validate:"required"`
}


```


#### Create Member Response struct

Create `member_responses.go` file in `Responses` folder



```go

package responses

type MemberResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}


```


#### Create Members Api Routes

Create `member_routes.go` file in `Routes` folder


```go

package routes

import "github.com/gorilla/mux"

func MemberRoutes(router *mux.Router) {

}



```


import routes to `main.go` file



```go
package main

import (
	"fmt"
	"log"
	"net/http"

	configs "github.com/newlinedeveloper/go-api/Configs"

	routes "github.com/newlinedeveloper/go-api/Routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// MongoDB Connection
	configs.ConnectDB()

    // Imported Members routes
	routes.MemberRoutes(router)

	fmt.Print("Server is running on port 8000 !!!!")
	log.Fatal(http.ListenAndServe(":8000", router))
}


```



#### Create Members Controller functions

Create `member_controllers.go` file in `Controllers` folder



Create `CreateMember` function


```go

package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	configs "github.com/newlinedeveloper/go-api/Configs"
	models "github.com/newlinedeveloper/go-api/Models"
	responses "github.com/newlinedeveloper/go-api/Responses"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var memberCollection *mongo.Collection = configs.GetCollection(configs.DB, "members")
var validate = validator.New()

func CreateMember() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var member models.Member
		defer cancel()

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.MemberResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&member); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.MemberResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newUser := models.Member{
			Id:    primitive.NewObjectID(),
			Name:  member.Name,
			Email: member.Email,
			City:  member.City,
		}
		result, err := memberCollection.InsertOne(ctx, newUser)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.MemberResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		response := responses.MemberResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
}

```

update routes file


```go
package routes

import (
	"github.com/gorilla/mux"
	controllers "github.com/newlinedeveloper/go-api/Controllers"
)

func MemberRoutes(router *mux.Router) {

	router.HandleFunc("/member", controllers.CreateMember()).Methods("POST")

}

```



Create `GetMember` function


```go
func GetMember() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["id"]
		var user models.Member

		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := memberCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.MemberResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.MemberResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(rw).Encode(response)

	}
}

```

update Member Routes `member_routes.go` file


```go
router.HandleFunc("/member/{id}", controllers.GetMember()).Methods("GET")

```




Create `GetAllMembers` function


```go
func GetAllMembers() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var members []models.Member
		defer cancel()

		results, err := memberCollection.Find(ctx, bson.M{})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.MemberResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.Member
			if err = results.Decode(&singleUser); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.MemberResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
			}
			members = append(members, singleUser)

		}

		rw.WriteHeader(http.StatusOK)
		response := responses.MemberResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": members}}
		json.NewEncoder(rw).Encode(response)

	}
}

```

update Member Routes `member_routes.go` file


```go
router.HandleFunc("/members", controllers.GetAllMembers()).Methods("GET")

```



Create `UpdateMember` function


```go
func UpdateMember() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["id"]
		var user models.Member

		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.MemberResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.MemberResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		update := bson.M{"name": user.Name, "email": user.Email, "city": user.City}

		result, err := memberCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.MemberResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Get Updated member details
		var updatedMember models.Member

		if result.MatchedCount == 1 {
			err := memberCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedMember)

			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.MemberResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
				return
			}

		}

		rw.WriteHeader(http.StatusOK)
		response := responses.MemberResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedMember}}
		json.NewEncoder(rw).Encode(response)


	}
}

```

update Member Routes `member_routes.go` file


```go
router.HandleFunc("/member/{id}", controllers.UpdateMember()).Methods("PUT")

```




Create `DeleteMember` function


```go
func DeleteMember() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["id"]
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := memberCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.MemberResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount < 1 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.MemberResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Member Id Not found"}}
			json.NewEncoder(rw).Encode(response)
			return

		}

		rw.WriteHeader(http.StatusOK)
		response := responses.MemberResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": "Member deleted successfully"}}
		json.NewEncoder(rw).Encode(response)

	}
}

```

update Member Routes `member_routes.go` file


```go
router.HandleFunc("/member/{id}", controllers.DeleteMember()).Methods("DELETE")

```