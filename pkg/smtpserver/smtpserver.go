package smtpserver

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net"
	"net/mail"
	"strconv"
	"strings"
	"suntek2telegram/pkg/config"
)

var (
	imagesReadersChan chan io.Reader

	usernameBase64 string
	passwordBase64 string
)

func Start(sc *config.SMTP, readersChan chan io.Reader) {
	imagesReadersChan = readersChan

	usernameBase64 = base64.StdEncoding.EncodeToString([]byte(sc.Username))
	passwordBase64 = base64.StdEncoding.EncodeToString([]byte(sc.Password))

	listener, err := net.Listen("tcp", sc.BindHost+":"+strconv.Itoa(sc.BindPort))
	if err != nil {
		log.Fatalln("Failed to start SMTP listener", err)
	}
	log.Println("SMTP TCP listener started on", sc.BindHost+":"+strconv.Itoa(sc.BindPort))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		log.Println("Connection established from", conn.RemoteAddr())

		go handleConn(conn)
	}
}

func sendAndReceive(writer *bufio.Writer, reader *bufio.Reader, message string, expected string) (string, error) {
	writer.WriteString(message + "\r\n")
	writer.Flush()
	log.Println("Sent:", message)

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	log.Println("Received:", line)

	if !strings.Contains(line, expected) {
		return "", err
	}

	return line, nil
}

func handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer conn.Close()

	_, err := sendAndReceive(writer, reader, "220 Hello", "EHLO")
	if err != nil {
		log.Println("Failed to read:", err)
		return
	}

	_, err = sendAndReceive(writer, reader, "250 OK", "AUTH LOGIN")
	if err != nil {
		log.Println("Failed to read:", err)
		return
	}

	_, err = sendAndReceive(writer, reader, "334 Username", usernameBase64)
	if err != nil {
		log.Println("Failed to read:", err)
		return
	}

	_, err = sendAndReceive(writer, reader, "334 Password", passwordBase64)
	if err != nil {
		log.Println("Failed to read:", err)
		return
	}

	_, err = sendAndReceive(writer, reader, "235 OK", "MAIL FROM")
	if err != nil {
		log.Println("Failed to read:", err)
		return
	}

	_, err = sendAndReceive(writer, reader, "250 OK", "RCPT TO")
	if err != nil {
		log.Println("Failed to read:", err)
		return
	}

	_, err = sendAndReceive(writer, reader, "250 OK", "DATA")
	if err != nil {
		log.Println("Failed to read:", err)
		return
	}

	writer.WriteString("354 Ready\r\n")
	writer.Flush()
	log.Println("Sent: 354 Ready")

	var mailBody strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Failed to read:", err)
			return
		}

		if strings.TrimRight(line, "\r\n") == "." {
			break
		}

		mailBody.WriteString(line)
	}
	handleMailBody(mailBody.String())

	_, err = sendAndReceive(writer, reader, "250 OK", "QUIT")
	if err != nil {
		log.Println("Failed to read:", err)
		return
	}

	writer.WriteString("221 Bye\r\n")
	writer.Flush()
	log.Println("Sent: 221 Bye")
}

func handleMailBody(mailBody string) {
	log.Println("Working on mail body data...")

	reader := strings.NewReader(mailBody)
	msg, err := mail.ReadMessage(reader)
	if err != nil {
		log.Println("Failed to parse mail body:", err)
		return
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		log.Println("Failed to parse 'Content-Type' in mail body:", err)
		return
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		log.Println("Multipart mail body detected")
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Error during multipart parsing:", err)
				return
			}

			//log.Println("Headers:", p.Header)

			slurp, err := io.ReadAll(p)
			if err != nil {
				log.Println("Error during multipart parsing (slurp):", err)
				return
			}

			disposition, params, err := mime.ParseMediaType(p.Header.Get("Content-Disposition"))
			if err != nil {
				log.Println("Failed to parse Content-Disposition:", err)
				continue
			}

			filename := params["file_name"]
			if filename != "" && disposition == "attachment" {
				log.Println("Attachment detected:", filename)
				sanitized := strings.ReplaceAll(string(slurp), "\n", "")
				decodedAttachment, err := base64.StdEncoding.DecodeString(sanitized)
				if err != nil {
					log.Println("Failed to decode attachment base64:", err)
					return
				}
				newReader := bytes.NewReader(decodedAttachment)
				imagesReadersChan <- newReader
			}
		}
	}
}
