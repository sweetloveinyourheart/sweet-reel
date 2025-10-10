package actions

import (
	"context"
	"database/sql"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/models"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/repos"
)

func (a *actions) UpsertOAuthUser(ctx context.Context, request *connect.Request[proto.UpsertOAuthUserRequest]) (*connect.Response[proto.UpsertOAuthUserResponse], error) {
	isNewUser := false

	var user *models.User
	user, err := a.userRepo.GetUserWithIdentity(ctx, request.Msg.Provider, request.Msg.ProviderUserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, grpc.InternalError(err)
	}

	if user == nil {
		newUserID := uuid.Must(uuid.NewV7())
		newUser := &models.User{
			ID:      newUserID,
			Email:   request.Msg.Email,
			Name:    request.Msg.Name,
			Picture: request.Msg.Picture,
		}

		newUserIdentity := &models.UserIdentity{
			ID:             uuid.Must(uuid.NewV7()),
			UserID:         newUserID,
			Provider:       request.Msg.Provider,
			ProviderUserID: request.Msg.ProviderUserId,
		}

		// Execute transaction
		err = a.dbConn.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
			return db.TransactionC(ctx, conn.Conn(), "CreateUserAndIdentity", func(tx pgx.Tx) error {
				userRepoTx := repos.NewUserRepository(tx)

				err := userRepoTx.CreateUser(ctx, newUser)
				if err != nil {
					return err
				}

				err = userRepoTx.CreateIdentity(ctx, newUserIdentity)
				if err != nil {
					return err
				}

				return nil
			})
		})

		if err != nil {
			return nil, grpc.InternalError(err)
		}

		user = newUser
		isNewUser = true
	}

	response := &proto.UpsertOAuthUserResponse{
		User: &proto.User{
			Id:      user.ID.String(),
			Email:   user.Email,
			Name:    user.Name,
			Picture: user.Picture,
		},
		IsNewUser: isNewUser,
	}

	return connect.NewResponse(response), nil
}
