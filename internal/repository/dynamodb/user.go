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

func (r *UserRepository) SaveUser(ctx context.Context, userInput user.GenericUser) (*user.UserData, error) {
	// Primero, intenta obtener el usuario por correo electrónico
	existingUser, err := r.GetUserByEmail(ctx, userInput.Email)
	if err != nil {
		// Si hay un error al obtener el usuario, devuélvelo
		r.logger.Error("error getting user by email", "error", err)
		return nil, err
	}

	if existingUser != nil {
		// Si el usuario existe, actualiza sus datos
		existingUser = updateOrAddProvider(userInput, existingUser)

		err := r.UpdateUserProviders(ctx, *existingUser)
		if err != nil {
			return nil, err
		}

		return existingUser, nil
	} else {
		newUser := user.UserData{
			Email: userInput.Email,
			Providers: []user.ProviderData{
				{
					Provider:       userInput.Provider,
					ProviderUserID: userInput.ProviderUserID,
					UserFullname:   userInput.UserFullname,
					Email:          userInput.Email,
				},
			},
		}
		// Si el usuario no existe, crea uno nuevo
		newUserID, err := r.CreateUser(ctx, &newUser)
		if err != nil {
			return nil, err
		}
		newUser.ID = newUserID
		return &newUser, nil
	}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.UserData, error) {
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

	var user user.UserData
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func updateOrAddProvider(gUser user.GenericUser, userData *user.UserData) *user.UserData {
	found := false

	// Revisar cada proveedor existente en UserData
	for i, provider := range userData.Providers {
		if provider.Provider == gUser.Provider && provider.ProviderUserID == gUser.ProviderUserID {
			// Actualiza el nombre si el proveedor y el ID del proveedor coinciden
			userData.Providers[i].UserFullname = gUser.UserFullname
			found = true
			break
		}
	}

	// Si no se encontró un proveedor existente, añade uno nuevo
	if !found {
		newProvider := user.ProviderData{
			Provider:       gUser.Provider,
			ProviderUserID: gUser.ProviderUserID,
			UserFullname:   gUser.UserFullname,
			Email:          gUser.Email,
		}
		userData.Providers = append(userData.Providers, newProvider)
	}

	return userData
}

func (r *UserRepository) UpdateUserProviders(ctx context.Context, updatedUser user.UserData) error {

	providersAV, err := dynamodbattribute.MarshalList(updatedUser.Providers)
	if err != nil {
		r.logger.Error("error marshalling providers", "error", err)
		return err
	}

	expressionValues := map[string]*dynamodb.AttributeValue{
		":p": {L: providersAV},
	}
	// Preparar el input para UpdateItem
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.config.AwsDynamoUserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(updatedUser.ID)},
		},
		UpdateExpression:          aws.String("set Providers = :p"),
		ExpressionAttributeValues: expressionValues,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	// Realizar la actualización en DynamoDB
	_, err = r.db.UpdateItemWithContext(ctx, input)
	if err != nil {
		r.logger.Error("error updating user providers", "error", err)
	}
	return err
	/*
		expressionValues := map[string]*dynamodb.AttributeValue{
			":p":  {S: aws.String(updatedUser.Provider)},
			":pu": {S: aws.String(updatedUser.ProviderUserID)},
			":uf": {S: aws.String(updatedUser.UserFullname)}, // Aquí ajustamos el nombre
			":e":  {S: aws.String(updatedUser.Email)},
		}

		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: expressionValues,
			TableName:                 aws.String(r.config.AwsDynamoUserTableName),
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {S: aws.String(updatedUser.ID)},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set Provider = :p, ProviderUserID = :pu, UserFullname = :uf, Email = :e"),
		}
		// , AccessToken = :at, RefreshToken = :rt
		var err error
		_, err = r.db.UpdateItemWithContext(ctx, input)
		if err != nil {
			r.logger.Error("error updating user", "error", err)
		}
		return err
	*/
}

func (r *UserRepository) CreateUser(ctx context.Context, user *user.UserData) (string, error) {
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
