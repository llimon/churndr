package main

import (
	"encoding/base64"
	//"os"
	"fmt"
	"io/ioutil"
	"log"

	"crypto/tls"
	"net"
	"net/smtp"

	"github.com/llimon/page-on-pod-restarts/common"
)

func SendTLSMail(e common.Email) {

	// Connect to the SMTP Server
	servername := "smtp.gmail.com:465"

	content, err := ioutil.ReadFile(".password")
	if err != nil {
		log.Fatal(err)
	}

	str := content
	data, err := base64.StdEncoding.DecodeString(string(str))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", "luislimon@gmail.com", string(data), host)

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = e.From
	headers["To"] = e.To
	headers["Subject"] = e.Subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + e.Body

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(e.From); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(e.To); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

}
