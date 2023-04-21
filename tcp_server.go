package main

import (
	"io"
	"log"
	"net"
	"sync"
	"strings"
"strconv"

	"github.com/issaalmusawi/repo3-crypt/mycrypt"
	conv"github.com/issaalmusawi/funtest/konv"
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.4:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
go func() {

	defer wg.Done()
	for {
	log.Println("før server.Accept() kallet")
	conn, err := server.Accept()
	if err != nil {
	return
	}

go func(c net.Conn) {
	defer c.Close()
	for {
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		if err != io.EOF {
				log.Println(err)
		}
			return // fra for løkke
			}
	
	encryptedMessage := []rune(string(buf[:n]))
	decryptedMessage, err := mycrypt.Krypter(encryptedMessage, -4)
	if err != nil{
		log.Fatal(err)
	}

	switch msg := string(decryptedMessage); msg {
		case "ping":
		encryptedResponse, err := mycrypt.Krypter([]rune("pong"), 4)
		if err != nil{
			log.Fatal(err)
		}

		_, err = c.Write([]byte(string(encryptedResponse)))
		if err != nil {
			log.Fatal(err)
		}
		case "Kjevik":
		line := "Kjevik;SN39040;18.03.2022 01:50;6"
		parts := strings.Split(line, ";")
		temperatureC, err := strconv.ParseFloat(parts[len(parts)-1], 64)
		if err != nil{
			log.Fatal(err)
		}
		temperatureF := conv.CelsiusToFahrenheit(temperatureC)

		encryptedResponse, err := mycrypt.Krypter([]rune(strconv.FormatFloat(temperatureF, 'f', 2, 64)), 4)
		if err != nil {
			log.Fatal(err)
		}

		_, err = c.Write([]byte(string(encryptedResponse)))
		if err != nil{
			log.Println(err)
			return
		}

	default:
		decryptedMsg, err := mycrypt.Krypter([]rune(msg), -4)
		if err != nil {
			log.Fatal(err)
		}

	//	_, err := []byte(string(decryptedMsg))
		_, err = c.Write([]byte(string(decryptedMsg)))

	//	_, err = c.Write([]byte(string(decryptedMsg)))
		if err != nil {
			log.Println(err)
			return
				}
			}
		}
		}(conn)
		}
	}()
	wg.Wait()
}

