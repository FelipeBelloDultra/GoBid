package user

import (
	"context"

	"github.com/FelipeBelloDultra/go-bid/internal/validator"
)

type CreateUserReq struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserReq) Valid(context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(
		validator.NotBlank(req.UserName),
		"user_name",
		"this field cannot be blank",
	)
	eval.CheckField(
		validator.NotBlank(req.Email),
		"email",
		"this field cannot be blank",
	)
	eval.CheckField(
		validator.NotBlank(req.Bio),
		"bio",
		"this field cannot be blank",
	)
	eval.CheckField(
		validator.MinChars(req.Bio, 10) && validator.MaxChars(req.Bio, 255),
		"bio",
		"this field must have length between 10 and 255 characters",
	)
	eval.CheckField(
		validator.MinChars(req.Password, 8),
		"password",
		"this field must have length at least 8 characters",
	)
	eval.CheckField(
		validator.Matches(
			req.Email,
			validator.EmailRegex,
		),
		"email",
		"this field must be a valid email address",
	)

	return eval
}
