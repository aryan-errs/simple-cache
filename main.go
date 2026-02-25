package main

import (
	"aryan-errs/simple-cache/cache"
	"fmt"
)

func main() {
	cache := cache.CreateCache(5)
	ok := cache.Put(1, 2)
	if !ok {
		fmt.Println("key not added")
	}
	fmt.Println(cache.Get(1))
}
