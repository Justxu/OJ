package models

import ()

type UserSolvedProb struct {
	UserId int64
	ProId  int64
	Solved bool
}

func (usp *UserSolvedProb) Save(userId, probId int64, ok bool) (int64, error) {
	usp.UserId = userId
	usp.ProId = probId
	usp.Solved = ok
	return engine.Insert(usp)
}
func (usp *UserSolvedProb) FindSovledProblems(userId int64) []Problem {
	var probs []Problem
	// slect * from problem where user_solved_problem.sovled = ture and
	// user_solved_problem.user_id = userId and user_solved_problem.pro_id
	// = problem.id
	return probs
}
