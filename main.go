package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type GameMsg struct {
	name string
	hit  string
}

type GameScore struct {
	p1 int
	p2 int
}

func (s *GameScore) Add(name string, msg GameMsg) {
	if msg.name == name {
		s.p1 += 1
	} else {
		s.p2 += 1
	}
}

const (
	regularHit = "удар"
	strongHit  = "сильный удар!"
)

func main() {
	p1Name, p2Name, finalScore := preparingToGame()
	s := GameScore{0, 0}
	game := make(chan GameMsg, 2)
	done := make(chan struct{})
	var wg sync.WaitGroup
	for s.p1 < finalScore && s.p2 < finalScore {
		wg.Add(2)
		go func() {
			defer wg.Done()
			Play(game, p1Name)
		}()
		go func() {
			defer wg.Done()
			Play(game, p2Name)
		}()
		select {
		case msg := <-game:
			if msg.hit == strongHit {
				fmt.Println(msg.name, msg.hit)
				s.Add(p1Name, msg)
				fmt.Println("Игрок ", msg.name, " получил +1 очко!")
				fmt.Println("Текущий счет: ", s.p1, " : ", s.p2)
				fmt.Println("Продолжаем!")
			} else {
				fmt.Println(msg.name, msg.hit)
			}
			time.Sleep(time.Millisecond * 500)
		case <-done:
			wg.Wait()
			return
		}
	}
	printGameOver(s, p1Name, p2Name)
}

func preparingToGame() (string, string, int) {
	var pn1, pn2 string
	var fs int
	fmt.Println("Игра в пинг-понг")
	fmt.Println("Введите имя первого игрока?")
	fmt.Scan(&pn1)

	fmt.Println("Введите имя второго игрока?")
	fmt.Scan(&pn2)

	fmt.Println("До скольки очков играем?")
	fmt.Scan(&fs)
	fmt.Println("Начинаем!")

	return pn1, pn2, fs
}

func Play(hit chan<- GameMsg, name string) {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) <= 20 {
		hit <- GameMsg{name, strongHit}
	} else {
		hit <- GameMsg{name, regularHit}
	}
}

func printGameOver(s GameScore, p1, p2 string) {
	fmt.Println("Игра окончена!")
	fmt.Println("Итоговый счет: ", s.p1, " : ", s.p2)
	var winner string
	if s.p1 > s.p2 {
		winner = p1
	} else {
		winner = p2
	}
	fmt.Println("Победитель: ", winner)

}
