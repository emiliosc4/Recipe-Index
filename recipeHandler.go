
/*
Author: Emilio Santa Cruz
File: recipeHandler.go
Course: CSC 372
Assignment: Final Project Part 3
Instructor: Dr. McCann
Due Date: 12/7/20
TA's: Josh Xiong, Tito Ferra, Martin Marquez, Christian Collberg
Purpose: Handles gettings recipes from a file, storing recipes to a file,
	and creating a loop for the user to pick out recipes from the information
	they are presented with.
Input: Taken from stdin
Language: Golang 1.15.4
*/

package main

import(
	"bufio"
	"os"	
	"strconv"
	"sort"
	"fmt"
	"github.com/pkg/browser"
)

/*
Name: saveRecipes
Parameters: ratingsMap is a map containing ratings and a list of recipe
	structs tied to that rating, term is a string from the user
Purpose: Saves the contents of ratingsMap to a text file
Pre-Conditions: ratingsMap is non-empty
Post-Conditions: A text file is written
*/
func saveRecipes(ratingsMap map[string][]recipe, term string){
	keys := make([]string, len(ratingsMap))
	file, err := os.Create(term + ".txt")
	check(err)
	defer file.Close()

	i := 0
	for key, recipeList := range ratingsMap{
		// add length for each rating to allow for easy parsing
		file.WriteString(key + "\n")
		for _, recipe := range(recipeList){
			file.WriteString(recipe.name + "\n" + recipe.ratingCount + "\n" + recipe.link + "\n")
		}
		file.WriteString("\n")
		keys[i] = key
		i++
	}
}

/*
Name: loadRecipes
Parameters: searchTerm is a string from the user.
Purpose: To read in a text file specified buy the user and collect the
	recipes contained.
Pre-Conditions: The user has chosen to reuse a file
Post-Conditions: The text file specified by the user is read in
*/
func loadRecipes(searchTerm string) (map[string][]recipe){
	ratingsMap := make(map[string][]recipe)
	file, err := os.Open(searchTerm + ".txt")
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var rating string
	readRating := true
	readName := true
	readCount := true
	var currRecipe recipe
	for scanner.Scan(){
		if readRating{	// reads rating of recipes
			rating = scanner.Text()[0:3]
			var recipeList []recipe
			ratingsMap[rating] = recipeList
			readRating = false
		} else if scanner.Text() == ""{
			readRating = true
		} else if readName{	// reads name of recipe
			currRecipe = recipe{name:scanner.Text()}
			readName = false
		} else if readCount{// reads count of rating
			currRecipe.ratingCount = scanner.Text()
			readCount = false
		} else{				// reads link of recipe
			currRecipe.link = scanner.Text()
			ratingsMap[rating] = append(ratingsMap[rating], currRecipe)
			readName = true
			readCount = true
		}
	}

	return ratingsMap
}

/*
Name: getCurrRecipeList
Parameters: scanner is a bufio.scanner, keys is a list of string holds keys to
	recipeMap, recipeMap is a map of ratings and list of recipe structs.
Purpose: To take input from the user to determine a rating to read in
Pre-Conditions: The program is in recipe access mode
Post-Conditions: A list of recipe structs is returned
*/
func getCurrRecipeList(scanner *bufio.Scanner, keys []string, recipeMap map[string][]recipe) []recipe{
	fmt.Println("Which rating would you like to look at?")
	for _, key := range keys{	// prints out ratings(keys)
		fmt.Println(key + ": " + strconv.Itoa(len(recipeMap[key])) + " recipes")
	}
	
	for scanner.Scan(){			// determines which rating
		if _, ok := recipeMap[scanner.Text()]; !ok{
			fmt.Println(scanner.Text() + " not found. Enter one rating above.")
		} else{
			break
		}
	}

	fmt.Println("Enter number next to recipe name to open that recipe.\n" + 
	"To look a different rating, enter diff.\nEnter exit to end execution.")
	for i, currRecipe := range recipeMap[scanner.Text()]{
		fmt.Print(i)
		fmt.Println("\t" + currRecipe.name + ": " + currRecipe.ratingCount + " ratings")
	}

	return recipeMap[scanner.Text()]
}

/*
Name: presentRecipes
Parameters: recipeMap is a map contains ratings and their attached list of
	recipe structs
Purpose: Presents the recipes collected and allows the user to pick rating
	of recipes and open recipes URL upon command.
Pre-Conditions: The recipes have been read in from the web
Post-Conditions: The user is presented with a loop in which they are allowed
	to pick recipes depending on rating and name. Allows user to swap ratings
	for readability and the ability to exit.
*/
func presentRecipes(recipeMap map[string][]recipe){
	var keys []string
	for rating := range recipeMap{
		keys = append(keys, rating)
	}
	sort.Strings(keys)
	
	scanner := bufio.NewScanner(os.Stdin)
	currRecipeList := getCurrRecipeList(scanner, keys, recipeMap)
	for scanner.Scan(){
		if scanner.Text() == "exit"{
			break
		} else if scanner.Text() == "diff"{
			currRecipeList = getCurrRecipeList(scanner, keys, recipeMap)
		} else{
			currI, err := strconv.Atoi(scanner.Text())
			if(err != nil || currI < 0 || currI >= len(currRecipeList)){
				fmt.Println("Enter only numbers that are listed.")
			} else{
				browser.OpenURL(currRecipeList[currI].link)
			}
		}		
	}
}