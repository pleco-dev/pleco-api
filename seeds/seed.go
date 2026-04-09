package seeds

import (
	"fmt"
	permissionModule "go-auth-app/modules/permission"
	roleModule "go-auth-app/modules/role"
	userModule "go-auth-app/modules/user"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Helper to check if db is nil and log fatal if so
func mustHaveDB(db *gorm.DB) {
	if db == nil {
		log.Fatal("DB connection failed: db is nil")
	}
}

func SeedRoles(db *gorm.DB) map[string]roleModule.Role {
	mustHaveDB(db)

	roleNames := []string{"superadmin", "admin", "user"}
	roleMap := make(map[string]roleModule.Role)

	for _, name := range roleNames {
		role := roleModule.Role{Name: name}
		if err := db.FirstOrCreate(&role, roleModule.Role{Name: name}).Error; err != nil {
			log.Printf("Failed to seed role %s: %v", name, err)
			continue
		}
		roleMap[name] = role
	}
	fmt.Println("Roles seeding done")
	return roleMap
}

func SeedPermissions(db *gorm.DB) map[string]permissionModule.Permission {
	mustHaveDB(db)

	permNames := []string{"create_user", "delete_user"}
	permMap := make(map[string]permissionModule.Permission)

	for _, name := range permNames {
		perm := permissionModule.Permission{Name: name}
		if err := db.FirstOrCreate(&perm, permissionModule.Permission{Name: name}).Error; err != nil {
			log.Printf("Failed to seed permission %s: %v", name, err)
			continue
		}
		permMap[name] = perm
	}
	fmt.Println("Permissions seeding done")
	return permMap
}

func SeedAdmin(db *gorm.DB) {
	mustHaveDB(db)

	var user userModule.User

	roleMap := SeedRoles(db)
	superadminRole, ok := roleMap["superadmin"]
	if !ok {
		log.Println("Superadmin role not found: cannot seed admin user")
		return
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Println("ADMIN_EMAIL or ADMIN_PASSWORD environment variables not set, skipping admin seeding")
		return
	}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Error checking for admin user: %v", err)
		return
	}
	if user.ID == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			log.Printf("Error hashing admin password: %v", err)
			return
		}

		admin := userModule.User{
			Name:       "Super Admin",
			Email:      email,
			Password:   string(hashedPassword),
			RoleID:     superadminRole.ID,
			Role:       superadminRole.Name,
			IsVerified: true,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Error creating admin user: %v", err)
			return
		} else {
			log.Printf("Super Admin user seeded with email: %s", email)
		}
	}
	fmt.Println("Admin seeding done")
}
