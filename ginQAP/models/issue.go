package models


// Issue 帖子
type Issue struct {
	IssueID     int64      `db:"id" json:"issue_id"`
	Title       string     `db:"title" json:"title"`
	Content     string     `db:"content" json:"content"`
	CreateTime  string     `db:"create_time" json:"create_time"`
	AuthorID    int64      `db:"author_id" json:"author_id"`
	AuthorName  string     `db:"author_name" json:"author_name"`
}

// Answer 帖子的答案
type Answer struct {
	AnswerID    int64      `db:"id" json:"answer_id"`
	Content     string     `db:"content" json:"content"`
	CreateTime  string     `db:"create_time" json:"create_time"`
	AuthorID    int64      `db:"author_id" json:"author_id"`
	AuthorName  string     `db:"author_name" json:"author_name"`
	IssueID     int64      `db:"issue_id" json:"issue_id"`
}


// IssueDetail 某个帖子的详情
type IssueDetail struct {
	*Issue
	Answers     []*Answer   `json:"answers"`
}


// IssuePost 发布帖子
type IssuePost struct {
	Title       string     `db:"title" json:"title" binding:"required"`
	Content     string     `db:"content" json:"content" binding:"required"`
	AuthorID    int64      `db:"author_id" json:"author_id"`
}


// AnswerPost 发布答案
type AnswerPost struct {
	Content     string     `db:"content" json:"content" binding:"required"`
	AuthorID    int64      `db:"author_id" json:"author_id"`
	IssueID     int64      `db:"issue_id" json:"issue_id"`
}