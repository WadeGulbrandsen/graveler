package main

import (
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

const rolls_per_game int = 231
const default_max_games int = 1000000000
const ones_to_win int = 177

func play_game() int {
	ones := 0
	for i := 0; i < rolls_per_game; i++ {
		if rand.Intn(4) == 0 { // since rand.Intn(4) will be from 0-3 consider 0 to be 1
			ones += 1
		}
	}
	return ones
}

func game_player(ch chan int) {
	for {
		ch <- play_game()
	}
}

func main() {
	max_games := default_max_games
	if len(os.Args) > 1 {
		x, err := strconv.Atoi(os.Args[1])
		if err == nil {
			max_games = x
		}
	}
	log.Printf("Starting Graveler Simulation for %v games using %v threads", max_games, runtime.NumCPU())
	start := time.Now()
	games := 0
	max_ones := 0
	ch := make(chan int)
	for i := 0; i < runtime.NumCPU(); i++ {
		go game_player(ch)
	}
	for games < max_games && max_ones < ones_to_win {
		max_ones = max(max_ones, <-ch)
		games += 1
		if games%1e8 == 0 {
			// print current highest roll every 100,000,000 games
			log.Printf("%v games simulated; current highest ones roll: %v", games, max_ones)
		}
	}
	elapsed := time.Since(start)
	log.Printf("Highest Ones Roll: %v", max_ones)
	log.Printf("Number of Roll Sessions: %v", games)
	log.Printf("Simulations Finished in: %s", elapsed)
}
