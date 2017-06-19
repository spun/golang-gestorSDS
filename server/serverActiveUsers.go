package server

import (
	"errors"
	"strconv"
	"time"

	"github.com/bertus193/gestorSDS/config"
	"github.com/bertus193/gestorSDS/model"
	"github.com/bertus193/gestorSDS/utils"
)

var activeUsers = make(map[string]*model.ActiveUser)

// CreateUserSession añade un usuario a la lista de usuarios
// activos asignandole un token de sesión generado.
func CreateUserSession(userEmail string, useA2F bool) (string, string) {

	cleanAllInactiveUsers()

	var token = generateSessionToken()

	// Valores por defecto para un usuario sin A2F
	var a2fresolved = true
	var a2fchallenge = ""
	var a2fexpiration time.Time
	if useA2F {
		a2fresolved = false
		a2fchallenge = generateA2FCode()
		a2fexpiration = time.Now().Add(time.Second * time.Duration(config.MaxA2FTime))
	}

	activeUsers[token] = &model.ActiveUser{
		UserEmail:          userEmail,
		SesssionExpireTime: time.Now().Add(time.Second * time.Duration(config.MaxTimeSession)),
		A2FResolved:        a2fresolved,
		A2FChallenge:       a2fchallenge,
		A2FExpiration:      a2fexpiration,
	}

	return token, a2fchallenge
}

// UnlockSessionWith2FA desbloquea la sesión de un usuario con 2FA que
// ha introducido correctamente el código del reto.
func UnlockSessionWith2FA(token string, code2FA string) error {
	var err error
	if userSession, ok := activeUsers[token]; !ok {
		err = errors.New("session not found")
	} else if userSession.A2FResolved {
		err = errors.New("2fa already resolved")
	} else if time.Now().After(userSession.A2FExpiration) {
		err = errors.New("2fa expired")
	} else if userSession.A2FChallenge != code2FA {
		err = errors.New("incorrect 2fa code")
	} else {
		userSession.A2FResolved = true
		resetSessionExpireTime(token)
	}
	return err
}

// GetUserFromSession recupera el correo electrónico del usuario
// si está activo a partir del token de sesión que se indica.
func GetUserFromSession(token string) (string, error) {

	var userEmail = ""
	var err error
	if tempUser, ok := activeUsers[token]; !ok {
		err = errors.New("session not found")
	} else if isSessionExpired(token) {
		delete(activeUsers, token)
		err = errors.New("session expired")
	} else if !tempUser.A2FResolved {
		err = errors.New("2fa not resolved")
	} else {
		userEmail = tempUser.UserEmail
		resetSessionExpireTime(token)
	}

	return userEmail, err
}

func resetSessionExpireTime(token string) {
	if tempUser, ok := activeUsers[token]; ok {
		tempUser.SesssionExpireTime = time.Now().Add(time.Second * time.Duration(config.MaxTimeSession))
	}
}

// cleanInactiveUsers recorre la lista de usuarios activos y
// elimina aquellos cuya sesión haya cadudado.
func cleanAllInactiveUsers() {
	for k := range activeUsers {
		if isSessionExpired(k) {
			delete(activeUsers, k)
		}
	}
}

// isSessionExpired comprueba si el tiempo de validez de una
// sesión sigue estando activo
func isSessionExpired(token string) bool {
	isExpired := false
	if tempUser, ok := activeUsers[token]; ok {
		currentTime := time.Now()
		if currentTime.After(tempUser.SesssionExpireTime) {
			isExpired = true
		} else if tempUser.A2FResolved == false && currentTime.After(tempUser.A2FExpiration) {
			isExpired = true
		}
	} else {
		isExpired = true
	}
	return isExpired
}

// generateSessionToken genera un token para la sesión que será
// el que use el usuario al realizar las peticiones.
func generateSessionToken() string {
	tokenRaw, _ := utils.GenerateRandomBytes(24)
	tokenSrc := utils.EncodeBase64(tokenRaw)
	return tokenSrc
}

// generateA2FCode genera un código númerico unico que serivrá
// de reto para la autenticacion en dos pasos.
func generateA2FCode() string {
	result := ""
	codeSize := config.SizeA2FCode
	for i := 0; i < codeSize; i++ {
		gen := utils.CryptoRandSecure(10)
		result += strconv.Itoa(int(gen))
	}
	return result
}
