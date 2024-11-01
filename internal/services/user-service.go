package services

import (
	"context"
	"errors"

	"github.com/FelipeBelloDultra/go-bid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

var (
	ErrDuplicatedEmailOrUsername = errors.New("username or email already exists")
	ErrInvalidCredentials        = errors.New("invalid credentials")
)

func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (us *UserService) CreateUser(ctx context.Context, userName, email, password, bio string) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.UUID{}, err
	}

	args := pgstore.CreateUserParams{
		UserName:     userName,
		Email:        email,
		PasswordHash: hash,
		Bio:          bio,
	}

	id, err := us.queries.CreateUser(ctx, args)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" { // <- postgres code of unique constraint violation
			return uuid.UUID{}, ErrDuplicatedEmailOrUsername
		}

		return uuid.UUID{}, err
	}

	return id, nil
}

func (us *UserService) AuthenticateUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	user, err := us.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, ErrInvalidCredentials
		}

		return uuid.UUID{}, err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.UUID{}, ErrInvalidCredentials
		}

		return uuid.UUID{}, err
	}

	return user.ID, nil
}
