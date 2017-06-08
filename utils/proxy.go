package utils

// func EncryptThenProxy(dst io.Writer, src io.Reader) (written int64, err error) {
// 	buf := make([]byte, 64*1024-100)
// 	encryptor := NewEncryptor("test")
// 	for {
// 		nr, er := src.Read(buf)
// 		if nr > 0 {
// 			nw, ew := EncryptThenWrite(dst, buf[:nr], encryptor)
//
// 			if nw > 0 {
// 				written += int64(nw)
// 			}
//
// 			if ew != nil {
// 				err = ew
// 				break
// 			}
// 		}
//
// 		if er != nil {
// 			if er != io.EOF {
// 				err = er
// 			}
// 			break
// 		}
// 	}
// 	return written, err
// }

// func EncryptThenWrite(dst io.Writer, rawData []byte, encryptor *Encryptor) (nw int, err error) {
// 	encryptBytes := encryptor.CFBEncrypter(rawData)
// 	encryptBytesLength := len(encryptBytes)
//
// 	encryptBytesLengthBytes := make([]byte, 2)
// 	binary.BigEndian.PutUint16(encryptBytesLengthBytes, uint16(encryptBytesLength))
// 	encryptBytes = append(encryptBytesLengthBytes, encryptBytes...)
// 	encryptBytesLength += 2
//
// 	// fmt.Printf("[debug] raw: %v, encrypted: %v\n", rawData, encryptBytes)
//
// 	nw, ew := dst.Write(encryptBytes)
//
// 	if ew != nil {
// 		err = ew
// 		return
// 	}
//
// 	if encryptBytesLength != nw {
// 		err = io.ErrShortWrite
// 		return
// 	}
//
// 	return
// }
//
// func ReadThenDecrypt(src io.Reader, buf []byte, encryptor *Encryptor) (decryptedBytes []byte, err error) {
// 	if buf == nil {
// 		buf = make([]byte, 64*1024)
// 	}
//
// 	if _, er := io.ReadFull(src, buf[:2]); er != nil {
// 		if er == io.EOF {
// 			return
// 		}
// 		println("read encrytped data length error")
// 		err = er
// 		return
// 	}
//
// 	dataLen := binary.BigEndian.Uint16(buf[:2])
// 	_, er := io.ReadFull(src, buf[:dataLen])
// 	if er != nil {
// 		if er == io.EOF {
// 			return
// 		}
// 		println("can't read full encrytped data")
// 		err = er
// 		return
// 	}
// 	// fmt.Printf("[debug] encrypted: %v, ", buf[:dataLen])
// 	decryptedBytes = encryptor.CFBDecrypter(buf[:dataLen])
// 	// fmt.Printf("raw: %v\n", decryptedBytes)
// 	fmt.Printf("%v %s", decryptedBytes, decryptedBytes)
// 	// decryptedBytesLength := len(decryptedBytes)
// 	return
// }

// func ProxyThenDecrypt(dst io.Writer, src io.Reader) (written int64, err error) {
// 	buf := make([]byte, 64*1024)
// 	encryptor := NewEncryptor("test")
//
// 	for {
// 		decryptedBytes, er := ReadThenDecrypt(src, buf, encryptor)
//
// 		if er != nil {
// 			println(er)
// 			break
// 		}
// 		// decryptedBytes := encryptor.CFBDecrypter(buf[:dataLen])
// 		decryptedBytesLength := len(decryptedBytes)
//
// 		// fmt.Printf("decrypt %X to %X", buf[0:nr], decryptedBytes)
// 		nw, ew := dst.Write(decryptedBytes)
//
// 		if nw > 0 {
// 			written += int64(nw)
// 		}
//
// 		if ew != nil {
// 			err = ew
// 			break
// 		}
//
// 		if decryptedBytesLength != nw {
// 			err = io.ErrShortWrite
// 			break
// 		}
// 	}
// 	return written, err
// }
