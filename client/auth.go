package client

import (
	"log"

	"github.com/davidqhr/socccks/utils"
)

func authentication(client *Client) bool {
	conn := client.Conn
	switch client.AuthMethod {
	case utils.AuthNo:
		return true
	case utils.AuthUsernamePassword:
		// http://www.rfc-base.org/txt/rfc-1929.txt
		buf := utils.Pool33K.Get()
		defer utils.Pool33K.Put(buf)

		_, err := conn.Read(buf)

		if err != nil {
			log.Printf("read AUTH_USERNAME_PASSWORD error")
			return false
		}

		usernameLen := int(buf[1])

		if usernameLen > 255 || usernameLen < 1 {
			log.Printf("username size error [1-255]")
			return false
		}

		username := buf[1+1 : 1+1+usernameLen]

		passwordLen := int(buf[1+1+usernameLen])

		if passwordLen > 255 || passwordLen < 1 {
			log.Printf("password size error [1-255]")
			return false
		}

		password := buf[1+1+usernameLen+1 : 1+1+usernameLen+1+passwordLen]

		log.Printf("username: %s, password: %s", string(username), string(password))

		// TODO real auth
		err = client.AuthSuccess()

		if err != nil {
			return false
		}

		return true
	}

	return false
}
