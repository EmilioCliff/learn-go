package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type Movie struct {
	Genre         string
	Title         string
	MovieLink     string
	Year          string
	ImageURL      string
	QualityLevels string
}

const (
	URL         = "https://www.goojara.to/watch-movies-genre"
	scraperType = "without"
)

var elapsedTime time.Duration

func main() {
	file, err := os.Create("movies.csv")
	if err != nil {
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

	genres := []string{"Action", "Adventure", "Comedy", "Drama", "Sci-Fi", "Horror", "Crime", "Thriller", "Romance", "Fantasy"}

	var wg sync.WaitGroup
	ch := make(chan []Movie, len(genres))

	startTime := time.Now()

	switch scraperType {
	case "with":
		for _, genre := range genres {
			wg.Add(1)
			go scrapeMoviesWithGoroutines(&wg, genre, ch)
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		elapsedTime = time.Since(startTime)

		for movieList := range ch {
			writeMoviesToCSV(writer, movieList)
		}

	case "without":
		var allMovieList [][]Movie
		for _, genre := range genres {
			returnMovieList := scrapeMoviesWithoutGoroutines(genre)
			allMovieList = append(allMovieList, returnMovieList)
		}

		elapsedTime = time.Since(startTime)

		for _, movieList := range allMovieList {
			writeMoviesToCSV(writer, movieList)
		}

	default:
		fmt.Println("Invalid input. Please enter 'with' or 'without'.")
	}

	fmt.Printf("Time Taken: %v\n", elapsedTime)

}

func writeMoviesToCSV(writer *csv.Writer, movieList []Movie) {
	for _, movie := range movieList {
		if err := writer.Write([]string{
			movie.Genre,
			movie.Title,
			movie.Year,
			movie.QualityLevels,
			movie.MovieLink,
			movie.ImageURL,
		}); err != nil {
			fmt.Println("Error writing data to CSV")
		}
	}

}

func scrapeMoviesWithGoroutines(wg *sync.WaitGroup, genre string, ch chan<- []Movie) {
	defer wg.Done()

	c := colly.NewCollector()

	movieList := []Movie{}

	newURL := fmt.Sprintf("%s-%s", URL, genre)

	c.OnHTML("div.dflex", func(e *colly.HTMLElement) {

		e.ForEach("div > a", func(_ int, a *colly.HTMLElement) {
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

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	err := c.Visit(newURL)
	if err != nil {
		log.Fatal(err)
	}

	ch <- movieList
}

func scrapeMoviesWithoutGoroutines(genre string) []Movie {
	c := colly.NewCollector()

	movieList := []Movie{}

	newURL := fmt.Sprintf("%s-%s", URL, genre)

	c.OnHTML("div.dflex", func(e *colly.HTMLElement) {

		e.ForEach("div > a", func(_ int, a *colly.HTMLElement) {
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

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	err := c.Visit(newURL)
	if err != nil {
		log.Fatal(err)
	}

	return movieList
}
