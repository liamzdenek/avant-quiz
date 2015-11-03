package main

import (
	"fmt"
	"github.com/mrjones/oauth"
	"bufio"
    "encoding/json"
    "time"
    "strings"
)

type Tweet struct {
    Text string `json:"text"`
}

var DEFAULT_FILTER = []string{
    "and",
    "the",
    "me",
    "i",
    "you",
}

func main() {
    done := make(chan struct{});

    stream := GetRawStream(done);

    parsedStream := ParseStream(stream);

    <-TweetStats(done, parsedStream, DEFAULT_FILTER, time.Minute * 5)
}

type internalTweetStats struct {
    Filter []string
    NumTweets uint64
    WordCount uint64
    PopularWords map[string]uint64
}

func(ts *internalTweetStats) Push(tweet *Tweet) {
    fmt.Printf("Got tweet: %s\n", tweet.Text);

    words := strings.Split(tweet.Text, " ");

    ts.NumTweets++;
    //ts.WordCount = ts.WordCount + uint64(len(words));

    OUTER: for _, word := range words {
        lc_word := strings.ToLower(word);
        lc_word = strings.Trim(lc_word, " ,.?!@#$%^&*()_+-=|\\}]{['\";:<>/`~");
        if lc_word == "" {
            continue OUTER;
        }
        for _, filter_word := range ts.Filter {
            lc_filter_word := strings.ToLower(filter_word);
            if(lc_word == lc_filter_word) {
                continue OUTER;
            }
        }
        ts.WordCount++;
        ts.PopularWords[lc_word]++;
    }
}

func(ts *internalTweetStats) PrintStats() {
    type WordLinkedList struct {
        Word string
        Count uint64
        Next *WordLinkedList
    }
    
    fmt.Printf("=== Stats ===\n");
    fmt.Printf("      Word Count: %d\n", ts.WordCount);
    fmt.Printf("Number of Tweets: %d\n", ts.NumTweets);

    var master *WordLinkedList = nil

    OUTER: for word, count := range ts.PopularWords {
        head := master;
        for {
            // head will never be null... hopefully
            // if the linked list is empty or the current head is the head of the linked list and we're greater than that node
            if head == nil || (head == master && head.Count < count) {
                master = &WordLinkedList{
                    Word: word,
                    Count: count,
                };
                continue OUTER;
            // else if the next item in the linked list is not empty
            } else if head.Next != nil {
                // if the current item is greater than us but the next item is less than us
                if(head.Next.Count <= count && head.Count > count) {
                    head.Next = &WordLinkedList{
                        Word: word,
                        Count: count,
                        Next: head.Next,
                    }
                    continue OUTER;
                // else advance the head
                } else {
                    head = head.Next;
                }
            // else we're at the tail of the linked list, put this at the end
            } else {
                head.Next = &WordLinkedList{
                    Word: word,
                    Count: count,
                }
                continue OUTER;
            }
        }
    }

    fmt.Printf("=== Popular Words ===\n");
    head := master;
    for i := 0; i < 10; i++ {
        fmt.Printf("\t%d => %s\n",head.Count, head.Word);
        if head.Next != nil {
            head = head.Next;
        } else {
            break;
        }
    }
}

func TweetStats(sourcedone chan struct{}, tweetstream chan *Tweet, word_filter []string, duration time.Duration)  chan struct{} {
    selfdone := make(chan struct{});

    go func() {
        defer close(selfdone);

        stats := &internalTweetStats{
            Filter: word_filter,
            PopularWords: map[string]uint64{},
        };
        defer stats.PrintStats();

        after := time.After(duration);

        OUTER: for {
            select {
            case <-after:
                fmt.Printf("Time has elapsed\n");
                close(sourcedone);
            case tweet, ok := <-tweetstream:
                if tweet != nil {
                    stats.Push(tweet);
                }
                if(!ok) {
                    break OUTER;
                }
            }
        }
    }()

    return selfdone;
}

func ParseStream(rawstream chan []byte) chan *Tweet {
    tweets := make(chan *Tweet);

    go func() {
        defer close(tweets)
        for raw := range rawstream {
            new_tweet := &Tweet{};
            err := json.Unmarshal(raw, new_tweet);
            if err != nil {
                fmt.Printf("An error occured while parsing a tweet. Silencing. %s\n", err);
            } else if len(new_tweet.Text) > 0 {
                tweets <- new_tweet;
            }
        }
    }()

    return tweets;
}

func GetRawStream(done chan struct{}) (chan []byte) {
    stream := make(chan []byte);

    go func() {
        defer close(stream);

        consumer := oauth.NewConsumer(
            "fUMzljWf3zdg40R6BjdaIdxam",
            "V7FjRIjHUcbeaeGeEH94CHooSWsNCPZr5S1V8lEGIKUwt4brIs",
            oauth.ServiceProvider{},
        )
        accessToken := &oauth.AccessToken{
            Token:  "1460496416-2vcYZZ2qT8uGpOuCwxfmoFA3OlTI5xoO1LOY1PW",
            Secret: "XrxntnaLHc5wcQXxJFRWQr2RNnRFMMUfKG9uP6vK7nFpD",
        }

        url := "https://stream.twitter.com/1.1/statuses/sample.json"

        response, err := consumer.Get(url, nil, accessToken)
        Check(err)

        r := bufio.NewReader(response.Body)

        OUTER: for {
            line, err := r.ReadBytes('\n')
            Check(err)
            select {
            case <-done:
                break OUTER;
            case stream <- line:
            }
        }
    }()

    return stream;
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
