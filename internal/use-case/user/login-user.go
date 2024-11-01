package user

import (
	"context"

	"github.com/FelipeBelloDultra/go-bid/internal/validator"
)

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginUserReq) Valid(context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(
		validator.Matches(
			req.Email,
			validator.EmailRegex,
		),
		"email",
		"this field must be a valid email address",
	)
	eval.CheckField(
		validator.NotBlank(req.Password),
		"password",
		"this field cannot be blank",
	)

	return eval
}
