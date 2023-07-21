package logic

import (
	"ginQAP/dao/mysql"
	"ginQAP/models"
	"go.uber.org/zap"
)


// GetIssues 获得所有帖子
func GetIssues() ([]*models.Issue, error) {
	// 查询数据库
	return mysql.GetIssues()
}


// GetIssueByID 获得某个帖子
func GetIssueByID(iid int64) (data *models.IssueDetail, err error) {
	// 查询某个帖子
	issue, err := mysql.GetIssueByID(iid)
	if err != nil {
		zap.L().Error("Issue query failed", zap.Error(err))
		return
	}

	// 查询某个帖子下的答案
	answers, err := mysql.GetAnswers(iid)
	if err != nil {
		zap.L().Error("Answers query failed", zap.Error(err))
		return
	}

	// 组合
	data = &models.IssueDetail{
		Issue:     issue,
		Answers:   answers,
	}
	return
}


// PostIssue 发布帖子
func PostIssue(p *models.IssuePost) (err error) {
	// 访问数据库
	return mysql.InsertIssue(p)
}


// PostAnswer 发布答案
func PostAnswer(p *models.AnswerPost) (err error) {
	// 访问数据库
	return mysql.InsertAnswer(p)
}
