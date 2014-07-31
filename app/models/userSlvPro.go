package models

type UserProb struct {
	UserId int64
	ProbId int64
	Solved bool
}

func (usp *UserProb) Save(userId, probId int64, ok bool) (int64, error) {
	usp.UserId = userId
	usp.ProbId = probId
	usp.Solved = ok
	return engine.Insert(usp)
}
func FindSovledProblems(userId int64) ([]UserProb, error) {
	var s []UserProb
	err := engine.Table("problem").Join("INNER", "user_prob", `"user_prob".prob_id = "problem".id`).Join("INNER", `"user"`, `"user".id = "user_prob".user_id`).Cols("user.solved", "user_id", "prob_id").Where("user_id = ?", userId).Find(&s)
	if err != nil {
		return nil, err
	} else {
		return s, nil
	}
}
