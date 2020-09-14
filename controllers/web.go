package controllers

import (
	"fmt"
	"strconv"
	"word-scrambler/libs"
	"word-scrambler/models"
	"word-scrambler/pkg"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

// GetQuestion -
func GetQuestion(c echo.Context) error {
	question, err := models.GetRandQuestion(pkg.DB)
	if err != nil {
		return pkg.InternalErrResp(c, err)
	}
	question.Question = libs.Scrambler(question.Question)

	return pkg.Data(c, question)
}

// Answers -
func Answers(c echo.Context) error {
	questionID, err := strconv.ParseUint(c.Param("question_id"), 10, 64)
	if err != nil {
		return pkg.NotFoundResponse(c)
	}
	var req struct {
		Answer string `json:"answer"`
	}

	if err := c.Bind(&req); err != nil {
		return pkg.BadRequestResp(c, "Invalid input format")
	}

	ans, err := models.GetAnswer(pkg.DB, questionID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return pkg.NotFoundResponse(c)
		}

		return pkg.InternalErrResp(c, err)
	}

	uans := models.UserAnswers{
		QuestionID: questionID,
		Answer:     req.Answer,
		UserID:     c.Get("user_id").(uint64),
	}
	fmt.Println(uans)

	if ans.Question == req.Answer {
		uans.Status = models.Correct
	} else {
		uans.Status = models.Wrong
	}

	tx := pkg.DB.Begin()
	if err := models.InsertUserAnswer(tx, &uans); err != nil {
		tx.Rollback()
		return pkg.InternalErrResp(c, err)
	}
	tx.Commit()

	return pkg.Data(c, uans)
}

// GetTotalScore -
func GetTotalScore(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	tot, err := models.TotalScore(pkg.DB, userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return pkg.InternalErrResp(c, err)
	}
	return pkg.Data(c, tot)

}

// Summary -
func Summary(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	data, err := models.UserSummaries(pkg.DB, userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return pkg.InternalErrResp(c, err)
	}
	return pkg.Data(c, data)
}
