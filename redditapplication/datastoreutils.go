package redditapplication

import (
	"fmt"
	"strings"

	go_linq "github.com/ahmetb/go-linq/v3"
	ds "github.com/soumitsalman/goredditor/socialmediadatastore"
	"github.com/soumitsalman/goredditor/utils"
)

func initializeDataStores() {
	ds.InitializeContentStoreClient()
	ds.InitializeProcessingQueues()
}

// TODO: known bug - this does not update the stats if they have changed. so there should be some kind of fix to keep the store up to date
func saveNewItemsToDB(user_id string, item_kind string, items []RedditData) {
	items = filterNewItems(user_id, items)
	cs_items, pq_items := getNormalizedData(user_id, item_kind, items)
	ds.AddToContentStore(cs_items)
	ds.AddToUserActionStore(pq_items)
	ds.BatchQue(ds.NEW, pq_items)
}

func filterNewItems(user_id string, items []RedditData) []RedditData {
	var new_items []RedditData
	existing_ids := ds.GetExistingUserActionsContentIds(user_id, REDDIT_PRETTY_URL)

	go_linq.From(items).
		Where(func(i interface{}) bool {
			return !go_linq.From(existing_ids).
				Contains(i.(RedditData).Name)
		}).ToSlice(&new_items)

	// TODO: delete this part. this is debug only
	fmt.Println(existing_ids)
	go_linq.From(items).ForEach(func(i interface{}) {
		fmt.Printf("[%s] %s\n", i.(RedditData).Name, i.(RedditData).Title)
	})

	return new_items
}

func getNormalizedData(user_id string, item_kind string, items []RedditData) ([]ds.ContentStoreData, []ds.UserActionData) {
	ds_items := make([]ds.ContentStoreData, len(items))
	pq_items := make([]ds.UserActionData, len(items))

	for i, v := range items {
		// transform for content store
		ds_items[i] = ds.ContentStoreData{
			//applies to all
			Id:          v.Name,
			UrlId:       v.Id,
			Title:       v.Title,
			Kind:        v.Kind,
			CreatedDate: v.CreatedDate,

			//applies to post and comment
			Channel:        v.Subreddit,
			Author:         v.Author,
			Score:          v.Score,
			UpvoteRatio:    v.UpvoteRatio,
			NumComments:    v.NumComments,
			Category:       v.PostCategory,
			NumSubscribers: v.SubredditSubscribers,
			Url:            ensureRedditDotCom(v.Link),
		}
		// special field overrides for content store data
		switch v.Kind {
		case SUBREDDIT:
			// overriding channel name, category, numsubscribers and url
			ds_items[i].Channel = v.DisplayName
			ds_items[i].Text = utils.ExtractTextFromHtml(v.PublicDescription) + "\n" + utils.ExtractTextFromHtml(v.Description)
			ds_items[i].Category = v.SubredditCategory
			ds_items[i].NumSubscribers = v.NumSubscribers
			ds_items[i].Url = ensureRedditDotCom(v.Url)
		case POST:
			if v.PostText == "" {
				// then this is a URL posting and maintain the url
				ds_items[i].Text = v.Url
			} else {
				// this post has content written by the author. The url doesnt  matter here
				ds_items[i].Text = v.PostText
			}
		case COMMENT:
			ds_items[i].Parent = v.Parent
			ds_items[i].Text = v.CommentBody
		}
		// trim text field
		// there is no conceivable reason for you to write a novel as part of a post or subreddit description.
		// If you do so, bruh take a chill. We are going to trim this
		ds_items[i].Text = utils.TruncateString(ds_items[i].Text, getMaxTextSize())

		// transform for processing queue and user content tracking
		pq_items[i] = ds.UserActionData{
			UserId:    user_id,
			Source:    REDDIT_PRETTY_URL,
			ContentId: v.Name,
			Processed: true,
		}
		// post processing id
		pq_items[i].RecordId = fmt.Sprintf("%s_processed_%s/%s", pq_items[i].UserId, pq_items[i].Source, pq_items[i].ContentId)
	}
	return ds_items, pq_items
}

func getUserActionsAndContents() ([]ds.UserActionData, []ds.ContentStoreData) {
	ua_data := ds.Deque(ds.USER_ACTION)
	var content_data []ds.ContentStoreData
	go_linq.From(ua_data).Select(func(ua any) any {
		return ds.GetContentFromStore(ua.(ds.UserActionData).ContentId)
	}).ToSlice(&content_data)
	return ua_data, content_data
}

// prefixes www.reddit.com to the URL it it is not there
func ensureRedditDotCom(url string) string {
	if !strings.HasPrefix(url, REDDIT_PRETTY_URL) {
		return REDDIT_PRETTY_URL + url
	}
	return url
}

const REDDIT_PRETTY_URL = "https://www.reddit.com"
