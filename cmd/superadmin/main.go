package main

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"sales-manager-back/pkg/domain/tenant"
	"sales-manager-back/pkg/domain/user"
)

func main() {
	email := "sebastianm.fuchilieri@gmail.com"
	dbURL := os.Getenv("CONNECTION_URL")
	if dbURL == "" {
		log.Fatal("Falta la variable de entorno CONNECTION_URL (Ej: root:pass@tcp(localhost:3306)/sales_manager?parseTime=true)")
	}

	firebaseKeyPath := os.Getenv("FIREBASE_JSON_PATH")
	if firebaseKeyPath == "" {
		log.Fatal("Falta la variable de entorno FIREBASE_JSON_PATH (Ruta al archivo .json que bajaste de Firebase)")
	}

	// 1. Conectar a MySQL
	fmt.Println("Conectando a MySQL...")
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("Error conectando a MySQL: %v", err)
	}

	// Correr las migraciones para asegurarnos de que las tablas existen
	fmt.Println("Migrando la base de datos (creando tablas si no existen)...")
	err = db.AutoMigrate(&tenant.Tenant{}, &user.User{})
	if err != nil {
		log.Fatalf("Error corriendo las migraciones: %v", err)
	}

	// 2. Conectar a Firebase
	fmt.Println("Conectando a Firebase...")
	opt := option.WithCredentialsFile(firebaseKeyPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error inicializando Firebase: %v", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error obteniendo cliente de Auth: %v", err)
	}

	// 3. Crear o buscar usuario en Firebase
	var uid string
	fbUser, err := client.GetUserByEmail(context.Background(), email)
	if err != nil {
		fmt.Println("Usuario no existe en Firebase. Creándolo...")
		params := (&auth.UserToCreate{}).
			Email(email).
			EmailVerified(true).
			DisplayName("Sebastian Fuchilieri (SuperAdmin)")
		
		newUser, err := client.CreateUser(context.Background(), params)
		if err != nil {
			log.Fatalf("Error creando usuario en Firebase: %v", err)
		}
		uid = newUser.UID
	} else {
		fmt.Println("Usuario ya existe en Firebase.")
		uid = fbUser.UID
	}

	// 4. Asegurar que exista el Tenant de Sistema (SuperAdmin) en la Base de Datos
	var t tenant.Tenant
	if err := db.FirstOrCreate(&t, tenant.Tenant{Email: "admin@lila.com.ar"}).Error; err != nil {
		// En caso de que se cree, actualizamos los datos
		db.Model(&t).Updates(tenant.Tenant{
			Name:   "Lila Software Studio",
			Owner:  "Sebastian Fuchilieri",
			Status: tenant.EstadoActivo,
			Plan:   tenant.PlanPro,
		})
	}

	// 5. Crear el Usuario SuperAdmin en MySQL
	var dbUser user.User
	result := db.Where("email = ?", email).First(&dbUser)
	if result.Error != nil {
		fmt.Println("Creando usuario SuperAdmin en MySQL...")
		dbUser = user.User{
			TenantID:           t.ID,
			Name:               "Sebastian Fuchilieri",
			Email:              email,
			TeamRole:           user.EquipoRolTitular,
			Role:               user.RoleSuperAdmin,
			Status:             user.EstadoActivo,
			FirebaseUID:        uid,
			SplitPercentageSub: 0,
		}
		if err := db.Create(&dbUser).Error; err != nil {
			log.Fatalf("Error creando usuario en MySQL: %v", err)
		}
	} else {
		fmt.Println("Usuario ya existe en MySQL. Actualizando rol a SuperAdmin y vinculando FirebaseUID...")
		db.Model(&dbUser).Updates(map[string]interface{}{
			"role":         user.RoleSuperAdmin,
			"firebase_uid": uid,
		})
	}

	// 6. Asignar Claims (Permisos) en Firebase
	fmt.Println("Asignando claims de SuperAdmin en Firebase...")
	claims := map[string]interface{}{
		"tenantId": float64(t.ID),
		"userId":   float64(dbUser.ID),
		"role":     string(user.RoleSuperAdmin),
	}
	
	if err := client.SetCustomUserClaims(context.Background(), uid, claims); err != nil {
		log.Fatalf("Error asignando claims: %v", err)
	}

	fmt.Println("✅ ¡LISTO! El usuario", email, "fue configurado exitosamente como SuperAdmin.")
	fmt.Println("Ya podés entrar a la web y hacer click en 'Ingresar con Google'.")
}
