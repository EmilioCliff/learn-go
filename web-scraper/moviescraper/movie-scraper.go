package main

import (
	"sync"
	"encoding/csv"
	"log"
	"fmt"
	"github.com/gocolly/colly"
	"os"
)

type Movie struct{
	Genre string
	Title string
	MovieLink string
	Year string
	ImageURL string
	QualityLevels string
}

var URL string = "https://www.goojara.to/watch-movies-genre"

func main(){
	file, err := os.Create("movies.csv")
	if err != nil{
		fmt.Println("Error while creating file: ", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"Genre",
		"Title",
		"Year",
		"QualityLevels",
		"MovieLink",
		"ImageURL",
	})

	var wg sync.WaitGroup

	genres := []string{"Action", "Adventure", "Comedy", "Drama", "Sci-Fi", "Horror", "Crime", "Thriller", "Romance", "Fantasy"}
	ch := make(chan []Movie, len(genres))

	for _, genre := range genres{
		wg.Add(1)
		go scrapeMovies(&wg, genre, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for movieList := range ch{
		for _, movie := range movieList{
			if err := writer.Write([]string{
				movie.Genre,
				movie.Title,
				movie.Year,
				movie.QualityLevels,
				movie.MovieLink,
				movie.ImageURL,
			}); err != nil{
				fmt.Println("Error writing data to CSV")
			}
		}
	}
}

func scrapeMovies(wg *sync.WaitGroup, genre string, ch chan<- []Movie){
	defer wg.Done()

	c := colly.NewCollector()

	movieList := []Movie{}

	newURL := fmt.Sprintf("%s-%s", URL, genre)

	c.OnHTML("div.dflex", func(e *colly.HTMLElement){

		e.ForEach("div > a", func(_ int, a *colly.HTMLElement){
			var movie Movie
			movie.Genre = genre
			movie.Title = a.ChildText("span.mtl")
			movie.MovieLink = a.Attr("href")
			movie.ImageURL = a.ChildAttr("img", "src")
			movie.QualityLevels = a.ChildText("span.hda")
			movie.Year = a.ChildText("span.hdy")
			movieList = append(movieList, movie)
		})
	})

	c.OnRequest(func(r *colly.Request){
		fmt.Println("Visiting: ", r.URL)
	})

	err := c.Visit(newURL)
	if err != nil{
		log.Fatal(err)
	}
	
	ch <- movieList
}