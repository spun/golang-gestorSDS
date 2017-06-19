package config

// AppName contiene el nombre de la aplicación
var AppName = "Gestor SDS"

// SecureServerPort Puerto Seguro Cliente
var SecureServerPort = ":10443"

// SecureURL Url segura
var SecureURL = "https://127.0.0.1"

// MaxTimeSession es el tiempo máximo de sesión (segundos)
var MaxTimeSession = 60 * 30

// MaxA2FTime es el tiempo máximo de espera para resolver
// el reto de segundo factor de autenticación (segundos)
var MaxA2FTime = 60 * 5

// SizeA2FCode es el número de digitos que contendrá la clave
// que se envía al los usuario con A2F activado
var SizeA2FCode = 6

// PassDBEncrypt es la clave de cifrado del fichero de base de datos
var PassDBEncrypt = []byte("a very very very very secret key")

// Account2FA contiene los datos de la cuenta de correo
// encargada de enviar los códigos de inicio de sesión
var Account2FA = map[string]string{
	"email":      "",
	"passw":      "",
	"smtpServer": "smtp.gmail.com",
	"smtpPort":   "587",
}

// EmailDebug permite comprobar las funcionalidades que
// hacen uso de correos electrónicos sin realizar el envío.
// Una vez puesto a "true", los correos se mostrarán por la salida estandar.
var EmailDebug = true

// EncryptLogs se encarga de indicar si se desea cifrar el log del servidor
var EncryptLogs = true

// PassEncryptLogs Clave de cifrado de los ficheros de logs
var PassEncryptLogs = []byte("a really difficult logg password")
