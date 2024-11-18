package main

import (
	"context"
	"fmt"

	"github.com/tai-kun/surrealdb.go"
	"github.com/tai-kun/surrealdb.go/pkg/models"
)

func main() {
	ctx := context.Background()
	sdb, err := surrealdb.New(surrealdb.WithContext(ctx))
	if err != nil {
		panic(err)
	}

	if err := sdb.Connect("http://localhost:8000"); err != nil {
		panic(err)
	}
	defer sdb.Close()

	a := surrealdb.NewRootUserAuth("surreal", "surreal")
	if _, err := sdb.SignIn(a); err != nil {
		panic(err)
	}

	if err := sdb.Use("foo", "bar"); err != nil {
		panic(err)
	}

	now := models.NewDatetime()
	res, err := sdb.Query(
		`
		RETURN 'Hello Go SDK';
		CREATE ONLY user SET name = 'tai-kun', time = $time;
		`,
		map[string]any{
			"time": now,
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Len()) // 2

	var stmt0 string
	if err := res.Remove(0, &stmt0); err != nil {
		panic(err)
	}

	fmt.Println(stmt0)     // Hello Go SDK
	fmt.Println(res.Len()) // 1

	var stmt1 map[string]any
	if err := res.Remove(0, &stmt1); err != nil {
		panic(err)
	}

	fmt.Println(stmt1) // map[id:{user vhlyl9m4sol7q7mnsy0u} name:tai-kun time:2024-11-18 11:23:47.160342431 +0000 UTC]
}
