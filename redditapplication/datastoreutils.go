package redditapplication

import (
	"strings"

	cs "github.com/soumitsalman/goredditor/socialmediadatastore/contentstore"
	pq "github.com/soumitsalman/goredditor/socialmediadatastore/processingqueue"
	"github.com/soumitsalman/goredditor/utils"
)

func initializeDataStores() {
	cs.InitializeContentStoreClient()
	pq.InitializeProcessingQueues()
}

func saveNewItemsToDB(user_id string, item_kind string, items []RedditData) {
	cs.AddNewItems(item_kind, getNormalizedDataForCS(item_kind, items))
	pq.BatchQue(pq.NEW, getNormalizedDataForPQ(user_id, "www.reddit.com", items))
}

func getNormalizedDataForCS(item_kind string, items []RedditData) []cs.ContentStoreData {
	ds_items := make([]cs.ContentStoreData, len(items))
	for i, v := range items {
		ds_items[i] = cs.ContentStoreData{
			//applies to all
			Name:        v.Name,
			Id:          v.Id,
			Title:       v.Title,
			Kind:        v.Kind,
			CreatedDate: v.CreatedDate,

			//applies to subreddit
			NumSubscribers: v.NumSubscribers,

			//applies to post and comment
			Channel:     v.Subreddit,
			Author:      v.Author,
			Score:       v.Score,
			UpvoteRatio: v.UpvoteRatio,
			NumComments: v.NumComments,
			Url:         ensureRedditDotCom(v.Link),
		}
		//special field overrides
		switch v.Kind {
		case SUBREDDIT:
			ds_items[i].Channel = v.DisplayName
			ds_items[i].Text = utils.ExtractTextFromHtml(v.PublicDescription) + "\n" + utils.ExtractTextFromHtml(v.Description)
			ds_items[i].Category = v.SubredditCategory
		case POST:
			ds_items[i].Text = v.PostText + " " + v.Url
		case COMMENT:
			ds_items[i].Parent = v.Parent
			ds_items[i].Text = v.CommentBody
			ds_items[i].Category = v.PostCategory
		}
	}
	return ds_items
}

func getNormalizedDataForPQ(user_id string, source string, items []RedditData) []pq.ContentStoreDataRef {
	pq_items := make([]pq.ContentStoreDataRef, len(items))
	for i, v := range items {
		pq_items[i] = pq.ContentStoreDataRef{
			UserId: user_id,
			Source: source,
			Name:   v.Name,
		}
	}
	return pq_items
}

// prefixes www.reddit.com to the URL it it is not there
func ensureRedditDotCom(url string) string {
	if !strings.HasPrefix(url, REDDIT_PRETTY_URL) {
		return REDDIT_PRETTY_URL + url
	}
	return url
}

const REDDIT_PRETTY_URL = "www.reddit.com"
