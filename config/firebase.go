package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	credentialsJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON")
	credentialsFile := os.Getenv("FIREBASE_CREDENTIALS_FILE")

	if credentialsJSON == "" && credentialsFile == "" {
		log.Fatal("FIREBASE_CREDENTIALS_JSON or FIREBASE_CREDENTIALS_FILE must be set")
	}

	var opt option.ClientOption
	if credentialsJSON != "" {
		opt = option.WithCredentialsJSON([]byte(credentialsJSON))
	} else {
		opt = option.WithCredentialsFile(credentialsFile)
	}

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Firebase init error: %v", err)
	}

	FirebaseApp = app
	log.Println("✅ Firebase initialized successfully")
}
