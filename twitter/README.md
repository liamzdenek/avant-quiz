There were a few decisions made in this that affect the behavior of the program:

* this is very eurocentric in terms of word count. I didn't want to go through the complexity of which unicode characters map to exactly what a "word" is, so word count is simply by space.

* the tweet wordcount stats does not scrub itself over time, and preparing the output data does not scrub results with very low word counts, so a long-lived scan will definitely run OOM while preparing its output.

* the word count considers hashtags as words.

* similarly, the word count considers emojis as characters, and an emoji surrounded by whitespace will be considered a word.

* this program handles EOF silently, so it may very well just abort the twitter data stream if twitter closes the socket on us -- which the docs say to expect and program for, and have a huge document for handling properly. I think this is sufficient for now.

* json that failed to decode is silently thrown away.


Sample Output
```
(a ton of Got Tweet: lines have been scrubbed for brevity)
Got tweet: I don't understand people who come ON Twitter, to tell people ON Twitter how they're doing nothing ON TWITTER!                                                 
=== Stats ===                                                                        
      Word Count: 16367                                                              
Number of Tweets: 1963                                                               
=== Popular Words ===                                                                
        744 => rt                                                                    
        44 => this                                                                   
        36 => i'm                                                                    
        23 => en                                                                     
        20 => aldubthevisitor                                                        
        20 => but                                                                    
        20 => what                                                                   
        18 => all                                                                    
        17 => out                                                                    
        16 => more                                                                   
```

Optional Part B) How would you implement it so that if you had to stop the program and restart,
 it could pick up from the total word counts that you started from?

 I would implement this by saving the internalTweetStats object, and reloading it upon resume. This would handle all of the aforementioned. 
