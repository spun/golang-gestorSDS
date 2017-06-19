package client

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bertus193/gestorSDS/config"
	"github.com/bertus193/gestorSDS/model"
	"github.com/bertus193/gestorSDS/utils"
	"github.com/fatih/color"
)

var httpClient *http.Client

func startUI(c *http.Client) {
	httpClient = c
	uiInicio("", "")
}

// Pantalla de bienvenida con las opciones de
// login, registro y cerrar aplicación
func uiInicio(showError string, showSuccess string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Bienvenido a %s\n\n", config.AppName)

	// Mensaje de confirmación de acción en caso de existir
	if showSuccess != "" {
		color.HiGreen("* %s\n\n", showSuccess)
	}
	// Opciones
	fmt.Println("1. Entrar")
	fmt.Println("2. Crear usuario")
	fmt.Println("0. Salir")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("\n* %s", showError)
	}

	// Lectura de opción elegida
	fmt.Printf("\nSeleccione una opción: ")
	inputSelectionStr := utils.CustomScanf()

	// Ejecución de la opción elegida
	switch {
	case inputSelectionStr == "1": // Login
		uiLoginUser("")
	case inputSelectionStr == "2": // Registro
		uiRegistroUsuario("")
	case inputSelectionStr == "0": // Salir
		os.Exit(0)
	default:
		uiInicio("La opción elegida no es correcta", "")
	}
}

// Pantalla de creación de usuario
func uiRegistroUsuario(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Registro de usuarios\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("* %s\n\n", showError)
	}

	// Lectura de datos del nuevo usuario
	fmt.Print("Email: ")
	inputUser := utils.CustomScanf()
	fmt.Print("Contraseña: ")
	inputPass := utils.CustomScanf()

	// Petición al servidor
	if err := registroUsuario(httpClient, inputUser, inputPass); err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "user already exists":
			uiRegistroUsuario("Ya existe un usuario con ese correo.")
		default:
			uiRegistroUsuario("Ocurrio un error al crear el usuario.")
		}
	} else {
		// Registro completado, volvemos a la pantalla de inicio
		uiInicio("", "Registro completado, ya puedes iniciar sesión")
	}
}

// Pantalla de entrada de usuarios
func uiLoginUser(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Entrada de usuarios\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("* %s\n\n", showError)
	}

	// Lectura de datos del usuario
	fmt.Print("Email: ")
	inputUser := utils.CustomScanf()
	fmt.Print("Contraseña: ")
	inputPass := utils.CustomScanf()

	// Petición al servidor
	if err := loginUsuario(httpClient, inputUser, inputPass); err != nil {

		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "user not found":
			uiInicio("No exite ningún usuario con esos datos.", "")
		case "passwords do not match":
			// No damos información detallada del error en este caso
			uiInicio("No exite ningún usuario con esos datos.", "")
		case "a2f required":
			uiUnlockA2F("")
		default:
			uiInicio("Ocurrio un error al realizar el login.", "")
		}

	} else {
		// Login completado, vamos a la pantalla principal del usuario
		uiUserMainMenu("", "")
	}
}

// Pantalla de reto de auntenticación en dos pasos
func uiUnlockA2F(showError string) {
	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Verificación en 2 pasos\n\n")

	// Mensaje de información
	color.HiGreen("> Te hemos enviado un correo a tu cuenta con el código que debes introducir para iniciar sesión.\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("* %s\n\n", showError)
	}

	// Lectura de opción elegida
	fmt.Print("Código: ")
	inputA2Fcode := utils.CustomScanf()

	// Petición al servidor
	if err := desbloquearA2F(httpClient, inputA2Fcode); err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "session not found":
			uiLoginUser("La sesión de usuario ha cadudado.")
		case "2fa already resolved":
			uiUserMainMenu("", "")
		case "2fa expired":
			uiLoginUser("El código de verificación en dos pasos ha caducado.")
		case "incorrect 2fa code":
			uiUnlockA2F("El código introducido no es valido.")
		default:
			uiUnlockA2F("Ocurrio un error verificar el código.")
		}
	} else {
		//Desbloqueado con exito
		uiUserMainMenu("", "")
	}
}

// Pantalla principal del usuario, listado de entradas
func uiUserMainMenu(showError string, showSuccess string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Página de usuario\n")

	// Mensaje de confirmación de acción en caso de existir
	if showSuccess != "" {
		color.HiGreen("\n* %s\n", showSuccess)
	}

	// Recuperamos las cuentas del usuario
	fmt.Printf("\n------ Listado de entradas ------\n\n")
	// Petición al servidor
	entradas, err := listarEntradas(httpClient)
	if err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("La sesión de usuario ha cadudado.")
		default:
			fmt.Println("Ocurrio un error al recuperar las entradas." + err.Error())
		}
	} else {
		// Mostramos la lista de cuentas de usuario guardadas

		boldBlue := color.New(color.FgHiBlue, color.Bold)
		if (entradas.Accounts != nil && len(entradas.Accounts) != 0) ||
			(entradas.Texts != nil && len(entradas.Texts) != 0) {

			// Mostramos la lista de cuentas de usuario guardadas
			if entradas.Accounts != nil && len(entradas.Accounts) != 0 {
				boldBlue.Printf(" Cuentas de usuario\n")
				// Imprimimos los resultados
				for c := range entradas.Accounts {
					fmt.Printf("    [%s]\n", entradas.Accounts[c])
				}
				fmt.Printf("\n")
			}

			// Mostramos la lista de textos guardados
			if entradas.Texts != nil && len(entradas.Texts) != 0 {
				boldBlue.Printf(" Textos guardados\n")
				// Imprimimos los resultados
				for c := range entradas.Texts {
					fmt.Printf("    [%s]\n", entradas.Texts[c])
				}
			}

		} else {
			boldBlue.Printf("* No hay nada guardado todavía\n")
		}
	}
	fmt.Printf("\n--------------------------------\n\n")

	// Opciones
	fmt.Println("1. Añadir entrada")
	fmt.Println("2. Ver detalle de entrada")
	fmt.Println("3. Configuración de mi cuenta")
	fmt.Println("0. Salir")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("\n* %s", showError)
	}

	// Lectura de opción elegida
	fmt.Printf("\nSeleccione una opción: ")
	inputSelectionStr := utils.CustomScanf()

	// Ejecución de la opción elegida
	switch {
	case inputSelectionStr == "1":
		uiAddNewEntry("")
	case inputSelectionStr == "2":
		fmt.Print("Escribe el nombre de la entrada: ")
		inputEntrySelectionStr := utils.CustomScanf()
		uiDetailsEntry("", inputEntrySelectionStr)
	case inputSelectionStr == "3":
		uiUserConfiguration("")
	case inputSelectionStr == "0":
		uiInicio("", "")
	default:
		uiUserMainMenu("La opción elegida no es correcta", "")
	}
}

// Pantalla de creación de nueva entrada
func uiAddNewEntry(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Añadir nueva entrada\n\n")

	// Solicitamos información de lo que queremos guardar de entre las posibles
	fmt.Println("1. Texto")
	fmt.Println("2. Cuenta de usuario")
	fmt.Println("0. Volver")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("\n* %s\n", showError)
	}

	// Lectura de opción elegida
	fmt.Printf("\nSeleccione una opción: ")
	inputEntryMode := utils.CustomScanf()

	switch inputEntryMode {
	case "1":
		uiAddNewTextEntry("")
	case "2":
		uiAddNewAccountEntry("")
	case "0":
		uiUserMainMenu("", "")
	default:
		uiAddNewEntry("La opción elegida no es correcta.")
	}
}

// Pantalla de creación de nueva entrada de tipo texto
func uiAddNewTextEntry(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Añadir nuevo texto\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("* %s\n\n", showError)
	}

	// Lectura de los datos de la nueva entrada
	fmt.Printf("\nEscribe el título del texto: ")
	inputTitle := utils.CustomScanf()
	fmt.Printf("\nEscribe el texto que quieres guardar (ENTER para terminar):\n\n")
	inputText := utils.CustomScanf()

	// Petición al servidor
	if err := crearEntradaDeTexto(httpClient, inputTitle, inputText); err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("La sesión de usuario ha cadudado.")
		case "user not found":
			uiLoginUser("Ha ocurrido un error al guardar la entrada en tu cuenta.")
		case "entry already exists":
			uiUserMainMenu("Ya existe una entrada con ese título.", "")
		default:
			uiUserMainMenu("Ocurrio un error al añadir la entrada el código.", "")
		}
	} else {
		uiUserMainMenu("", "Texto ["+inputTitle+"] añadido correctamente")
	}
}

// Pantalla de creación de nueva entrada de tipo cuenta de usuario
func uiAddNewAccountEntry(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Añadir nueva cuenta de usuario\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("* %s\n\n", showError)
	}

	// Lectura de los datos de la nueva entrada
	fmt.Print("Título de la entrada (twitter, facebook, etc): ")
	inputAccountType := utils.CustomScanf()
	fmt.Print("Usuario: ")
	inputAccountUser := utils.CustomScanf()
	fmt.Print("¿Deseas generar una contraseña? (si, no): ")
	inputGeneratePassw := utils.CustomScanf()

	var finalPassw string
	if inputGeneratePassw == "si" || inputGeneratePassw == "s" {

		// Solicitamos información de como se desea generar la contraseña
		for {
			// Tamaño de la contraseña
			var genLenght int
			for {
				fmt.Print("¿Que tamaño de contraseña deseas? ")
				inputLenght := utils.CustomScanf()
				if convLenght, err := strconv.Atoi(inputLenght); err == nil {
					genLenght = convLenght
					break
				}
			}

			// La contraseña generada puede tener números
			fmt.Print("¿Deseas que tenga números? (si, no): ")
			inputWithNums := utils.CustomScanf()
			genWithNums := inputWithNums == "si" || inputWithNums == "s"

			// La contraseña generada puede tener simbolos
			fmt.Print("¿Deseas que tenga símbolos? (si, no): ")
			inputWithSymbols := utils.CustomScanf()
			genWithSymbols := inputWithSymbols == "si" || inputWithSymbols == "s"

			// Mostramos la contraseña y preguntamos al usuario si está de acuerdo
			finalPassw = utils.GeneratePassword(genLenght, true, genWithNums, genWithSymbols)
			fmt.Printf("La contraseña es: %s\n¿Estás de acuerdo? (si, no): ", finalPassw)
			inputConfirm := utils.CustomScanf()
			if inputConfirm == "si" || inputConfirm == "s" {
				break
			}
		}

	} else {
		fmt.Print("Contraseña: ")
		finalPassw = utils.CustomScanf()
	}

	// Petición al servidor
	if err := crearEntradaDeCuenta(httpClient, inputAccountType, inputAccountUser, finalPassw); err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("La sesión de usuario ha cadudado.")
		case "user not found":
			uiLoginUser("Ha ocurrido un error al guardar la entrada en tu cuenta.")
		case "entry already exists":
			uiUserMainMenu("Ya existe una entrada con ese título.", "")
		default:
			uiUserMainMenu("Ocurrio un error al añadir la entrada el código.", "")
		}
	} else {
		uiUserMainMenu("", "Cuenta ["+inputAccountType+"] añadida correctamente")
	}
}

// Pantalla de visualización de detalles de una entrada
func uiDetailsEntry(showError string, entryName string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Detalles de la entrada [%s]\n\n", entryName)

	// Petición al servidor
	fmt.Printf("--------------------------------\n\n")
	entry, err := detallesEntrada(httpClient, entryName)
	if err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("La sesión de usuario ha cadudado.")
		case "not found":
			uiUserMainMenu("No se han podido obtener detalles de la cuenta elegida.", "")
		default:
			fmt.Println("Ocurrió un error al recuperar las entradas." + err.Error())
		}
	} else {
		// Si los detalles de la cuenta están vacios
		if (model.VaultEntry{}) == entry {
			// Volvemos al menú del usuario
			uiUserMainMenu("No se han podido obtener detalles de la cuenta elegida.", "")
		}

		// Comprobamos el tipo de entrada (texto, cuenta) y la mostramos
		if entry.Mode == 0 {
			// Si es una entrada de tipo texto
			fmt.Printf("[Texto] \n\n%s\n", entry.Text)

		} else if entry.Mode == 1 {
			// Si es una entrada de tipo cuenta de usuario
			fmt.Printf("[Usuario] -> %s \n", entry.User)
			fmt.Printf("[Contraseña] -> %s \n", entry.Password)
		}
	}
	fmt.Printf("\n--------------------------------\n\n")

	// Opciones
	fmt.Println("1. Borrar entrada")
	fmt.Println("0. Volver")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("\n* %s", showError)
	}

	// Lectura de opción elegida
	fmt.Printf("\nSeleccione una opción: ")
	inputSelectionStr := utils.CustomScanf()

	switch {
	case inputSelectionStr == "1":
		fmt.Print("¿Estás seguro? (si, no): ")
		inputDecission := utils.CustomScanf()
		if inputDecission == "si" || inputDecission == "s" {

			// Petición al servidor para eliminar la entrada de la BD
			if errDel := eliminarEntrada(httpClient, entryName); errDel != nil {
				// Si hay un error, mostramos el mensaje de error adecuado
				switch errDel.Error() {
				case "unauthorized":
					uiLoginUser("La sesión de usuario ha cadudado.")
				case "not found":
					uiUserMainMenu("No se ha podido borrar la entrada.", "")
				default:
					fmt.Println("Ocurrió un error al borrar la entrada." + err.Error())
				}

			} else {
				// Se ha eliminado correctamente
				uiUserMainMenu("", "Entrada ["+entryName+"] eliminada correctamente")
			}
		} else {
			uiDetailsEntry("", entryName)
		}
	case inputSelectionStr == "0":
		uiUserMainMenu("", "")
	default:
		uiDetailsEntry("La opción elegida no es correcta", entryName)
	}
}

// Pantalla de visualización de detalles de un usuario
func uiUserConfiguration(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Página de configuración de usuario\n\n")

	// Petición al servidor
	fmt.Printf("------ Información de la cuenta ------\n\n")
	userDetails, err := detallesUsuario(httpClient)
	if err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("La sesión de usuario ha cadudado.")
		case "user not found":
			uiLoginUser("No se ha podido obtener la información del usuario.")
		case "unable to unmarshal":
			uiUserMainMenu("No se ha podido mostrar la información del usuario.", "")
		default:
			fmt.Println("Ocurrió un error al recuperar los detalles." + err.Error())
		}
	} else {

		// Mostramos los detalles de la cuenta
		boldBlue := color.New(color.FgHiBlue, color.Bold)
		fmt.Printf("Correo electrónico: ")
		boldBlue.Println(userDetails.Email)
		fmt.Printf("Número de entradas guardadas: ")
		boldBlue.Println(userDetails.NumEntries)
		fmt.Printf("Segundo factor de autenticación: ")
		if userDetails.A2FEnabled {
			boldGreen := color.New(color.FgBlack, color.BgHiGreen, color.Bold)
			boldGreen.Println(" Activado ")
		} else {
			boldRed := color.New(color.FgWhite, color.BgHiRed, color.Bold)
			boldRed.Println(" Desactivado ")
		}
	}
	fmt.Printf("\n------------------------------------\n\n")

	// Opciones
	fmt.Println("1. Modificar contraseña (to-do)")
	fmt.Println("2. Eliminar mi usuario")
	if userDetails.A2FEnabled {
		fmt.Println("3. Desactivar 2FA")
	} else {
		fmt.Println("3. Activar 2FA")
	}
	fmt.Println("0. Volver")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("\n* %s", showError)
	}

	// Lectura de opción elegida
	fmt.Printf("\nSeleccione una opción: ")
	inputSelectionStr := utils.CustomScanf()

	switch {
	case inputSelectionStr == "1":
		uiUserConfiguration("To-do")
	case inputSelectionStr == "2":
		fmt.Print("¿Estás seguro? (si, no): ")
		inputDecission := utils.CustomScanf()
		if inputDecission == "si" || inputDecission == "s" {

			// Petición al servidor para eliminar la entrada de la BD
			if errDel := eliminarUsuario(httpClient); errDel != nil {
				// Si hay un error, mostramos el mensaje de error adecuado
				switch errDel.Error() {
				case "unauthorized":
					uiLoginUser("La sesión de usuario ha cadudado.")
				case "user not found":
					uiLoginUser("La cuenta de usuario que desea borrar no existe.")
				default:
					fmt.Println("Ocurrió un error al borrar la entrada." + err.Error())
				}

			} else {
				// Se ha eliminado correctamente
				uiInicio("Tu cuenta de usuario se ha borrado correctamente. Para usar "+config.AppName+" debes crear un nuevo usuario.", "")
			}
		} else {
			uiUserConfiguration("")
		}
	case inputSelectionStr == "3":
		if errUpdate := updateA2F(httpClient, !userDetails.A2FEnabled); errUpdate != nil {
			// Si hay un error, mostramos el mensaje de error adecuado
			switch errUpdate.Error() {
			case "unauthorized":
				uiLoginUser("La sesión de usuario ha cadudado.")
			case "user not found":
				uiLoginUser("No se ha podido obtener la configuración.")
			default:
				uiUserMainMenu("No se ha podido cambiar la configuración.", "")
			}
		} else {
			uiUserConfiguration("")
		}
	case inputSelectionStr == "0":
		uiUserMainMenu("", "")
	default:
		uiUserConfiguration("La opción elegida no es correcta")
	}
}
