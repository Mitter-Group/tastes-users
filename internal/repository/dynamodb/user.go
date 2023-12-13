package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chunnior/users/internal/domain"
	"github.com/chunnior/users/internal/domain/user"
	"github.com/chunnior/users/pkg/config"
	"github.com/google/uuid"
)

type UserRepository struct {
	db     *dynamodb.DynamoDB
	config *config.Config
	logger domain.Logger
}

func NewUserRepository(config *config.Config, logger domain.Logger) *UserRepository {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile:           config.AwsProfile,
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		panic(err)
	}

	return &UserRepository{
		db:     dynamodb.New(sess),
		config: config,
		logger: logger,
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, user user.GenericUser) (string, error) {
	// Primero, intenta obtener el usuario por correo electrónico
	existingUser, err := r.GetUserByEmail(ctx, user.Email)
	if err != nil {
		// Si hay un error al obtener el usuario, devuélvelo
		r.logger.Error("error getting user by email", "error", err)
		return "", err
	}

	if existingUser != nil {
		// Si el usuario existe, actualiza sus datos
		existingUser.Provider = user.Provider
		existingUser.ProviderUserID = user.ProviderUserID
		existingUser.UserFullname = user.UserFullname
		existingUser.AccessToken = user.AccessToken
		existingUser.RefreshToken = user.RefreshToken

		err := r.UpdateUser(ctx, *existingUser)
		if err != nil {
			return "", err
		}

		return existingUser.ID, nil
	} else {
		// Si el usuario no existe, crea uno nuevo
		newUserID, err := r.CreateUser(ctx, &user)
		if err != nil {
			return "", err
		}

		return newUserID, nil
	}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.GenericUser, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.config.AwsDynamoUserTableName),
		IndexName: aws.String("EmailIndex"), // Asegúrate de que este nombre coincide con el nombre del índice global secundario que creaste en Terraform
		KeyConditions: map[string]*dynamodb.Condition{
			"Email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(email),
					},
				},
			},
		},
	}

	result, err := r.db.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var user user.GenericUser
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, updatedUser user.GenericUser) error {
	expressionValues := map[string]*dynamodb.AttributeValue{
		":p":  {S: aws.String(updatedUser.Provider)},
		":pu": {S: aws.String(updatedUser.ProviderUserID)},
		":uf": {S: aws.String(updatedUser.UserFullname)}, // Aquí ajustamos el nombre
		":e":  {S: aws.String(updatedUser.Email)},
		":at": {S: aws.String(updatedUser.AccessToken)},
		":rt": {S: aws.String(updatedUser.RefreshToken)},
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expressionValues,
		TableName:                 aws.String(r.config.AwsDynamoUserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(updatedUser.ID)},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Provider = :p, ProviderUserID = :pu, UserFullname = :uf, Email = :e, AccessToken = :at, RefreshToken = :rt"),
	}
	var err error
	_, err = r.db.UpdateItemWithContext(ctx, input)
	if err != nil {
		r.logger.Error("error updating user", "error", err)
	}
	return err
}

func (r *UserRepository) CreateUser(ctx context.Context, user *user.GenericUser) (string, error) {
	newID := uuid.New().String()
	user.ID = newID

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return "", err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.config.AwsDynamoUserTableName),
		Item:      item,
	}

	_, err = r.db.PutItemWithContext(ctx, input)
	if err != nil {
		return "", err
	}

	// Devuelve el nuevo ID del usuario
	return newID, nil
}
