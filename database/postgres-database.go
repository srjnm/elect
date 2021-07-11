package database

import (
	"crypto/sha256"
	"elect/models"
	"elect/roles"
	"encoding/hex"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type postgresDatabase struct {
	connection *gorm.DB
}

func NewPostgresDatabase() (Database, *http.ServeMux) {
	source := os.Getenv("DATABASE_URL")
	db, err := gorm.Open("postgres", source)
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.User{}, &models.Election{}, &models.Participant{}, &models.Blacklist{}, &models.Candidate{}, &models.ResetToken{})

	count := 0
	if db.Model(models.User{}).Where("email = ?", os.Getenv("ADMIN_EMAIL")).Count(&count); count == 0 {
		hashedPassword, err := HashPassword(os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			panic("Failed to initialize database!")
		}

		t := sha256.Sum256([]byte(strconv.Itoa(rand.Int()) + os.Getenv("ADMIN_EMAIL") + strconv.Itoa(rand.Int())))
		token := t[:]

		ret := db.Create(&models.User{
			FirstName:   "Suraj",
			LastName:    "N M",
			Email:       os.Getenv("ADMIN_EMAIL"),
			Password:    hashedPassword,
			Role:        roles.SuperAdmin,
			VerifyToken: hex.EncodeToString(token),
			Verified:    true})
		if ret.Error != nil {
			panic(ret.Error.Error())
		}
	}

	mux := SetUpQORAdmin(db)

	return &postgresDatabase{
		connection: db,
	}, mux
}
