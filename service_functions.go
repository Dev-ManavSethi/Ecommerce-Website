package main

import (
	"context"
	"errors"

	"github.com/fatih/color"

	"github.com/DevManavSethi/EcommerceWebsite/service"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func (*server) Signup(ctx context.Context, req *service.SignupRequest) (*service.SignupResponse, error) {

	user := req.GetUser()

	pass := user.GetPassword()

	//--------------------------------
	EncryptedPass, err00 := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err00 != nil {
		FatalOnError("Error encrypting password", err00)
	}
	user.Password = string(EncryptedPass)
	//-------------------------------
	uuid, err000 := uuid.NewV4()
	FatalOnError("Error ", err000)
	user.ID = uuid.String()

	//---------------------------------------------------------------------------------

	mongoDBclient, err01 := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	FatalOnError("Error connecting to mongoDB.", err01)

	_, err02 := mongoDBclient.Database("ecommerce").Collection("users").InsertOne(context.TODO(), user)
	if err02 != nil {
		FatalOnError("error insering in mongodb", err02)
		return nil, err02
	}

	//------------------------------------------------------------------------------------
	return &service.SignupResponse{
		User: user,
	}, nil

}

func (*server) AddToCart(ctx context.Context, req *service.AddToCartRequest) (*service.AddToCartResponse, error) {

	productID := req.ProductID
	Quantity := req.Quantity
	comments := req.Comments
	userID := req.UserID

	color.Yellow("params", productID, Quantity, comments, userID)

	mongoDBclient, err01 := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err01 != nil {
		FatalOnError("Error connecting to mongoDB.", err01)
		return nil, err01
	}

	result := mongoDBclient.Database("ecommerce").Collection("users").FindOne(context.TODO(), bson.M{"id": userID})
	if result == nil {

		return nil, errors.New("User not found")

	}
	var user *service.User

	err02 := result.Decode(&user)
	color.Green("user", *user)
	if err02 != nil {
		FatalOnError("Error decoding", err02)
		return nil, err02

	}

	cart := user.GetCart()
	//	color.Red("cart", *cart)

	if cart == nil {

		FistCart := service.Cart{}
		FistCart.ProductIDs = append(FistCart.ProductIDs, productID)
		FistCart.Quantities = append(FistCart.Quantities, Quantity)
		FistCart.Comments = append(FistCart.Comments, comments)

		SingleResult := mongoDBclient.Database("ecommerce").Collection("users").FindOneAndUpdate(context.TODO(), bson.M{"id": userID}, bson.M{"$set": bson.M{"cart": FistCart}})

		color.Blue("singleresult", *SingleResult)

		if SingleResult.Err() != nil {
			return nil, SingleResult.Err()
		}
		return &service.AddToCartResponse{
			Success: true,
		}, nil

	} else {

		cart.ProductIDs = append(cart.ProductIDs, productID)
		cart.Quantities = append(cart.Quantities, Quantity)
		cart.Comments = append(cart.Comments, comments)

		SingleResult := mongoDBclient.Database("ecommerce").Collection("users").FindOneAndUpdate(context.TODO(), bson.M{"id": userID}, bson.M{"$set": bson.M{"cart": *cart}})

		color.Blue("singleresult", *SingleResult)

		if SingleResult.Err() != nil {
			return nil, SingleResult.Err()
		}
		return &service.AddToCartResponse{
			Success: true,
		}, nil
	}

}
func (*server) Login(ctx context.Context, req *service.LoginRequest) (*service.LoginResponse, error) {
	email := req.GetEmail()
	pass := req.GetPassword()
	//KeepLoggedIn := req.ToKeepLoggedIn

	mongoDBclient, err01 := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	FatalOnError("Error connecting to mongoDB.", err01)

	//EncryptedPass, err00 := bcrypt.GenerateFromPassword([]byte(pass), 10)

	//FatalOnError("Error encrypting password", err00)

	result := mongoDBclient.Database("ecommerce").Collection("users").FindOne(context.TODO(), bson.M{"email": email})
	if result == nil {

		return nil, errors.New("User not found")

	}
	var user *service.User

	err02 := result.Decode(&user)
	if err02 != nil {
		FatalOnError("Error decoding", err02)
		return nil, err02

	}

	error001 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

	if user.Email != email || error001 != nil {
		return nil, error001
	}

	return &service.LoginResponse{
		User: user,
	}, nil

}
func (*server) Logout(ctx context.Context, req *service.LogoutRequest) (*service.LogoutResponse, error) {

	//erase cookie

	return &service.LogoutResponse{
		Success: true,
	}, nil

}

func (*server) ReadCart(ctx context.Context, req *service.ReadCartRequest) (*service.ReadCartResponse, error) {
	return nil, nil
}
func (*server) UpdateCart(ctx context.Context, req *service.UpdateCartRequest) (*service.UpdateCartResponse, error) {
	return nil, nil
}
func (*server) EmptyCart(ctx context.Context, req *service.EmptyCartRequest) (*service.EmptyCartResponse, error) {
	return nil, nil
}
func (*server) Checkout(ctx context.Context, req *service.CheckoutRequest) (*service.CheckoutResponse, error) {
	return nil, nil
}
func (*server) EditUser(ctx context.Context, req *service.EditUserRequest) (*service.EditUserResponse, error) {

	// userBefore := req.GetUserBefore()
	// userAfter := req.GetUserAfter()

	return nil, nil
}
func (*server) Pay(ctx context.Context, req *service.PayRequest) (*service.PayResponse, error) {
	return nil, nil
}
