package socks5

func authentication(client *Client) bool {
	conn := client.Conn
	switch client.AuthMethod {
	case AuthNo:
		return true
	case AuthUsernamePassword:
		// http://www.rfc-base.org/txt/rfc-1929.txt
		buf := bufferPool.Get()
		defer bufferPool.Put(buf)

		_, err := conn.Read(buf)

		if err != nil {
			logger.Info(client, "read AUTH_USERNAME_PASSWORD error")
			return false
		}

		usernameLen := int(buf[1])

		if usernameLen > 255 || usernameLen < 1 {
			logger.Error(client, "username size error [1-255]")
			return false
		}

		username := buf[1+1 : 1+1+usernameLen]

		passwordLen := int(buf[1+1+usernameLen])

		if passwordLen > 255 || passwordLen < 1 {
			logger.Error(client, "password size error [1-255]")
			return false
		}

		password := buf[1+1+usernameLen+1 : 1+1+usernameLen+1+passwordLen]

		logger.Info(client, "username: %s, password: %s", string(username), string(password))

		// TODO real auth
		err = client.AuthSuccess()

		if err != nil {
			return false
		}

		return true
	}

	return false
}
