package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	width       = 20
	height      = 10
	snacksCount = 5 // Adjust difficulty by number of snacks 1 being hardest
	scoreFile   = "top_score.txt"
)

type point struct {
	x, y int
}

type snake struct {
	body      []point
	direction point
}

type game struct {
	snake          *snake
	snacks         []point
	speed          time.Duration
	snacksCaptured int
	topScore       int
}

func newSnake() *snake {
	return &snake{
		body:      []point{{x: width / 2, y: height / 2}},
		direction: point{x: 1, y: 0}, // Initially moving right
	}
}

func newGame() *game {
	s := newSnake()
	g := &game{
		snake:  s,
		speed:  200, // initial speed in milliseconds
		snacks: make([]point, 0, snacksCount),
	}
	g.loadTopScore()
	g.populateSnacks()
	return g
}

func (g *game) loadTopScore() {
	data, err := ioutil.ReadFile(scoreFile)
	if err == nil {
		score, err := strconv.Atoi(string(data))
		if err == nil {
			g.topScore = score
		}
	}
}

func (g *game) saveTopScore() {
	if g.snacksCaptured > g.topScore {
		g.topScore = g.snacksCaptured
		ioutil.WriteFile(scoreFile, []byte(strconv.Itoa(g.topScore)), 0644)
	}
}

func (g *game) populateSnacks() {
	for len(g.snacks) < snacksCount {
		snack := point{x: rand.Intn(width), y: rand.Intn(height)}
		if !g.isOccupied(snack) {
			g.snacks = append(g.snacks, snack)
		}
	}
}

func (g *game) isOccupied(p point) bool {
	for _, b := range g.snake.body {
		if b == p {
			return true
		}
	}
	for _, s := range g.snacks {
		if s == p {
			return true
		}
	}
	return false
}

func (s *snake) move(g *game) {
	head := s.body[0]
	newHead := point{x: head.x + s.direction.x, y: head.y + s.direction.y}

	if newHead.x < 0 || newHead.x >= width || newHead.y < 0 || newHead.y >= height {
		fmt.Printf("Game Over! Snacks captured: %d\n", g.snacksCaptured)
		fmt.Printf("Top score: %d\n", g.topScore)
		g.saveTopScore()
		os.Exit(0)
	}

	for _, b := range s.body[1:] {
		if b == newHead {
			fmt.Printf("Game Over! Snacks captured: %d\n", g.snacksCaptured)
			fmt.Printf("Top score: %d\n", g.topScore)
			g.saveTopScore()
			os.Exit(0)
		}
	}

	s.body = append([]point{newHead}, s.body...)

	for i, snack := range g.snacks {
		if newHead == snack {
			s.body = append(s.body, newHead) // Grow the snake
			g.snacksCaptured++
			g.snacks = append(g.snacks[:i], g.snacks[i+1:]...) // Remove the captured snack
			break
		}
	}

	if len(g.snacks) < snacksCount {
		g.populateSnacks()
	}
}

func changeDirection(g *game) {
	if char, _, err := keyboard.GetSingleKey(); err == nil {
		switch char {
		case 'a':
			if g.snake.direction.x == 0 {
				g.snake.direction = point{x: -1, y: 0}
			}
		case 'd':
			if g.snake.direction.x == 0 {
				g.snake.direction = point{x: 1, y: 0}
			}
		case 'w':
			if g.snake.direction.y == 0 {
				g.snake.direction = point{x: 0, y: -1}
			}
		case 's':
			if g.snake.direction.y == 0 {
				g.snake.direction = point{x: 0, y: 1}
			}
		case 'q': // Quit the game
			fmt.Printf("Exiting game\n")
			os.Exit(0)
		}
	}
}

func render(g *game) {
	clearScreen()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			printed := false
			for _, b := range g.snake.body {
				if b.x == x && b.y == y {
					fmt.Print("#")
					printed = true
					break
				}
			}
			if !printed {
				snackHere := false
				for _, s := range g.snacks {
					if s.x == x && s.y == y {
						fmt.Print("O")
						snackHere = true
						break
					}
				}
				if !snackHere {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	g := newGame()
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Game instructions
	fmt.Println("Use 'WASD' to move the snake. Press 'Q' to quit.")
	fmt.Println("Snake is # and will grow when snacks caputured. Capture the snacks use O symbol")
	fmt.Println("Starts in 5 seconds")
	time.Sleep(4 * time.Second)

	clearScreen()

	for {
		render(g)
		changeDirection(g)
		g.snake.move(g)
		time.Sleep(g.speed * time.Millisecond)
	}
}
