package main

import (
	"context"
	"flag"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	auth "firebase.google.com/go/v4/auth"
)

func createUser(authClient *auth.Client) *auth.UserRecord {
	params := (&auth.UserToCreate{}).
		Email("user@example.com").
		EmailVerified(false).
		Password("secretPassword").
		Disabled(false)

	u, err := authClient.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}

	log.Printf("Successfully created user: %v\n", u)
	log.Printf("User ID (UID): %s\n", u.UID)  // Output the user ID
	return u
}

func getUser(authClient *auth.Client, userId string) *auth.UserRecord {
	u, err := authClient.GetUser(context.Background(), userId)
	if err != nil {
		log.Fatalf("error getting user %s: %v\n", userId, err)
	}

	log.Printf("Successfully fetched user data: %v\n", u)
	return u
}

func updateUser(authClient *auth.Client, userId string) *auth.UserRecord {
	params := (&auth.UserToUpdate{}).
		Email("newemail@example.com").
		EmailVerified(true).
		Password("newPassword").
		Disabled(true)

	u, err := authClient.UpdateUser(context.Background(), userId, params)
	if err != nil {
		log.Fatalf("error updating user: %v\n", err)
	}

	log.Printf("Successfully updated user: %v\n", u)
	return u
}

func deleteUser(authClient *auth.Client, userId string) {
	err := authClient.DeleteUser(context.Background(), userId)
	if err != nil {
		log.Fatalf("error deleting user: %v\n", err)
	}

	log.Println("Successfully deleted user")
}

func main() {
	opt := option.WithCredentialsFile("./serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	action := flag.String("action", "", "The action to perform: 'create', 'get', 'update', or 'delete'")
	userId := flag.String("userId", "", "The user ID for 'get', 'update', and 'delete' actions")
	flag.Parse()

	switch *action {
	case "create":
		createUser(authClient)
	case "get":
		if *userId == "" {
			log.Fatal("userId must be provided for 'get' action")
		}
		getUser(authClient, *userId)
	case "update":
		if *userId == "" {
			log.Fatal("userId must be provided for 'update' action")
		}
		updateUser(authClient, *userId)
	case "delete":
		if *userId == "" {
			log.Fatal("userId must be provided for 'delete' action")
		}
		deleteUser(authClient, *userId)
	default:
		log.Fatalf("Invalid action: %s", *action)
	}
}
