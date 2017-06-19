package utils

import (
	"log"
	"net/smtp"

	"github.com/bertus193/gestorSDS/config"
)

// SendWelcome envía un correo electrónico de bienvenida
// a la dirección indicada (sin uso).
func SendWelcome(sendTo string) {

	subject := "Te damos la bienvenida a " + config.AppName
	body := "Gracias por crear una cuenta de " + config.AppName + ".\n" +
		"Ya puedes empezar a utilizar la aplicación usando tu correo y contraseña para iniciar sesión.\n\n" +
		"Gracias, \n" +
		"El equipo de cuentas de " + config.AppName + "."

	sendEmail(sendTo, subject, body)
}

// Send2FACode envía un correo electrónico a la dirección indicada
// con la información necesaria para iniciar sesión usando 2FA
func Send2FACode(sendTo string, authCode string) {

	subject := "Código de seguridad de inicio de sesión en " + config.AppName
	body := "Se ha realizado un inicio de sesión en su cuenta. \n\n" +
		"Use el siguiente código de seguridad para continuar. \n" +
		"Código de seguridad: " + authCode + "\n\n" +
		"Si no has iniciado sesión recientemente, es posible que su contraseña haya sido comprometida. \n" +
		"Pongase en contacto con nosotros lo antes posible para solucionarlo. \n\n" +
		"Gracias, \n" +
		"El equipo de cuentas de " + config.AppName + "."

	sendEmail(sendTo, subject, body)
}

func sendEmail(sendTo string, subject string, body string) {

	from := config.Account2FA["email"]
	pass := config.Account2FA["passw"]
	smtpServer := config.Account2FA["smtpServer"]
	smtpPort := config.Account2FA["smtpPort"]

	msg := "From: " + from + "\n" +
		"To: " + sendTo + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	if config.EmailDebug == true {
		// Mostramos en terminal
		log.Println(sendTo)
		log.Println(subject)
		log.Println(body)
	} else {
		// Enviamos el correo
		err := smtp.SendMail(smtpServer+":"+smtpPort,
			smtp.PlainAuth("", from, pass, smtpServer),
			from, []string{sendTo}, []byte(msg))

		if err != nil {
			log.Printf("smtp error: %s", err)
			return
		}
	}

	log.Print("email sent")
}
