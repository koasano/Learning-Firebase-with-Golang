package main

import (
	"context"
	"flag"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"google.golang.org/api/iterator"
	"cloud.google.com/go/firestore"
)

func createDoc(ctx context.Context, client *firestore.Client) {
	_, _, err := client.Collection("users").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	log.Println("Successfully created document")
}

func readDoc(ctx context.Context, client *firestore.Client) {
	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		log.Printf("Document data: %#v\n", doc.Data())
	}
	log.Println("Successfully read document")
}

func updateDoc(ctx context.Context, client *firestore.Client, docId string) {
	_, err := client.Collection("users").Doc(docId).Set(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Byron",
		"born":  1815,
	}, firestore.MergeAll)
	if err != nil {
		log.Fatalf("Failed updating alovelace: %v", err)
	}
	log.Println("Successfully update document")
}

func deleteDoc(ctx context.Context, client *firestore.Client, docId string) {
	_, err := client.Collection("users").Doc(docId).Delete(ctx)
	if err != nil {
		log.Fatalf("Failed deleting alovelace: %v", err)
	}
	log.Println("Successfully delete document")
}

func main() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccountKey.json")
	conf := &firebase.Config{ProjectID: "fir-godemo-e0791"}
	app, err := firebase.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
	}
	defer client.Close()

	action := flag.String("action", "", "The action to perform: 'create', 'read', 'update', or 'delete'")
	docId := flag.String("docId", "", "The document ID for 'update' and 'delete' actions")
	flag.Parse()

	switch *action {
	case "create":
		createDoc(ctx, client)
	case "read":
		readDoc(ctx, client)
	case "update":
		if *docId == "" {
			log.Fatal("docId must be provided for 'update' action")
		}
		updateDoc(ctx, client, *docId)
	case "delete":
		if *docId == "" {
			log.Fatal("docId must be provided for 'delete' action")
		}
		deleteDoc(ctx, client, *docId)
	default:
		log.Fatalf("Invalid action: %s", *action)
	}
}

