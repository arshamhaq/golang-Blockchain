package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// This is a very good example of go routines and why we need waitgroups
// if we don't use waitgroups "the miner %v done!" won't print.
// when we say defer wg.done() this line will execute after the function is exited.

//don't put wg.Add() inside go func. If you call wg.Add(1) inside the goroutine, there’s a race condition:
//The main goroutine might call wg.Wait() before your Add happens
// Here in this code it won't happen even though I put ad inside the start_mining bc of the logic.

//you can also use wg.go() which is the short term of the add and done: less understandable but also less coding

// func (wg *WaitGroup) Go(f func()) {
//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         f()
//     }()
// }

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type Block struct {
	Data  string
	Nonce int
	Hash  [32]byte
}

type Results struct {
	b          Block
	miner_name string
}

func main() {
	c := make(chan Block)
	c2 := make(chan Results)

	miners_list := []string{"m1", "m2", "m3"}

	for range 5 {
		var new_block Block
		new_block.Data = randomString(7)
		go func(b Block) {
			c <- b
		}(new_block)
	}

	var wg sync.WaitGroup
	go miners_start_working(c, c2, miners_list, &wg)

	b_counter := 0
	var full_res [5]Results

	for b_counter <= 4 {
		full_res[b_counter] = <-c2
		fmt.Printf("Mined Block: %+v", full_res[b_counter].b)
		fmt.Printf("\n Miners's name: %v", full_res[b_counter].miner_name)
		fmt.Println("\n---------------------------------")
		b_counter++
	}

	close(c)
	wg.Wait()
	close(c2)
	fmt.Println("\n--------------------------------- \n All blocks done")
}

func mine(c chan Block, c2 chan Results, name string) bool {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var b Block
	b, ok := <-c
	if !ok {
		return true
	}
	for {
		b.Nonce = r.Int()

		data := fmt.Sprintf("%s%d", b.Data, b.Nonce)
		hash := sha256.Sum256([]byte(data))

		// Difficulty: first 6 hex characters must be "000000"
		// That means first 24 bits must be zero.
		if hash[0] == 0 && hash[1] == 0 && hash[2] == 0 {
			b.Hash = hash
			var res Results
			res.b = b
			res.miner_name = name
			c2 <- res
			return false
		}
	}
}

func miners_start_working(c chan Block, c2 chan Results, miners_lst []string, wg *sync.WaitGroup) {
	for i := range 3 {
		wg.Add(1) // could've also used wg.go func ...
		go func(i int) {
			// defer wg.done() ---> this is safer bc if we have panic or sth this still executes.
			for {
				not_ok := mine(c, c2, miners_lst[i])
				if not_ok {
					fmt.Printf("\nminer %v Done!", miners_lst[i])
					wg.Done()
					break
				}
			}
		}(i)
	}
}
