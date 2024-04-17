package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Player struct {
	name string
}

func (p *Player) Play(hit chan<- GameMsg) {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) <= 20 {
		hit <- GameMsg{p, strongHit}
	} else {
		hit <- GameMsg{p, regularHit}
	}
}

type GameMsg struct {
	p   *Player
	hit string
}

type GameScore struct {
	p1 int
	p2 int
}

func (s *GameScore) Add(p1 Player, msg GameMsg) {
	if msg.p.name == p1.name {
		s.p1 += 1
	} else {
		s.p2 += 1
	}
}

const (
	regularHit = "удар"
	strongHit  = "Сильный удар"
)

func main() {

	fmt.Println("Игра в пинг-понг")
	var p1Name string
	fmt.Println("Введите имя первого игрока?")
	fmt.Scan(&p1Name)

	var p2Name string
	fmt.Println("Введите имя второго игрока?")
	fmt.Scan(&p2Name)

	var finalScore int
	fmt.Println("До скольки очков играем?")
	fmt.Scan(&finalScore)
	fmt.Println("Начинаем!")

	p1 := Player{p1Name}
	p2 := Player{p2Name}
	s := GameScore{0, 0}
	game := make(chan GameMsg, 2)
	done := make(chan struct{})

	var wg sync.WaitGroup

	for s.p1 < finalScore && s.p2 < finalScore {
		wg.Add(2)

		go func() {
			defer wg.Done()
			p1.Play(game)
		}()

		go func() {
			defer wg.Done()
			p2.Play(game)
		}()

		select {
		case msg := <-game:
			if msg.hit == strongHit {
				fmt.Println(msg.p.name, msg.hit)
				s.Add(p1, msg)
				fmt.Println("Игрок ", msg.p.name, " получил +1 очко!")
				fmt.Println("Текущий счет: ", s.p1, " : ", s.p2)
				fmt.Println("Продолжаем!")
			} else {
				fmt.Println(msg.p.name, msg.hit)
			}
			time.Sleep(time.Millisecond * 500)
		case <-done:
			fmt.Println("мы в done")
			wg.Wait()
			return
		}
	}

	fmt.Println("Игра окончена!")
	fmt.Println("Итоговый счет: ", s.p1, " : ", s.p2)
	fmt.Println("Победитель: ", Winner(s, p1, p2).name)

}

func Winner(s GameScore, p1, p2 Player) Player {
	if s.p1 > s.p2 {
		return p1
	} else {
		return p2
	}

}
