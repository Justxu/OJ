package models

type Solve struct {
	UserId int64
	ProbId int64
	Solved bool
}

func (sov *Solve) Save(userId, probId int64, ok bool) (int64, error) {
	s := &Solve{}
	err := engine.Where("user_id = ? and prob_id = ?", userId, probId).Find(s)
	sov.UserId = userId
	sov.ProbId = probId
	sov.Solved = ok
	if err == nil {
		return engine.Update(sov)
	} else {

		return engine.Insert(sov)
	}
}

func FindSovledProblems(userId int64) ([]Solve, error) {
	var s []Solve
	err := engine.Table("problem").Join("INNER", "solve", `"solve".prob_id = "problem".id`).Join("INNER", `"user"`, `"user".id = "solve".user_id`).Cols("user.solved", "user_id", "prob_id").Where("user_id = ?", userId).Find(&s)
	if err != nil {
		return nil, err
	} else {
		return s, nil
	}
}
