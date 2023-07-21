package mysql

import (
	"database/sql"
	"ginQAP/models"
	"go.uber.org/zap"
)


// GetIssues 查询数据库中所有帖子，以列表形式返回
func GetIssues() (issues []*models.Issue, err error) {
	sqlStr := `select issues.id as id, title, content, issues.create_time as create_time, author_id, user.username as author_name  
			   from issues join user on issues.author_id = user.user_id`
	err = db.Select(&issues, sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		} else {
			zap.L().Error("Issues query failed", zap.Error(err))
			return
		}
	}
	return
}


// GetIssueByID 获得某个帖子
func GetIssueByID(iid int64) (issue *models.Issue, err error) {
	issue = new(models.Issue)
	sqlStr := `select issues.id as id, title, content, issues.create_time as create_time, author_id, user.username as author_name  
			   from issues join user on issues.author_id = user.user_id where issues.id = ?`
	err = db.Get(issue, sqlStr, iid)
	return
}


// GetAnswers 获得某个帖子的答案
func GetAnswers(iid int64) (answers []*models.Answer, err error) {
	sqlStr := `select answers.id as id, content, answers.create_time as create_time, author_id, 
			   user.username as author_name, issue_id from answers 
			   join user on answers.author_id = user.user_id where issue_id = ?`
	err = db.Select(&answers, sqlStr, iid)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no answer in db")
			err = nil
		} else {
			zap.L().Error("Answers query failed", zap.Error(err))
			return
		}
	}
	return
}


// InsertIssue 向数据库插入一条新的帖子
func InsertIssue(p *models.IssuePost) (err error) {
	// 插入
	sqlStr := `insert into issues (title, content, author_id) values(?, ?, ?)`
	_, err = db.Exec(sqlStr, p.Title, p.Content, p.AuthorID)
	if err != nil {
		zap.L().Error("Insert issue failded", zap.Error(err))
	}
	return
}


// InsertAnswer 向数据库插入一条新的答案
func InsertAnswer(p *models.AnswerPost) (err error) {
	// 插入
	sqlStr := `insert into answers (content, author_id, issue_id) values(?, ?, ?)`
	_, err = db.Exec(sqlStr, p.Content, p.AuthorID, p.IssueID)
	if err != nil {
		zap.L().Error("Insert answer failded", zap.Error(err))
	}
	return
}