/*
Author: Emilio Santa Cruz
File: recipeNetworker.go
Course: CSC 372
Assignment: Final Project Part 3
Due Date: 12/7/20
TA's: Josh Xiong, Tito Ferra, Martin Marquez, Christian Collberg
Purpose: Serves as the file that deals with accessing and parsing relevant
	information from the web.
Input: Taken from stdin
Language: Golang 1.15.4
*/

package main

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"
	"mvdan.cc/xurls/v2"
	"fmt"
)

/*
Name: recipe
Fields:
	link - holds the link to the recipe the struct is representing
	name - holds the name of the recipe
	ratingCount - holds the number of rating the recipe has received
Purpose: Represents a cluster of points and their centroid
*/
type recipe struct{
	link string
	name string
	ratingCount string
}

/*
Name: getRecipeName
Paramters: string, line is a line from the HTML from a webpage
Purpose: Gets the recipe name from a line picked out from a webpage
Pre-Conditions: A webpage has been accessed
Post-Conditions: The name of the recipe within line is returned
*/
func getRecipeName(line string) (name string){
	strptr := len(line) - 2

	for line[strptr] != '/'{
		name = string(line[strptr]) + name
		strptr--
	}

	return name
}

/*
Name: parseRating
Parameters: string, line is a line from the HTML from a webpage
Purpose: Gets the rating number from a line picked out from the webpage
Pre-Conditions: A webpage has been accessed.
Post-Conditions: The rating of the recipe within line is returned
*/
func parseRating(line string) string{
	ratingIndex := strings.Index(line, "ratingstars=")
	ratingIndex = ratingIndex + 13
	start := ratingIndex
	for line[ratingIndex] != '"'{
		ratingIndex++
	}
	rating, err := strconv.ParseFloat(line[start:ratingIndex], 32)
	check(err)

	return fmt.Sprintf("%.1f", rating)
}

/*
Name: parseRatingCount
Parameters: string, line is a line from the HTML from a webpage
Purpose: Gets the rating amount from a line picked out from the webpage
Pre-Conditions: A webpage has been accessed.
Post-Conditons: The rating count of the recipe within line is returned
*/
func parseRatingCount(line string) string{
	ratingStr := ""
	rateCountIndex := strings.Index(line, ">") + 1
	var stopChar byte
	stopChar = '<'
	if line[rateCountIndex] == stopChar{	// handles rating count >999
		rateCountIndex += 29
		stopChar = '"'
	}
	for line[rateCountIndex] != stopChar{
		ratingStr += string(line[rateCountIndex])
		rateCountIndex++
	}

	return ratingStr
}

/*
Name: processPage
Parameters: page is a page from the web, term is a string from the user,
	pageNum is an int representing the number of pages to process, ratingsMap
	is a map holding the recipes from the web pages.
Purpose: To get all relevant information from page. Gets ratings, name, and
	rating counts, puts them in recipe structs, and stores them in ratingsMap
Pre-Conditions: A new page has been accessed and downloaded
Post-Conditions: ratingsMap is updated with recipes from the webpage
*/
func processPage(page *http.Response, term string, pageNum int, 
	ratingsMap map[string][]recipe){
	currStr := ""
	rxRelaxed := xurls.Relaxed()
	linkNotFound := true
	scanner := bufio.NewScanner(page.Body)
	var newRecipe recipe
	var rating string
	for scanner.Scan(){
		line := scanner.Text()
		// gets recipe name
		if linkNotFound && strings.Contains(line, "https://www.allrecipes.com/recipe/"){ 
			currStr = rxRelaxed.FindString(line)
			newRecipe = recipe{link:currStr, name:getRecipeName(currStr)}
			linkNotFound = false
		}
		// gets rating
		if strings.Contains(line, "data-ratingstars="){	
			rating = parseRating(line)
		}
		// gets rating counts and adds to ratingsMap
		if strings.Contains(line, "fixed-recipe-card__reviews"){ 
			newRecipe.ratingCount = parseRatingCount(line)
			if _, ok := ratingsMap[rating]; !ok{	// if rating does not exist
				ratingsMap[rating] = []recipe{newRecipe}
			} else{									// if rating exists
				ratingsMap[rating] = append(ratingsMap[rating], newRecipe)
			}
			linkNotFound = true
		}
	}
	fmt.Printf("%d page(s) processed\n", pageNum - 1)
	page.Body.Close()
}

/*
Name: getPages
Parameters: term is a string from the user, pageCount is an int from the user
Purpose: Processes pageNum amount of pages and puts their relevant contents
	into ratingsMap. Uses goroutines to process pages while the next page is downloaded
Pre-Conditions: The mode of the program is to grab new recipes rather than
	use an old file.
Post-Conditions: Returns ratingsMap which contains the recipes from the webpages
*/
func getPages(term string, pageCount int) (map[string][]recipe){
	pageNum := 1
	page, err := http.Get("https://www.allrecipes.com/search/results/?wt=" + 
	term + "&sort=re&page=" + strconv.Itoa(pageNum))
	check(err)
	ratingsMap := make(map[string][]recipe)

	for pageNum <= pageCount{
		pageNum++
		go processPage(page, term, pageNum, ratingsMap)
		page, err = http.Get("https://www.allrecipes.com/search/results/?wt=" + 
		term + "&sort=p&page=" + strconv.Itoa(pageNum))
		check(err)
	}

	return ratingsMap
}