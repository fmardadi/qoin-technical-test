package main

import (
	"fmt"
	"math/rand"
	"qoin-technical-test/entity"
	"sort"
)

func main() {
	var totalPlayer int
	var totalDice int

	fmt.Print("Enter total player: ")
	fmt.Scanln(&totalPlayer)

	fmt.Print("Enter total dice: ")
	fmt.Scanln(&totalDice)

	players := generatePlayers(totalPlayer, totalDice)

	DiceRollGame := DiceRollGame{
		Players:     players,
		Leaderboard: make(map[string]int),
	}

	PlayerLeft = len(players)

	DiceRollGame.startGame()
	DiceRollGame.endGame()

}

var PlayerLeft int

type DiceRollGame struct {
	Players     []entity.Player
	Leaderboard map[string]int
	Round       int
	RoundResult map[string][]int
}

func generatePlayers(totalPlayer int, totalDice int) []entity.Player {
	var players []entity.Player

	for i := 1; i <= totalPlayer; i++ {
		player := entity.Player{
			ID:       i,
			Name:     fmt.Sprintf("Player #%d", i),
			DiceLeft: totalDice,
			IsActive: true,
		}
		players = append(players, player)
	}

	return players
}

func (d DiceRollGame) startGame() {
	fmt.Println("----- Game is starting -----")

	d.Round = 1
	for PlayerLeft > 1 {
		d.startRound()
		d.Round++
		fmt.Println("------------------------------")
	}

}

func (d DiceRollGame) startRound() {
	fmt.Printf("Round %d is starting \n", d.Round)

	d.RoundResult = make(map[string][]int)

	for _, player := range d.Players {
		for j := 0; j < player.DiceLeft; j++ {
			d.RoundResult[player.Name] = append(d.RoundResult[player.Name], d.rollTheDice())
		}
	}

	d.printResult(false)
	d.evaluateRound()
	d.adjustDiceLeft()
	d.printResult(true)
	d.printLeaderboard()
	d.checkActivePlayers()
}

func (d DiceRollGame) rollTheDice() int {
	var dice = []int{1, 2, 3, 4, 5, 6}
	return dice[rand.Intn(len(dice))]
}

func (d DiceRollGame) printResult(isEvaluate bool) {
	if isEvaluate {
		fmt.Printf("Evaluate Round %d \n", d.Round)
	} else {
		fmt.Printf("Round %d result \n", d.Round)
	}

	players := make([]string, 0, len(d.RoundResult))
	for k := range d.RoundResult {
		players = append(players, k)
	}
	sort.Strings(players)

	for _, player := range players {
		dice := d.RoundResult[player]
		fmt.Printf("%s (%d) : %d \n", player, d.Leaderboard[player], dice)

	}
}

func (d DiceRollGame) evaluateRound() {
	for player, dice := range d.RoundResult {

		var diceLeft int

		// // get player dice left
		for _, val := range d.Players {
			if val.Name == player {
				diceLeft = val.DiceLeft
				break
			}
		}

		var newDice []int
		for _, val := range dice {
			if val == 6 {
				d.Leaderboard[player] += 1
				diceLeft--
			} else if val == 1 {
				// pass the dice to the next player
				d.passTheDice(player)
				diceLeft--
			} else {
				newDice = append(newDice, val)
			}
		}

		d.RoundResult[player] = newDice

		// update dice left
		for i, val := range d.Players {
			if val.Name == player {
				differ := diceLeft - len(newDice)
				d.Players[i].DiceLeft = len(newDice) + differ
				break
			}
		}
	}
}

func (d DiceRollGame) passTheDice(playerName string) {
	for i, val := range d.Players {
		if playerName == val.Name {

			if i == len(d.Players)-1 {
				// Pass to next player from the start index, because its the last index player
				for j := 0; j < len(d.Players); j++ {
					if d.Players[j].IsActive && j != i {
						d.Players[j].DiceLeft += 1
						break
					}
				}
				break
			} else {
				// Pass to the next player

				// counter for safety
				maxCounter := len(d.Players)
				count := 0
				for j := i + 1; j < len(d.Players); j++ {
					if d.Players[j].IsActive && j != i {
						d.Players[j].DiceLeft += 1
						break
					}

					if j == len(d.Players)-1 {
						j = 0
					}

					if count == maxCounter {
						break
					}
					count++
				}
				break
			}

		}
	}
}

// this function is for adding dice to another player when dice is passed
func (d DiceRollGame) adjustDiceLeft() {

	for _, val := range d.Players {
		for player, dice := range d.RoundResult {

			if val.Name == player {
				differ := val.DiceLeft - len(dice)
				if differ > 0 {
					for i := 0; i < differ; i++ {
						d.RoundResult[player] = append(d.RoundResult[player], 1)
					}
				}

			}
		}
	}

}

func (d DiceRollGame) printLeaderboard() {
	fmt.Println("----- Leaderboard -----")

	players := make([]string, 0, len(d.Leaderboard))
	for k := range d.Leaderboard {
		players = append(players, k)
	}
	sort.Strings(players)

	for _, player := range players {
		points := d.Leaderboard[player]
		fmt.Printf("%s : %d \n", player, points)
	}
}

func (d DiceRollGame) checkActivePlayers() {
	var activePlayerID []int
	for i, val := range d.Players {
		if val.DiceLeft > 0 {
			activePlayerID = append(activePlayerID, val.ID)
		} else {
			// set not active player that has 0 dice left
			d.Players[i].IsActive = false
		}
	}

	var updatedPlayers []entity.Player

	for _, val := range d.Players {
		for _, id := range activePlayerID {
			if val.ID == id {
				updatedPlayers = append(updatedPlayers, val)
			}
		}
	}

	d.Players = updatedPlayers
	PlayerLeft = len(updatedPlayers)
}

func (d DiceRollGame) endGame() {
	var topScore int
	for _, value := range d.Leaderboard {
		if value > topScore {
			topScore = value
		}
	}

	var winners []string
	for key, value := range d.Leaderboard {
		if value == topScore {
			winners = append(winners, key)
		}
	}

	fmt.Printf("The winner is %s by %d point(s) \n", winners, topScore)
	fmt.Println("----- Game ended -----")
}
