
/*
Author: Emilio Santa Cruz
File: main.go
Course: CSC 372
Assignment: Final Project Part 3
Due Date: 12/7/20
TA's: Josh Xiong, Tito Ferra, Martin Marquez, Christian Collberg
Purpose: Serves as the main file for the recipe app. The recipe app takes input
	from the user to determine which recipes to look for and the amount of pages
	to scan for said recipe. In the case the recipe has been searched for before,
	the user will be prompted of a choice to use the previous file. Calls
	recipeNetworker to gather pages of recipes and recipeHandler to handle recipe
	accessing and presenting.
Input: Taken from stdin
Language: Golang 1.15.4
*/

package main

import(
	"fmt"
	"os"
	"strconv"
	"bufio"
	"log"
	"time"
)

/*
Name: check
Parameters: An error object containing with conversion error or file error
Purpose: To end the program with log.Fatal when any error is come across.
Pre-Conditions: e is either nil or an error
Post-Conditions: Either the program ends or nothing happens
*/
func check(e error){
	if e != nil{
		log.Fatal(e)
	}
}

/*
Name: fileExists
Parameters: term: string, represents the search term from the user
Purpose: Checks if the term the user inputted has been searched before.
Pre-Conditions: term is an inputted string
Post-Conditions: True is returned if the file exists, false if not
*/
func fileExists(term string) bool{
	if _, err := os.Stat(term + ".txt"); os.IsNotExist(err){
		return false
	}

	return true
}

/*
Name: formatMultiTerm
Parameters: term: string, represents the search term from the user
Purpose: Checks if any spaces exist in term, if so, they are replaced
	with %20 to allow for URL usage later on. Returns new string.
Pre-Conditions: term is an inputted string
Post-Conditions: newStr is returned
*/
func formatMultiTerm(term string) string{
	var newStr string
	for _, char := range term{
		if char == ' '{
			newStr += "%20"
		} else{
			newStr += string(char)
		}
	}

	return newStr
}

/*
Name: determineMode
Purpose: Takes input from the user to determine which mode the program
	to be in. Takes in the search term and the number of pages from the user.
Pre-Condition: Program starts
Post-Conditions: A string, int, and a boolean is returned
*/
func determineMode() (string, int, bool){
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter a search term for recipes.")
	for scanner.Scan(){
		if len(scanner.Text()) < 1{
			fmt.Println("Enter a term of length >0")
			continue
		} 
		break
	}
	term := scanner.Text()
	if fileExists(term){
		fmt.Println(term + ".txt found. Would you like to use this file?(y/n)")
		for scanner.Scan(){
			if scanner.Text() != "y" && scanner.Text() != "n" {
				fmt.Println("Enter y or n for yes or no.")
				continue
			}
			break
		}
		if scanner.Text() == "y"{
			return term, 0, true	// returns early due to reuse of file
		}
	}

	fmt.Println("How many pages to process? Enter an integer.")
	var num int
	var err error
	for scanner.Scan(){
		num, err = strconv.Atoi(scanner.Text())
		if err != nil{
			fmt.Println("Enter an integer.")
			continue
		} else if num < 1 {
			fmt.Println("Enter an integer > 0")
			continue
		}
		break
	}
	return term, num, false
}

func main() {
	var recipeMap map[string][]recipe
	term, pageNum, mode := determineMode()
	if mode{	// reuse
		recipeMap = loadRecipes(term)
	} else{		// new file
		recipeMap = getPages(formatMultiTerm(term), pageNum)
		if len(recipeMap) == 0{
			fmt.Print("No recipes found, try again with new restart.")
			time.Sleep(1 * time.Second)
			fmt.Print(".")
			time.Sleep(1 * time.Second)
			fmt.Print(".")
			time.Sleep(1 * time.Second)
			fmt.Print(".")
			time.Sleep(1 * time.Second)
			fmt.Print(".")
			time.Sleep(1 * time.Second)
			fmt.Print(".\n")
			os.Exit(1)
		}
		saveRecipes(recipeMap, term)
	}
	presentRecipes(recipeMap)
}