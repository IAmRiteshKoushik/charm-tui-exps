package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	xstrings "github.com/charmbracelet/x/exp/strings"
)

type Spice int

const (
	Mild Spice = iota + 1
	Medium
	Hot
)

func (s Spice) String() string {
	switch s {
	case Mild:
		return "Mild "
	case Medium:
		return "Medium-Spicy"
	case Hot:
		return "Spicy-Hot"
	default:
		return ""
	}
}

type Order struct {
	Burger       Burger
	Side         string
	Name         string
	Instructions string
	Discount     bool
}

type Burger struct {
	Type     string
	Toppings []string
	Spice    Spice
}

func main() {
	var burger Burger
	order := Order{Burger: burger}

	form := huh.NewForm(

		huh.NewGroup(huh.NewNote().
			Title("Charmburger").
			Description("Welcome to Charmburger.\n\nHow may we take your order?").
			Next(true).
			NextLabel("Next"),
		),

		// Choose a burger
		huh.NewGroup(

			// Field 1: Choose your burger
			huh.NewSelect[string]().
				Options().
				Title("Choose your burger").
				Description("At Charm we truly have a burger for everyone.").
				Options(huh.NewOptions(
					"Charmburger Classic",
					"Chickwich",
					"Fishburger",
					"Charmpossible Burger")...,
				).
				Validate(func(t string) error {
					if t == "Fishburger" {
						return fmt.Errorf("no fish today, sorry")
					}
					return nil
				}).
				Value(&order.Burger.Type),

			// Field 2:
			huh.NewMultiSelect[string]().
				Title("Toppings").
				Description("Choose upto 4.").
				Options(
					huh.NewOption("Lettuce", "Lettuce").Selected(true),
					huh.NewOption("Tomatoes", "Tomatoes").Selected(true),
					huh.NewOption("Charm Sauce", "Charm Sauce"),
					huh.NewOption("Jalapenos", "Jalapenos"),
					huh.NewOption("Cheese", "Cheese"),
					huh.NewOption("Vegan Cheese", "Vegan Cheese"),
					huh.NewOption("Nutella", "Nutella"),
				).
				Validate(func(t []string) error {
					if len(t) <= 0 {
						return fmt.Errorf("at least one topping is required")
					}
					return nil
				}).
				Value(&order.Burger.Toppings).
				Filterable(true).
				Limit(4),
		),

		// Prompt for toppings and special instructions
		huh.NewGroup(
			// Field 1: Spice level
			huh.NewSelect[Spice]().
				Title("Spice Level").
				Options(
					huh.NewOption("Mild", Mild).Selected(true),
					huh.NewOption("Medium", Medium),
					huh.NewOption("Hot", Hot),
				).
				Value(&order.Burger.Spice),

			// Field 2: Choice of sides
			huh.NewSelect[string]().
				Title("Sides").
				Description("You get one free size with this order.").
				Options(huh.NewOptions("Fries", "Disco Fries", "R&B Fries", "Carrots")...).
				Value(&order.Side),
		),

		// Prompt for final details
		huh.NewGroup(

			huh.NewInput().
				Title("What's your name?").
				Placeholder("Margaret Thatcher").
				Description("For when your order is ready.").
				Validate(func(s string) error {
					if s == "Frank" {
						return errors.New("no franks, sorry")
					}
					return nil
				}).
				Value(&order.Name),

			huh.NewText().
				Title("Special Instructions").
				Description("Anything we should know ?").
				Placeholder("Just put it in the mailbox please").
				CharLimit(400).
				Lines(5).
				Value(&order.Instructions),

			huh.NewConfirm().
				Title("Would you like 15% off?").
				Value(&order.Discount).
				Affirmative("Yes!").
				Negative("No."),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	prepareBurger := func() {
		time.Sleep(5 * time.Second)
	}

	_ = spinner.New().Title("Preparing your burger...").Action(prepareBurger).Run()

	{
		// Print order summary
		var sb strings.Builder
		keyword := func(s string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
		}
		fmt.Fprintf(&sb,
			"%s\n\nOne %s%s, topped with %s with %s on the side.",
			lipgloss.NewStyle().Bold(true).Render("BURGER RECEIPT"),
			keyword(order.Burger.Spice.String()),
			keyword(order.Burger.Type),
			keyword(xstrings.EnglishJoin(order.Burger.Toppings, true)),
			keyword(order.Side),
		)

		name := order.Name
		if name != "" {
			name = ", " + name
		}
		fmt.Fprintf(&sb, "\n\nThanks for you order%s!", name)

		if order.Discount {
			fmt.Fprint(&sb, "\n\nEnjoy 15% off .")
		}

		fmt.Println(
			lipgloss.NewStyle().
				Width(40).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Render(sb.String()),
		)
	}
}
