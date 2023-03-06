package app

import (
	"errors"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

func (s *courseService) Apply(cmd *PlayerApplyCmd) (code string, err error) {
	course, err := s.courseRepo.FindCourse(cmd.CourseId)
	if err != nil {
		return
	}

	if course.IsOver() {
		err = errors.New("course is over")

		return
	}

	if course.IsPreliminary() {
		err = errors.New("course is preparing")

		return
	}

	p := cmd.toPlayer()
	p.CreateToday()
	p.NewId()

	if err = s.playerRepo.SavePlayer(&p); err != nil {
		if repoerr.IsErrorDuplicateCreating(err) {
			code = errorDuplicateApply
		}

		return
	}

	if err = s.userCli.AddUserRegInfo(&p.Student); err != nil {
		return
	}

	return
}
