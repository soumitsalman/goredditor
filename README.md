This is an open source golang library for using https://www.reddit.com/dev/api/

**Current functionalities**
- Authenticate and re-authenticate: Requires going through basic auth first
- Retrieving data of "me"
- Retrieving `hot`, `top`, `best` posts from the users landing page and also for a specified `subreddit`
- Retrieving the list of `subreddit`s that the user is already subscribing to
- Retrieving the list of `subreddit`s that are similar to the `subreddit`s the user is subscribed to. This does NOT retrun a unique list and may contain duplicates
- Subscribing to a new 'subreddit`
- Creating a new post with either free-form markdown text or link in a given `subreddit`
- Commenting on a given post or a comment
- 

**Known gaps**
- There is no input check or validation
- There is no output check or error handling either
- There is no unit test

**Upcoming future improvements**
- [x] Adding search functionality `subreddit`s
- [x] Adding search functionality to retrieve posts
- [ ] Stabilizing the code with better error handling and input validation
- [x] Retrieving comments of a given thread
- [ ] Publish the api-s as golang libraries for everyone to use

**Contact**
If you want to contribute to the code base please feel free to do so. If you have any question feel free to reachout to me
