package ftpserver

import (
	"bytes"
	"io"
	"log"
	"suntek2telegram/pkg/config"

	"goftp.io/server/v2"
)

type (
	MyAuth   struct{}
	MyPerm   struct{}
	MyDriver struct{}
)

var (
	ftpConf           *config.FTP
	imagesReadersChan chan io.Reader
)

func Start(fc *config.FTP, imgReadersChan chan io.Reader) {
	ftpConf = fc
	imagesReadersChan = imgReadersChan

	myDriver := &MyDriver{}
	myPerm := &MyPerm{}
	myAuth := &MyAuth{}
	opts := &server.Options{
		Driver: myDriver,
		Perm:   myPerm,
		Auth:   myAuth,

		Port:     fc.BindPort,
		Hostname: fc.BindHost,

		PassivePorts: fc.PassivePorts,
		PublicIP:     fc.PublicIP,
	}

	server, err := server.NewServer(opts)
	if err != nil {
		log.Fatalln(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Failed to start FTP server:", err)
	}
}

func (auth *MyAuth) CheckPasswd(ctx *server.Context, s1, s2 string) (bool, error) {
	result := s1 == ftpConf.Username && s2 == ftpConf.Password
	return result, nil
}

func (driver *MyDriver) PutFile(ctx *server.Context, destPath string, data io.Reader, appendData int64) (int64, error) {
	var buf bytes.Buffer
	copied, err := io.Copy(&buf, data)
	if err != nil {
		return copied, err
	}

	newReader := bytes.NewReader(buf.Bytes())
	imagesReadersChan <- newReader

	// Workaround - because camera do not have keep-alive mechanism (and NAT is ALWAYS used with
	// mobile operators), so let's just close connection after each image, so it re-establishes it
	// for each image
	ctx.Sess.Close()

	return copied, nil
}
