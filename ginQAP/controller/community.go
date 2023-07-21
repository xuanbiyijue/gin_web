package controller

import (
	"fmt"
	"ginQAP/logic"
	"ginQAP/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)


// HandleGetIssues 获取所有帖子
func HandleGetIssues(c *gin.Context)  {
	// 因为没有参数，不需要进行参数校验
	// 查询所有帖子的 (id, title, content, author_id, create_time)
	issues, err := logic.GetIssues()
	if err != nil {
		zap.L().Error("Get issues failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	ResponseSuccess(c, issues)
}


// HandleGetIssueByID 获取某个帖子
func HandleGetIssueByID(c *gin.Context)  {
	// 获取参数，得到帖子的id
	issueIDStr := c.Param("issueID")
	issueID, err := strconv.ParseInt(issueIDStr, 10, 64)
	if err != nil {
		zap.L().Error("get issue detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 根据帖子id取数据
	issue, err := logic.GetIssueByID(issueID)
	if err != nil {
		zap.L().Error("logic.GetIssueByID(issueID) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, issue)
}


// HandlePostIssue 发布帖子
func HandlePostIssue(c *gin.Context)  {
	// 获取参数以及参数校验
	var param models.IssuePost
	err := c.ShouldBindJSON(&param)
	if err != nil {
		zap.L().Error("Post issue with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	authorID, exist := c.Get("userID")
	if !exist {
		zap.L().Error("Post issue with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	authorIDnum, err := strconv.ParseInt(fmt.Sprintf("%v",authorID), 10, 64)
	if err != nil {
		zap.L().Error("Change the authorIDnum failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	param.AuthorID = authorIDnum

	// 业务处理
	err = logic.PostIssue(&param)
	if err != nil {
		zap.L().Error("logic.PostIssue(param) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 成功返回响应
	ResponseSuccess(c, nil)
}


// HandlePostAnswer 发布答案
func HandlePostAnswer(c *gin.Context)  {
	// 获取参数以及参数校验
	var param models.AnswerPost
	err := c.ShouldBindJSON(&param)
	if err != nil {
		zap.L().Error("Post answer with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	issueIDStr := c.Param("issueID")
	issueID, err := strconv.ParseInt(issueIDStr, 10, 64)
	if err != nil {
		zap.L().Error("get issue id failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	param.IssueID = issueID
	authorID, exist := c.Get("userID")
	if !exist {
		zap.L().Error("Post issue with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	authorIDnum, err := strconv.ParseInt(fmt.Sprintf("%v",authorID), 10, 64)
	if err != nil {
		zap.L().Error("Change the authorIDnum failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	param.AuthorID = authorIDnum
	// 业务处理
	err = logic.PostAnswer(&param)
	if err != nil {
		zap.L().Error("logic.PostAnswer(&param) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 成功返回响应
	ResponseSuccess(c, nil)
}
