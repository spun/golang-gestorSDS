package model

import "time"

/* ----------- DATABASE ----------- */

type Usuario struct {
	UserPassword     string
	UserPasswordSalt string
	A2FEnabled       bool
	Vault            map[string]VaultEntry
}

type VaultEntry struct {
	Mode int
	// Mode 0 - Plain text
	Text string
	// Mode 1 - Account
	User     string
	Password string
}

/* Demo estructura en json (sin cifrados)
   "alu@alu.ua.es": {
       "UserPassword": "accoutPass",
       "UserPasswordSalt": "accoutSalrPass",
       "Vault": {
           "memoria": {
               "Mode": "0",
               "Text": "texto de la entrada"
           },
           "twitter": {
               "Mode": "1",
               "User": "usuarioTwitter",
               "Password": "54321"
           }
       }
   }
*/

/* -------------------------------- */

/*  ----- USUARIO ACTIVO ----- */
type ActiveUser struct {
	UserEmail          string
	SesssionExpireTime time.Time

	A2FResolved   bool
	A2FChallenge  string
	A2FExpiration time.Time
}

/*  -------------------------- */

/*  ----- PETICIONES ----- */

type DetallesUsuario struct {
	Email      string
	A2FEnabled bool
	NumEntries int
}

type ListaEntradas struct {
	Texts    []string
	Accounts []string
}

/* ----------------------- */
