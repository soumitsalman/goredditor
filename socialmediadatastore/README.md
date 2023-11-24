This is the **social media data store** as a service for brandy.io. This provides the following functions
* Storing/caching social media contents such as posts, comments, channels, subreddits
* Message queuing for social media actions such as
    * newly discovered posts, subreddits, comments
    * suggested user actions like joining a subreddit, making a post or comment (this part is yet to be implemented)
Currently this is primarily used by goredditor scraper

There is currently no test code and I am using the main.go for manually testing this out


