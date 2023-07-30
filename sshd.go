package main

import (
	"errors"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/ajpikul-com/wsssh/wsconn"
)

func ReadAuthorizedKeys() {
	// Right now it's doing this everytime, which is not great TODO
	authorizedKeysBytes, err := os.ReadFile(clients)
	if err != nil {
		panic("Failed to load auth keys file " + err.Error())
	}
	globalConfig.authorizedKeysMap := map[string]string{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, comment, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			panic(err.Error())
		}
		globalConfig.authorizedKeysMap[string(pubKey.Marshal())] = comment
		authorizedKeysBytes = rest
	}

	// TODO we should update globa state with possible clients
	globalClientList.UpdateClientList(possibleClients) // This probably needs to be processed, not sure what we accept
}

func GetServer(wsconn *wsconn.WSConn, clients string, privateKey string) (*ssh.ServerConn, <-chan ssh.NewChannel, <-chan *ssh.Request, error) {

	ReadAuthorizedKeys() // This isn't good to have here!

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			defaultLogger.Debug("Running Public Key Callback")
			comment, ok := globalConfig.authorizedKeysMap[string(pubKey.Marshal())]
			if ok {
				return &ssh.Permissions{
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
						"comment":   comment,
					},
				}, nil
			}
			return nil, errors.New("No access for " + c.User())
		},
	}
	privateBytes, err := os.ReadFile(privateKey)
	if err != nil {
		panic("Problem loading private key file")
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Couldn't parse private key")
	}

	config.AddHostKey(private)

	return ssh.NewServerConn(wsconn, config)
}
