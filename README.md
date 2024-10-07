# Golang-Masterclass
Golang Masterclass - Basic Concepts

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



