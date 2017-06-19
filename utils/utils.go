package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"

	"golang.org/x/crypto/salsa20"
	"golang.org/x/crypto/scrypt"
)

// HashScrypt Genera un hash a partir de una contraseña y un salt
func HashScrypt(pass, salt []byte) ([]byte, error) {

	key, err := scrypt.Key(pass, salt, 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("%x", key)), nil
}

// HashSha512 Genera un hash a partir de una contraseña
func HashSha512(pass []byte) [64]byte {
	return sha512.Sum512([]byte(pass))
}

// GenerateRandomBytes Genera cadenas aleatorias
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// err == nil only if len(b) == n
	if err != nil {
		return nil, err
	}

	return b, nil
}

// CryptoRandSecure Genera enteros aleatorios
func CryptoRandSecure(max int64) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
	}
	return int(nBig.Int64())
}

// GeneratePassword generador de contraseñas
func GeneratePassword(tamano int, letras bool, numeros bool, simbolos bool) string {
	arrayLetras := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	arrayNumeros := "0123456789"
	arraySimbolos := "+-/&<%#@*{!"

	arrayBase := ""
	salida := ""

	if letras == true {
		arrayBase = arrayBase + arrayLetras
	}
	if numeros == true {
		arrayBase = arrayBase + arrayNumeros
	}
	if simbolos == true {
		arrayBase = arrayBase + arraySimbolos
	}

	if arrayBase != "" {
		for i := 0; i < tamano; i++ {
			salida += string(arrayBase[CryptoRandSecure(int64(len(arrayBase)))])
		}
	} else {
		return "error"
	}

	return string(salida)
}

// EncodeBase64 función para codificar de []bytes a string (Base64)
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // sólo utiliza caracteres "imprimibles"
}

// DecodeBase64 función para decodificar de string a []bytes (Base64)
func DecodeBase64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) // recupera el formato original
	if err != nil {                              // comprobamos el error
		panic(err)
	}
	return b // devolvemos los datos originales
}

// ZLibCompress función para comprimir usando la libraria Zlib
func ZLibCompress(text string) string {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write([]byte(text))
	w.Close()
	return string(b.Bytes())
	// Output: [120 156 202 72 205 201 201 215 81 40 207 47 202 73 225 2 4 0 0 255 255 33 231 4 147]
}

// ZLibCompress función para descomprimir usando la libraria Zlib
func ZLibDecompress(buff []byte) string {
	b := bytes.NewReader(buff)
	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	if b, err := ioutil.ReadAll(r); err == nil {
		return string(b)
	}
	return "error"
}

// EncryptAES función para cifrar con AES, adjunta el IV al principio
func EncryptAES(data, key []byte) (out []byte) {
	out = make([]byte, len(data)+16) // reservamos espacio para el IV al principio
	rand.Read(out[:16])              // generamos el IV
	blk, err := aes.NewCipher(key)   // cifrador en bloque (AES), usa key
	if err != nil {                  // comprobamos el error
		panic(err)
	}
	ctr := cipher.NewCTR(blk, out[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out[16:], data)    // ciframos los datos
	return
}

// DecryptAES función para descifrar con AES
func DecryptAES(data, key []byte) (out []byte) {
	out = make([]byte, len(data)-16) // la salida no va a tener el IV
	blk, err := aes.NewCipher(key)   // cifrador en bloque (AES), usa key
	if err != nil {                  // comprobamos el error
		panic(err)
	}
	ctr := cipher.NewCTR(blk, data[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out, data[16:])     // desciframos (doble cifrado) los datos
	return
}

// CipherSalsa20 función para cifrar con Salsa20
func CipherSalsa20(dataIN []byte, key []byte, nonceIN []byte) (out []byte) {

	out = make([]byte, len(dataIN))

	// Requisitos de tamaño de nonce
	nonce := HashSha512(nonceIN)
	subnonce := nonce[0:24]

	// Requisitos de tamaño de key
	var subKey [32]byte
	copy(subKey[:], key[0:32])

	salsa20.XORKeyStream(out, dataIN, subnonce, &subKey)
	return
}
