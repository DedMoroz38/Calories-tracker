package dto

// PhotoResponse is a single photo with a presigned, temporary view URL. Used
// for the current user's own photo grid on the profile page.
type PhotoResponse struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

// FeedItemResponse is a photo in the public feed, including minimal author
// info so the feed can show who posted it. No likes or comments by design.
type FeedItemResponse struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	UserID    uint   `json:"user_id"`
	AuthorName string `json:"author_name"`
	AuthorAvatar string `json:"author_avatar"`
	CreatedAt string `json:"created_at"`
}

// FeedResponse wraps a page of feed items with the cursor to fetch the next
// page (the id of the last item). NextCursor is 0 when there are no more.
type FeedResponse struct {
	Items      []FeedItemResponse `json:"items"`
	NextCursor uint               `json:"next_cursor"`
}
