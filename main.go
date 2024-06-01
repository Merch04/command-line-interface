package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"golang.org/x/term"
)

type PageItem struct {
	name string
	link *Page
}
type Page struct {
	content string
}
type PageController struct {
	currentItem  uint
	selectedItem uint

	itemsArr *[]PageItem
	count    uint
}

func (c *PageController) increase() {

	if c.selectedItem >= c.count-1 {
		c.selectedItem = 0
	} else {
		c.selectedItem++
	}
}

func (c *PageController) reduce() {

	if c.selectedItem == 0 {
		c.selectedItem = c.count - 1
	} else {
		c.selectedItem--
	}
}

func main() {

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	//Создаем контроллер
	pageController, err := initStructures()
	if err != nil {
		panic(err)
	}

	//читаем нажатия
	for {
		drawPage(pageController)
		inputCheck(pageController)
	}

}

func initStructures() (*PageController, error) {

	items := []PageItem{
		{"Главная", &Page{content: "Это главная страница"}},
		{"Инфо", &Page{content: "Это информация"}},
		{"Выход", &Page{content: "Это типа выход"}},
	}

	controller := PageController{currentItem: 0, itemsArr: &items, count: uint(len(items))}
	return &controller, nil
}

func inputCheck(control *PageController) {

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC, keys.Escape:
			os.Exit(0)
		case keys.Up:
			control.reduce()
		case keys.Down:
			control.increase()
		case keys.Enter:
			if control.selectedItem == control.count-1 { //последний элемент, Item - Выход
				os.Exit(0)
			}
			control.currentItem = control.selectedItem
		case keys.RuneKey:
			if strings.ToLower(key.String()) == "w" || strings.ToLower(key.String()) == "ц" {
				control.reduce()
			} else if strings.ToLower(key.String()) == "s" || strings.ToLower(key.String()) == "ы" {
				control.increase()
			}
		default:
			return false, nil // Return false to continue listening
		}
		return true, nil
	})

}

func clearDisplay() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

}
func drawPage(control *PageController) {

	clearDisplay()

	items := *control.itemsArr
	for index, item := range items {
		if index == int(control.currentItem) {
			fmt.Print("*")
		}
		if index == int(control.selectedItem) {
			fmt.Printf("%d. %s  <----\n", index, item.name)
		} else {
			fmt.Printf("%d. %s \n", index, item.name)
		}
	}
	fmt.Print("\n\n\n")
	fmt.Printf("Content: %s\n", items[control.currentItem].link.content)
}
