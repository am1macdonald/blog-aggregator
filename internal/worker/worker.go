package worker

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/am1macdonald/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type Worker struct {
	Limit    int
	Interval time.Duration
	DB       database.Queries
}

func (w *Worker) work(wg *sync.WaitGroup, df *database.Feed) error {
	log.Println("Worker fetching feed")
	defer wg.Done()
	res, err := http.Get(df.Url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("failed to fetch feed")
	}
	if !strings.Contains(res.Header.Get("Content-Type"), "xml") {
		return errors.New("unknown data format")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	rss := Rss{}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return err
	}
	log.Printf("Worker: got feed %v", rss.Channel.Title)
	err = w.processFeed(&rss, df.ID)
	if err != nil {
		return err
	}
	return w.DB.MarkFeedFetched(context.Background(), df.ID)
}

func (w *Worker) processFeed(feed *Rss, feed_id uuid.UUID) error {
	for _, val := range feed.Channel.Item {
		pubDate, err := time.Parse(time.RFC1123Z, val.PubDate)
		if err != nil {
			log.Printf("%v: has invalid date", val.Title)
			continue
		}
		err = w.DB.CreatePost(context.Background(), database.CreatePostParams{
			Title:       val.Title,
			Url:         val.Link,
			Description: val.Description,
			PublishedAt: pubDate,
			FeedID:      feed_id,
		})
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (w *Worker) FetchFeeds() error {
	var wg sync.WaitGroup
	for {
		feeds, err := w.DB.GetNextFeedsToFetch(context.Background(), int32(w.Limit))
		fmt.Println(feeds)
		if err != nil {
			continue
		}
		for _, feed := range feeds {
			wg.Add(1)
			go w.work(&wg, &feed)
		}
		time.Sleep(w.Interval)
	}
}
