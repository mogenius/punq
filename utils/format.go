package utils

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/jaevor/go-nanoid"
	"github.com/mogenius/punq/logger"
)

var FirstNameList = []string{
	"Liam", "Emma", "Noah", "Olivia", "Ava",
	"Sophia", "Jackson", "Aiden", "Lucas", "Muhammad",
	"Amelia", "Mateo", "Ethan", "Harper", "Evelyn",
	"Mia", "Ella", "Riley", "Aria", "Logan",
	"Zoe", "Benjamin", "Oliver", "Lily", "Leo",
	"Charlotte", "Mason", "Isabella", "Layla", "Isaac",
	"Mila", "Sophie", "Elijah", "Emily", "Daniel",
	"James", "Aiden", "Abigail", "Levi", "Chloe",
	"Henry", "Alexander", "Sebastian", "Jack", "Hannah",
	"Jayden", "Gabriel", "Matthew", "Alice", "Oscar",
	"Josiah", "Evie", "Theo", "Isla", "Jaxon",
	"Grace", "Eva", "Samuel", "Owen", "Victoria",
	"Joseph", "Zachary", "Violet", "John", "William",
	"Ezra", "Ellie", "Freya", "Dylan", "Penelope",
	"Michael", "Scarlett", "Luna", "Max", "Alyssa",
	"Isabelle", "Eliza", "Luca", "Thomas", "Poppy",
	"David", "Ruby", "Christopher", "Jade", "Rose",
	"Sienna", "George", "Harvey", "Kaylee", "Annie",
	"Nathan", "Madison", "Jacob", "Noelle", "Parker",
	"Sarah", "Evelina", "Leo", "Ruby", "Abigail",
}

var MiddleNamesList = []string{
	"Bumblebee", "Rainbow", "Whiz", "Jolly", "Bubbles",
	"Sparkle", "Noodle", "Waffle", "Pickle", "Jiggle",
	"Twinkle", "Giggle", "Fizzle", "Muffin", "Pumpkin",
	"Squiggle", "Tofu", "Jazz", "Fizz", "Sunny",
	"Fluffy", "Peanut", "Jellybean", "Snicker", "Ripple",
	"Glimmer", "Cupcake", "Pudding", "Tinker", "Pebble",
	"Cuddle", "Bumpkin", "Dizzy", "Lolly", "Nugget",
	"Twirl", "Fizzypop", "Wiggles", "Snuggles", "Squishy",
	"Blinky", "Bubblegum", "Frodo", "Sizzle", "Taco",
	"Smiley", "Snickerdoodle", "Wobble", "Popsicle", "Zigzag",
	"Sprinkles", "Doodle", "Pizzazz", "Quicksilver", "Razzmatazz",
	"Duckling", "Hiccup", "Pumpernickel", "Zoodle", "Quizzical",
	"Flitter", "Whisper", "Mustard", "Wacky", "Scooter",
	"Moose", "Tizzy", "Bamboo", "Zephyr", "Rolo",
	"Sniffle", "Gobble", "Beep", "Cobweb", "Twizzle",
	"Bizz", "Fuddle", "Puzzle", "Rumble", "Rover",
	"Squabble", "Tumbleweed", "Vroom", "Whizzle", "YoYo",
}

var LastNamesList = []string{
	"Smith", "Kim", "Johnson", "Lee", "Brown",
	"Patel", "Garcia", "Rodriguez", "Martinez", "Chen",
	"Jones", "Nguyen", "Williams", "Lopez", "Gonzalez",
	"Perez", "Hernandez", "Tanaka", "Silva", "Santos",
	"Cohen", "Kumar", "Wang", "Meyer", "Schneider",
	"Taylor", "Anderson", "White", "Young", "Harris",
	"Clark", "Lewis", "Turner", "Walker", "Hall",
	"Allen", "Roberts", "Wright", "King", "Hill",
	"Scott", "Green", "Baker", "Adams", "Nelson",
	"Campbell", "Mitchell", "Robinson", "Carter", "Thomas",
	"Mueller", "Fernandez", "Oliveira", "Sharma", "Singh",
	"Liu", "Lin", "Ali", "Khan", "Jackson",
	"Parker", "Phillips", "Davis", "Murphy", "Price",
	"Suzuki", "Ross", "Reyes", "Jenkins", "Morris",
	"Sanchez", "Perry", "Powell", "Russell", "Moore",
	"Ramirez", "Gray", "James", "Watson", "Brooks",
	"Kelly", "Sanders", "Foster", "Evans", "Barnes",
}

func RandomFirstName() string {
	return FirstNameList[RandomInt(0, len(FirstNameList))]
}

func RandomMiddleName() string {
	return MiddleNamesList[RandomInt(0, len(MiddleNamesList))]
}

func RandomLastName() string {
	return LastNamesList[RandomInt(0, len(LastNamesList))]
}

func RandomInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func FatalError(message string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf(red("Error: %s\n"), message)
	os.Exit(0)
}

func PrintError(message string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Println(red(message))
}

func PrintInfo(message string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println(yellow(message))
}

func StatusEmoji(works bool) string {
	if works {
		return "✅"
	}
	return "❌"
}

func NanoId() string {
	id, err := nanoid.Custom("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 21)
	if err != nil {
		logger.Log.Error(err)
	}
	return id()
}

func NanoIdSmallLowerCase() string {
	id, err := nanoid.Custom("abcdefghijklmnopqrstuvwxyz1234567890", 10)
	if err != nil {
		logger.Log.Error(err)
	}
	return id()
}

func NanoIdExtraLong() string {
	id, err := nanoid.Custom("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 21)
	if err != nil {
		logger.Log.Error(err)
	}
	return id()
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func QuickHash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprint(h.Sum32())
}

func BytesToHumanReadable(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func NumberToHumanReadable(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %c",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func JsonStringToHumanDuration(jsonTime string) string {
	parsedTime, err := time.Parse(time.RFC3339, jsonTime)
	if err != nil {
		fmt.Printf("Error parsing date '%s'. Must be RFC3339 conform.", jsonTime)
	}
	return HumanDuration(time.Since(parsedTime))
}

// TAKEN FROM Kubernetes apimachineryv0.25.1
func HumanDuration(d time.Duration) string {
	// Allow deviation no more than 2 seconds(excluded) to tolerate machine time
	// inconsistence, it can be considered as almost now.
	if seconds := int(d.Seconds()); seconds < -1 {
		return "<invalid>"
	} else if seconds < 0 {
		return "0s"
	} else if seconds < 60*2 {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := int(d / time.Minute)
	if minutes < 10 {
		s := int(d/time.Second) % 60
		if s == 0 {
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%dm%ds", minutes, s)
	} else if minutes < 60*3 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := int(d / time.Hour)
	if hours < 8 {
		m := int(d/time.Minute) % 60
		if m == 0 {
			return fmt.Sprintf("%dh", hours)
		}
		return fmt.Sprintf("%dh%dm", hours, m)
	} else if hours < 48 {
		return fmt.Sprintf("%dh", hours)
	} else if hours < 24*8 {
		h := hours % 24
		if h == 0 {
			return fmt.Sprintf("%dd", hours/24)
		}
		return fmt.Sprintf("%dd%dh", hours/24, h)
	} else if hours < 24*365*2 {
		return fmt.Sprintf("%dd", hours/24)
	} else if hours < 24*365*8 {
		dy := int(hours/24) % 365
		if dy == 0 {
			return fmt.Sprintf("%dy", hours/24/365)
		}
		return fmt.Sprintf("%dy%dd", hours/24/365, dy)
	}
	return fmt.Sprintf("%dy", int(hours/24/365))
}

func MergeMaps(maps ...map[string]string) map[string]string {
	resultMap := make(map[string]string)

	// Iterate over the slice of maps
	for _, m := range maps {
		// Add all elements from each map, potentially overwriting
		for key, value := range m {
			resultMap[key] = value
		}
	}
	return resultMap
}

func PrintLogo() {
	logo := `   ______  __    __ _______   ______  
  /      \|  \  |  \       \ /      \ 
  |  ▓▓▓▓▓▓\ ▓▓  | ▓▓ ▓▓▓▓▓▓▓\  ▓▓▓▓▓▓\
  | ▓▓  | ▓▓ ▓▓  | ▓▓ ▓▓  | ▓▓ ▓▓  | ▓▓
  | ▓▓__/ ▓▓ ▓▓__/ ▓▓ ▓▓  | ▓▓ ▓▓__| ▓▓
  | ▓▓    ▓▓\▓▓    ▓▓ ▓▓  | ▓▓\▓▓    ▓▓
  | ▓▓▓▓▓▓▓  \▓▓▓▓▓▓ \▓▓   \▓▓ \▓▓▓▓▓▓▓
  | ▓▓                             | ▓▓
  | ▓▓                             | ▓▓
   \▓▓                              \▓▓`
	fmt.Print(logo)
}
