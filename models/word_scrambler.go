package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// AnswerStatus -
type AnswerStatus string

const (
	// Correct -
	Correct AnswerStatus = "correct"
	// Wrong -
	Wrong AnswerStatus = "wrong"
)

// Question -
type Question struct {
	ID           uint64     `gorm:"column:id; type:int(11) unsigned auto_increment; not null" json:"id"`
	Question     string     `gorm:"column:question; type:varchar(50); not null" json:"question"`
	CorrectPoint uint64     `gorm:"column:correct_point; type:int(11); default:1" json:"correct_point"`
	WrongPoint   int64      `gorm:"column:wrong_point; type:int(11); default:1" json:"wrong_point"`
	CreatedAt    *time.Time `gorm:"column:created_at; type:datetime; not null; default:current_timestamp" json:"created_at"`
}

// TableName -
func (Question) TableName() string {
	return "word_scrambler.questions"
}

// User -
type User struct {
	ID        uint64     `gorm:"column:id; type:int(11) unsigned auto_increment; not null" json:"id"`
	Username  string     `gorm:"column:username; type:varchar(50); not null" json:"username"`
	CreatedAt *time.Time `gorm:"column:created_at; type:datetime; not null; default:current_timestamp" json:"created_at"`
}

// TableName -
func (User) TableName() string {
	return "word_scrambler.users"
}

// UserAnswers -
type UserAnswers struct {
	ID         uint64       `gorm:"column:id; type:int(11) unsigned auto_increment; not null" json:"id"`
	Question   Question     `gorm:"foreignKey:QuestionID; references:ID" json:"-"`
	QuestionID uint64       `gorm:"column:question_id; type:int(11) unsigned; not null" json:"question_id"`
	Answer     string       `gorm:"column:answer; type:varchar(50); default:null" json:"answer"`
	Status     AnswerStatus `gorm:"column:status; type:ENUM('correct', 'wrong')" json:"status"`
	User       User         `gorm:"foreignKey:UserID; references:ID" json:"-"`
	UserID     uint64       `gorm:"column:user_id; type:int(11) unsigned; not null" json:"user_id"`
	CreatedAt  *time.Time   `gorm:"column:created_at; type:datetime; not null; default:current_timestamp" json:"created_at"`
}

// TableName -
func (UserAnswers) TableName() string {
	return "word_scrambler.user_answers"
}

// UserSummary -
type UserSummary struct {
	ID         uint64     `gorm:"column:id; type:int(11) unsigned auto_increment; not null" json:"id"`
	User       User       `gorm:"foreignKey:UserID; references:ID" json:"-"`
	UserID     uint64     `gorm:"column:user_id; type:int(50) unsigned; not null" json:"user_id"`
	TotalScore int64      `gprm:"column:total_score; type:int(11); not null; default:0" json:"total_score"`
	CreatedAt  *time.Time `gorm:"column:created_at; type:datetime; not null; default:current_timestamp" json:"created_at"`
}

// TableName -
func (UserSummary) TableName() string {
	return "word_scrambler.user_summary"
}

// GetRandQuestion - function to get randome question
func GetRandQuestion(DB *gorm.DB) (Question, error) {
	query := `
	SELECT 
		id,
		question,
		correct_point,
		wrong_point,
		created_at
	FROM word_scrambler.questions
	ORDER BY RAND()
	LIMIT 1  
	`
	var res Question
	if err := DB.Raw(query).Scan(&res).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		fmt.Println(err)
		return res, err
	}
	return res, nil
}

// GetAnswer -
func GetAnswer(DB *gorm.DB, questionID uint64) (Question, error) {
	query := `
		SELECT 
			*
		FROM word_scrambler.questions
		WHERE id = ?
	`

	var res Question
	if err := DB.Raw(query, questionID).Scan(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

// InsertUserAnswer -
func InsertUserAnswer(DB *gorm.DB, data *UserAnswers) error {
	if err := DB.Create(data).Error; err != nil {
		return err
	}
	return nil
}

//GetUserID -
func GetUserID(DB *gorm.DB, name string) (User, error) {
	var u User
	if err := DB.Raw(`
		SELECT * FROM word_scrambler.users WHERE username = ?
	`, name).Scan(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}

type TotalScores struct {
	UserID     uint64 `gorm:"-" json:"user_id"`
	TotalScore int64  `gorm:"column:total_score" json:"total_score"`
}

// TotalScore -
func TotalScore(DB *gorm.DB, userID uint64) (TotalScores, error) {
	var tot TotalScores
	query := `
		SELECT SUM(score) as total_score FROM (
			SELECT 
				IF(wsua.status = 'wrong', wsc.wrong_point, wsc.correct_point) as score
			FROM word_scrambler.user_answers wsua
			INNER JOIN word_scrambler.questions wsc ON wsc.id = wsua.question_id
			WHERE wsua.user_id = ?
		) sub
	`

	if err := DB.Raw(query, userID).Scan(&tot).Error; err != nil {
		return tot, err
	}
	tot.UserID = userID
	return tot, nil
}

type SummaryData struct {
	UserAnswers
	Scores int64 `gorm:"column:score" json:"score"`
}

// UserSummaries -
func UserSummaries(DB *gorm.DB, userID uint64) ([]SummaryData, error) {
	var data []SummaryData
	query := `
		SELECT 
			user_answers.*,
			IF(user_answers.status = 'wrong', questions.wrong_point, questions.correct_point) as score
		FROM word_scrambler.user_answers
		LEFT JOIN word_scrambler.questions on questions.id = user_answers.question_id
		where user_answers.user_id = ?
	`

	if err := DB.Raw(query, userID).Scan(&data).Error; err != nil {
		return data, err
	}

	return data, nil
}
