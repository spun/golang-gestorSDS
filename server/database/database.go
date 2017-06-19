package database

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bertus193/gestorSDS/config"
	"github.com/bertus193/gestorSDS/model"
	"github.com/bertus193/gestorSDS/utils"
)

// Base de datos de la aplicación
var gestor = make(map[string]*model.Usuario)

func init() {
	// Leer el fichero de la base de datos
	before()
}

// Descomprime y descifra el fichero de la base de datos
func before() {
	result := make(map[string]*model.Usuario)

	bytesEntrada, err := ioutil.ReadFile("./server/database/bd.txt")
	error := false
	if err != nil {
		error = true
	}
	if error == true || len(string(bytesEntrada)) == 0 {
		//fileData := []byte("{}")
		ioutil.WriteFile("./server/database/bd.txt", []byte(""), 0644)
	} else {
		decompress := []byte(utils.ZLibDecompress(bytesEntrada))
		decompress = utils.DecryptAES(decompress, config.PassDBEncrypt)
		if err := json.Unmarshal(decompress, &result); err != nil {
			panic("Error al leer fichero de entrada")
		}
	}

	gestor = result
}

// CreateUser guarda un nuevo usuario en la BD
func CreateUser(email string, passw string) error {

	var errResult error

	// Comprobamos si existe el email en la BD
	if _, ok := gestor[email]; ok {
		// Si existe el email, no modificamos nada
		errResult = errors.New("user already exists")
	} else if salt, errSalt := utils.GenerateRandomBytes(64); errSalt != nil {
		// Error al generar "salt"
		errResult = errors.New("unable to save")
	} else {
		// Hash de la contraseña también en servidor
		bytePass := []byte(passw)
		hashPass, _ := utils.HashScrypt(bytePass, salt)
		saltBase64 := utils.EncodeBase64(salt)

		// Guardamos el nuevo usuario
		gestor[email] = &model.Usuario{
			UserPassword:     string(hashPass),
			UserPasswordSalt: saltBase64,
			A2FEnabled:       false,
			Vault:            make(map[string]model.VaultEntry)}
	}
	return errResult
}

// ReadUser el usuario indicado a partir del email
func ReadUser(email string) (*model.Usuario, error) {

	var userResult *model.Usuario
	var errResult error

	if user, okUser := gestor[email]; !okUser {
		// Si no existe el el usuario indicado, no modificamos nada
		errResult = errors.New("user not found")
	} else {
		userResult = user
	}

	return userResult, errResult
}

// GetUser recupera un usuario de la BD que contenta el mismo
// email y contraseña que las indicads
func GetUser(email string, passw string) (*model.Usuario, error) {

	var userResult *model.Usuario
	var errResult error

	// Comprobamos si existe el email en la BD
	if user, ok := gestor[email]; !ok {
		// Si no existe el el usuario indicado
		errResult = errors.New("user not found")
	} else if salt, errSalt := base64.StdEncoding.DecodeString(user.UserPasswordSalt); errSalt != nil {
		// Error al recuperar el "salt"
		errResult = errors.New("unable to recover")
	} else {
		// Regeneramos el hash de servidor de la contraseña
		bytePass := []byte(passw)
		if hashPass, errHash := utils.HashScrypt(bytePass, salt); errHash != nil {
			// Error al regenerar el hash
			errResult = errors.New("unable to recover")
		} else if user.UserPassword != string(hashPass) {
			// Las contraseñas no coinciden
			errResult = errors.New("passwords do not match")
		} else {
			userResult = user
		}
	}
	return userResult, errResult
}

// CreateTextVaultEntry crea una entrada de tipo texto en el usaurio
func CreateTextVaultEntry(email string, entryTitle string, entryText string) error {
	var errResult error

	if user, okUser := gestor[email]; !okUser {
		// Si no existe el el usuario indicado, no modificamos nada
		errResult = errors.New("user not found")
	} else if _, okEntry := user.Vault[entryTitle]; okEntry {
		// Si ya existe una entrada con el mismo título
		errResult = errors.New("entry already exists")
	} else {
		user.Vault[entryTitle] = model.VaultEntry{
			Mode: 0, // Text
			Text: entryText,
		}
	}

	return errResult
}

// CreateAccountVaultEntry crea una entrada de tipo cuenta en el usaurio
func CreateAccountVaultEntry(email string, entryTitle string, userAccount string, passwAccount string) error {
	var errResult error

	if user, okUser := gestor[email]; !okUser {
		// Si no existe el el usuario indicado, no modificamos nada
		errResult = errors.New("user not found")
	} else if _, okEntry := user.Vault[entryTitle]; okEntry {
		// Si ya existe una entrada con el mismo título
		errResult = errors.New("entry already exists")
	} else {
		user.Vault[entryTitle] = model.VaultEntry{
			Mode:     1, // Account
			User:     userAccount,
			Password: passwAccount,
		}
	}

	return errResult
}

// ReadVaultEntry recupera la lista de entradas (sin detalles)
// de un usuario
func ReadVaultEntry(email string, entryTitle string) (model.VaultEntry, error) {

	var entryResult model.VaultEntry
	var errResult error

	if user, okUser := gestor[email]; !okUser {
		// Si no existe el el usuario indicado, no modificamos nada
		errResult = errors.New("user not found")
	} else if entry, okEntry := user.Vault[entryTitle]; !okEntry {
		// Si no existe una entrada con el mismo título
		errResult = errors.New("entry not found")
	} else {
		entryResult = entry
	}

	return entryResult, errResult
}

// DeleteVaultEntry eliina una entrada concreta del usuario
func DeleteVaultEntry(email string, entryTitle string) error {

	var errResult error

	if user, okUser := gestor[email]; !okUser {
		// Si no existe el el usuario indicado, no modificamos nada
		errResult = errors.New("user not found")
	} else if _, okEntry := user.Vault[entryTitle]; !okEntry {
		// Si no existe una entrada con el mismo título
		errResult = errors.New("entry not found")
	} else {
		delete(gestor[email].Vault, entryTitle)
	}

	return errResult
}

// UpdateA2F cambia el estado de activación de A2F para el usuario
func UpdateA2F(email string, newState bool) error {

	var errResult error

	if user, okUser := gestor[email]; !okUser {
		// Si no existe el el usuario indicado, no modificamos nada
		errResult = errors.New("user not found")
	} else {
		user.A2FEnabled = newState
	}

	return errResult
}

// DeleteUser Elimina cuenta de usuario
func DeleteUser(email string) error {

	var errResult error

	if _, okUser := gestor[email]; !okUser {
		// Si no existe el el usuario indicado, no modificamos nada
		errResult = errors.New("user not found")
	} else {
		delete(gestor, email)
	}

	return errResult

}

// After Persistencia Base de Datos
func After() {
	salida, err := os.Create("./server/database/bd.txt")
	if err != nil {
		panic(0)
	}

	// todo: comprobar y validar contraseña

	j, err := json.Marshal(gestor)

	if err != nil {
		fmt.Println(err)
	}

	//usuarios := string(j)
	usuarios := string(utils.EncryptAES(j, config.PassDBEncrypt)) //Encriptar
	usuarios = utils.ZLibCompress(usuarios)                       //Comprimir

	salida.Write([]byte(usuarios))
}
