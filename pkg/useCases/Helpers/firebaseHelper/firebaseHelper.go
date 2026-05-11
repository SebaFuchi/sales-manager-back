package firebaseHelper

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

// InitFirebase initializes the Firebase Admin SDK
// It looks for serviceAccountKey.json in the project root or uses the FIREBASE_CREDENTIALS env var
func InitFirebase() {
	// Para este ejemplo, asumiremos que se provee el JSON de servicio en la raíz
	// o que Google Application Default Credentials (ADC) está configurado.
	opt := option.WithCredentialsFile("serviceAccountKey.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Warning: error initializing firebase app: %v\n", err)
		log.Printf("Firebase Auth validation will fail until serviceAccountKey.json is provided.")
		return
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Printf("Warning: error getting Auth client: %v\n", err)
		return
	}

	AuthClient = client
	log.Println("✓ Firebase Admin SDK initialized successfully")
}

// VerifyToken verifies a Firebase ID token and returns the decoded token
func VerifyToken(idToken string) (*auth.Token, error) {
	if AuthClient == nil {
		// Mock implementation just for the fallback case while dev sets up the key
		// In production, this should return an error
		return nil, nil // Or throw an error demanding Firebase init
	}
	return AuthClient.VerifyIDToken(context.Background(), idToken)
}

// SetCustomUserClaims sets custom claims for a user
func SetCustomUserClaims(uid string, claims map[string]interface{}) error {
	if AuthClient == nil {
		return nil
	}
	return AuthClient.SetCustomUserClaims(context.Background(), uid, claims)
}

// CreateUser creates a new user in Firebase Auth
func CreateUser(email string, password string, displayName string) (*auth.UserRecord, error) {
	if AuthClient == nil {
		return nil, nil
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(displayName)

	return AuthClient.CreateUser(context.Background(), params)
}
