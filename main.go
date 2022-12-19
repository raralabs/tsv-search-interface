package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/raralabs/rara-search/pkg/client"
)

func main() {
	var wg sync.WaitGroup
	// id := 1
	// for i := 1; i <= 100; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		client := client.NewClient()
	// 		defer client.CloseConnection()
	// 		for j := 1; j <= 100000; j++ {
	// 			client.Index(fmt.Sprintf("rara%d", 1), uuid.NewString(), generate(16), map[string]interface{}{"id": id}, map[string]interface{}{"first_name": generate(8), "last_name": generate(8)})
	// 			id++
	// 		}
	// 	}()
	// }
	// wg.Wait()

	// fmt.Println(id, err)
	// d, err := client.SearchByField("rara1", map[string]interface{}{"first": 1})

	client := client.NewClient()
	defer client.CloseConnection()
	wg.Add(5)
	go func() {
		d, _ := client.Search("rara8", "2oWL", 1, 10)
		fmt.Printf("%+v", d)
	}()
	go func() {
		d, _ := client.Search("rara1", "abc", 1, 10)
		fmt.Printf("%+v", d)
	}()
	go func() {
		d, _ := client.Search("rara2", "abc", 1, 10)
		fmt.Printf("%+v", d)
	}()
	go func() {
		d, _ := client.Search("rara3", "abc", 1, 10)
		fmt.Printf("%+v", d)
	}()
	go func() {
		d, _ := client.Search("rara4", "abc", 1, 10)
		fmt.Printf("%+v", d)
	}()

	wg.Wait()
	// fmt.Printf("%+v %v", d, err)

	// id, err := client.Delete("rara1", "abcd", "member")
	// fmt.Println(id, err)

}

func generate(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}
