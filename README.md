# sessions

[![Run CI Lint](https://github.com/forthxu/sessions/actions/workflows/lint.yml/badge.svg?branch=master)](https://github.com/forthxu/sessions/actions/workflows/lint.yml)
[![Run Testing](https://github.com/forthxu/sessions/actions/workflows/testing.yml/badge.svg?branch=master)](https://github.com/forthxu/sessions/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/forthxu/sessions/branch/master/graph/badge.svg)](https://codecov.io/gh/forthxu/sessions)
[![Go Report Card](https://goreportcard.com/badge/github.com/forthxu/sessions)](https://goreportcard.com/report/github.com/forthxu/sessions)
[![GoDoc](https://godoc.org/github.com/forthxu/sessions?status.svg)](https://godoc.org/github.com/forthxu/sessions)

Gin middleware for session management with multi-backend support:

- [cookie-based](#cookie-based)
- [Redis](#redis)
- [memcached](#memcached)
- [MongoDB](#mongodb)
- [GORM](#gorm)
- [memstore](#memstore)
- [PostgreSQL](#postgresql)

## Usage

### Start using it

Download and install it:

```bash
go get github.com/forthxu/sessions
```

Import it in your code:

```go
import "github.com/forthxu/sessions"
```

## Basic Examples

### single session

```go
package main

import (
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/cookie"
  "github.com/forthxu/gin"
)

func main() {
  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/hello", func(c *gin.Context) {
    session := sessions.Default(c)

    if session.Get("hello") != "world" {
      session.Set("hello", "world")
      session.Save()
    }

    c.JSON(200, gin.H{"hello": session.Get("hello")})
  })
  r.Run(":8000")
}
```

### multiple sessions

```go
package main

import (
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/cookie"
  "github.com/forthxu/gin"
)

func main() {
  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  sessionNames := []string{"a", "b"}
  r.Use(sessions.SessionsMany(sessionNames, store))

  r.GET("/hello", func(c *gin.Context) {
    sessionA := sessions.DefaultMany(c, "a")
    sessionB := sessions.DefaultMany(c, "b")

    if sessionA.Get("hello") != "world!" {
      sessionA.Set("hello", "world!")
      sessionA.Save()
    }

    if sessionB.Get("hello") != "world?" {
      sessionB.Set("hello", "world?")
      sessionB.Save()
    }

    c.JSON(200, gin.H{
      "a": sessionA.Get("hello"),
      "b": sessionB.Get("hello"),
    })
  })
  r.Run(":8000")
}
```

## Backend Examples

### cookie-based

```go
package main

import (
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/cookie"
  "github.com/forthxu/gin"
)

func main() {
  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### Redis

```go
package main

import (
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/redis"
  "github.com/forthxu/gin"
)

func main() {
  r := gin.Default()
  store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### Memcached

#### ASCII Protocol

```go
package main

import (
  "github.com/bradfitz/gomemcache/memcache"
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/memcached"
  "github.com/forthxu/gin"
)

func main() {
  r := gin.Default()
  store := memcached.NewStore(memcache.New("localhost:11211"), "", []byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

#### Binary protocol (with optional SASL authentication)

```go
package main

import (
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/memcached"
  "github.com/forthxu/gin"
  "github.com/memcachier/mc"
)

func main() {
  r := gin.Default()
  client := mc.NewMC("localhost:11211", "username", "password")
  store := memcached.NewMemcacheStore(client, "", []byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### MongoDB

#### mgo

```go
package main

import (
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/mongo/mongomgo"
  "github.com/forthxu/gin"
  "github.com/globalsign/mgo"
)

func main() {
  r := gin.Default()
  session, err := mgo.Dial("localhost:27017/test")
  if err != nil {
    // handle err
  }

  c := session.DB("").C("sessions")
  store := mongomgo.NewStore(c, 3600, true, []byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

#### mongo-driver

```go
package main

import (
  "context"
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/mongo/mongodriver"
  "github.com/forthxu/gin"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
  r := gin.Default()
  mongoOptions := options.Client().ApplyURI("mongodb://localhost:27017")
  client, err := mongo.NewClient(mongoOptions)
  if err != nil {
    // handle err
  }

  if err := client.Connect(context.Background()); err != nil {
    // handle err
  }

  c := client.Database("test").Collection("sessions")
  store := mongodriver.NewStore(c, 3600, true, []byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### memstore

```go
package main

import (
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/memstore"
  "github.com/forthxu/gin"
)

func main() {
  r := gin.Default()
  store := memstore.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### GORM

```go
package main

import (
  "github.com/forthxu/sessions"
  gormsessions "github.com/forthxu/sessions/gorm"
  "github.com/forthxu/gin"
  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
)

func main() {
  db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic(err)
  }
  store := gormsessions.NewStore(db, true, []byte("secret"))

  r := gin.Default()
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### PostgreSQL

```go
package main

import (
  "database/sql"
  "github.com/forthxu/sessions"
  "github.com/forthxu/sessions/postgres"
  "github.com/forthxu/gin"
)

func main() {
  r := gin.Default()
  db, err := sql.Open("postgres", "postgresql://username:password@localhost:5432/database")
  if err != nil {
    // handle err
  }

  store, err := postgres.NewStore(db, []byte("secret"))
  if err != nil {
    // handle err
  }

  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```
