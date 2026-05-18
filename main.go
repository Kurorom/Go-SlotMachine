package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	Text       string
	Answer     string
	difficulty string
}

type QuizDatabase struct {
	Categories map[string][]Question
}
type DifficultyConfig struct {
	Name       string
	Multiplier uint
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getName() string {

	for {
		fmt.Print("what should i call you?\n")
		fmt.Printf("Enter your name: ")
		name := readInput()
		if name == "nouhaila" {
			fmt.Printf("you are a lovely girl! ⸜(｡˃ ᵕ ˂ )⸝♡\n")
		}
		if name != "" {
			fmt.Printf("Welcome to the game, %s!\n", name)
			return name
		}
		fmt.Println("Error: You must enter a valid name.")
	}
}

func getbet(balance uint) uint {

	for {

		fmt.Printf("Enter your bet or 0 to exit: ")
		input := readInput()
		val, err := strconv.ParseUint(input, 10, 32)
		bet := uint(val)

		if err != nil {
			fmt.Println("invalid input! please enter a positive whole number")
			continue
		}
		if bet == 0 {
			return 0
		}
		if bet > balance {
			fmt.Printf("Insuficient bakance! Your current balance is: %d\n", balance)
			continue
		}
		return bet
	}
}

func generateSymbolArray(Symbols map[string]uint) []string {
	symbolArr := []string{}
	for symbol, count := range Symbols {
		for i := uint(0); i < count; i++ {
			symbolArr = append(symbolArr, symbol)
		}
	}
	return symbolArr
}

func getRandomNumber(min int, max int) int {
	randomNumber := min + rand.Intn(max-min+1)
	return randomNumber
}

func spin(reel []string, rows int, cols int) [][]string {
	result := [][]string{}
	for i := 0; i < rows; i++ {
		result = append(result, []string{})
	}

	for col := 0; col < cols; col++ {
		selected := map[int]bool{}
		for row := 0; row < rows; row++ {
			randomIndex := getRandomNumber(0, len(reel)-1)
			_, exists := selected[randomIndex]
			if !exists {
				selected[randomIndex] = true
				result[row] = append(result[row], reel[randomIndex])

			}
		}
	}
	return result
}
func printSpinResult(gameResult [][]string) {
	fmt.Println("Spinning...")
	time.Sleep(1 * time.Second)
	fmt.Println("--- SPIN RESULTS ---")
	for _, row := range gameResult {
		for _, symbol := range row {
			fmt.Printf("| %s ", symbol)
		}
		fmt.Println("|")
	}
	fmt.Println("--------------------")
}

func calculateWinnings(gameResult [][]string, multipliers map[string]uint) uint {
	winning := uint(0)
	for _, row := range gameResult {
		if row[0] == row[1] && row[1] == row[2] {
			symbol := row[0]
			if symbol == "🌀" {
				println("congrats! you got another spin! 🔄")
				return 1
			}
			winning += multipliers[symbol]
		}

	}
	return winning
}

func getCategory(db QuizDatabase) string {
	fmt.Printf(" chose a category , if you answer the question correctly a sum will be added to your balance depending on the diffuclty of the question! \n")
	fmt.Printf("available categories :\n")
	for category := range db.Categories {
		fmt.Printf("%s\n", category)
	}
	categoryInput := strings.ToLower(readInput())
	_, exists := db.Categories[categoryInput]
	for !exists {
		fmt.Println("Invalid category chosen! Enter a valid category:")
		categoryInput = strings.ToLower(readInput())
		_, exists = db.Categories[categoryInput]

	}
	return categoryInput
}

func getDifficlty() DifficultyConfig {
	fmt.Printf("now choose a diffculty , type 1 for easy , 2 for medium , 3 for hard :\n")
	difficultyInput := strings.ToLower(readInput())
	diffConfigMap := map[string]DifficultyConfig{
		"1": {Name: "easy", Multiplier: 1},
		"2": {Name: "medium", Multiplier: 2},
		"3": {Name: "hard", Multiplier: 5},
	}

	targetDifficulty, ok := diffConfigMap[difficultyInput]
	for !ok {
		fmt.Println("Invalid difficulty choice!")
		difficultyInput := readInput()
		targetDifficulty, ok = diffConfigMap[difficultyInput]

	}
	return targetDifficulty
}

func getQuestions(db QuizDatabase) ([]Question, DifficultyConfig) {

	questions := db.Categories[getCategory(db)]
	targetDifficulty := getDifficlty()
	var matchingQuestions []Question
	for _, questions := range questions {
		if questions.difficulty == targetDifficulty.Name {
			matchingQuestions = append(matchingQuestions, questions)
			break
		}
	}
	for len(matchingQuestions) == 0 {
		fmt.Printf("No %s questions found in this category. chose a diffirent difficulty:\n", targetDifficulty.Name)
		targetDifficulty := getDifficlty()
		var matchingQuestions []Question
		for _, questions := range questions {
			if questions.difficulty == targetDifficulty.Name {
				matchingQuestions = append(matchingQuestions, questions)
				break
			}
		}

	}
	return matchingQuestions, targetDifficulty
}

func getQuizz(db QuizDatabase, balance *uint) bool {
	fmt.Printf("oh no you ran out of money ! , would you like more ? 🤭 ( 1 to procceed 0 to exit )\n")
	inpout := readInput()
	if inpout != "1" {
		return false
	}

	matchingQuestions, targetDifficulty := getQuestions(db)
	selectedQuestion := matchingQuestions[rand.Intn(len(matchingQuestions))]
	fmt.Printf("Questions: %s\n", selectedQuestion.Text)
	playerAnswer := readInput()
	if strings.EqualFold(playerAnswer, selectedQuestion.Answer) {

		reward := 100 * targetDifficulty.Multiplier
		*balance += reward
		fmt.Printf("٩(^ᗜ^ )و ´- Correct! Your balance has been increased by %d. Your new balance is: %d\n", reward, *balance)
		return true
	} else {
		fmt.Println("Incorrect!")

		return false
	}
}

func main() {

	db := QuizDatabase{
		Categories: map[string][]Question{
			"math": {
				{Text: "2 + 2?", Answer: "4", difficulty: "easy"},
				{Text: "5 * 5?", Answer: "25", difficulty: "medium"},
				{Text: "Square root of 144?", Answer: "12", difficulty: "hard"},
			},
			"geography": {
				{Text: "Capital of France?", Answer: "Paris", difficulty: "easy"},
				{Text: "Largest desert?", Answer: "Sahara", difficulty: "medium"},
				{Text: "fastest animal on earth?", Answer: "Cheetah", difficulty: "hard"},
			},
		},
	}

	symbols := map[string]uint{

		"🍒": 4,
		"🍋": 7,
		"🌀": 15,
		"🍊": 12,
		"🍇": 20,
		"🍓": 30}
	multipliers := map[string]uint{
		"🍒": 20,
		"🍋": 10,
		"🍊": 6,
		"🍇": 4,
		"🍓": 2}

	symbolARR := generateSymbolArray(symbols)
	run := 1
	balance := uint(500)
	name := getName()
	fmt.Printf("You have a balance of %d$\n", balance)
	for run == 1 {
		bet := getbet(balance)
		if bet == 0 {
			break
		}
		balance -= bet
		gameResult := spin(symbolARR, 3, 3)
		printSpinResult(gameResult)
		wins := calculateWinnings(gameResult, multipliers)
		if wins == 1 {
			gameResult := spin(symbolARR, 3, 3)
			printSpinResult(gameResult)
			wins = calculateWinnings(gameResult, multipliers)
		}
		if wins > 1 {
			balance += wins * bet
			fmt.Printf("Congratulations! you won %d\n, you're new balance is %d \n", wins*bet, balance)
		} else {
			fmt.Printf("Unlucky!, your new balance is %d\n", balance)
		}
		if balance == 0 {

			if !getQuizz(db, &balance) {
				break
			}

		}
	}

	fmt.Printf("Thank you for playing, %s! Your final balance is %d\n", name, balance)
}
