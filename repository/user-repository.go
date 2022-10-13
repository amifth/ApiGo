package repository

import (
	"log"

	"github.com/amifth/gorest/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) entity.User
	FetchUserById(userID int) entity.User
	FindAll() []entity.User
	DeleteById(userId int) entity.User
}

type userConnecttion struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnecttion{
		connection: db,
	}
}

func (db *userConnecttion) InsertUser(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))
	db.connection.Save(&user)
	return user
}

func (db *userConnecttion) UpdateUser(user entity.User) entity.User {
	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tempUser entity.User
		db.connection.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	}

	db.connection.Save(&user)
	return user
}

func (db *userConnecttion) VerifyCredential(email string, _ string) interface{} {
	var user entity.User
	res := db.connection.Where("email = ?", email).Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (db *userConnecttion) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User
	return db.connection.Where("email = ?", email).Take(&user)
}

func (db *userConnecttion) FindByEmail(email string) entity.User {
	var user entity.User
	db.connection.Where("email = ?", email).Take(&user)
	return user
}

func (db *userConnecttion) FetchUserById(userID int) entity.User {
	var user entity.User
	db.connection.Find(&user, userID)
	return user
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash password")
	}
	return string(hash)
}

func (db *userConnecttion) FindAll() []entity.User {
	var users []entity.User
	db.connection.Find(&users)
	// fmt.Println(users)
	return users
}

func (db *userConnecttion) DeleteById(userId int) entity.User {
	var user entity.User
	res := db.FetchUserById(userId)
	db.connection.Delete(&user, userId)
	return res
}
