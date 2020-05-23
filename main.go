package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-ldap/ldap"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

var (
	ldapBindUsername = os.Getenv("LDAP_BIND_USERNAME")
	ldapBindPassword = os.Getenv("LDAP_BIND_PASSWORD")
	ldapDC           = os.Getenv("LDAP_DC")
	crtFileName      = os.Getenv("CRT_FILENAME")
	keyFileName      = os.Getenv("KEY_FILENAME")
	ldapServer       = os.Getenv("LDAP_SERVER")
	ldapPort         = os.Getenv("LDAP_PORT")
	radiusSecret     = os.Getenv("RADIUS_SECRET")
	debug            = os.Getenv("DEBUG")
	ldapConn         *ldap.Conn
	tlsConf          *tls.Config
	cert             tls.Certificate
)

/**
 * Reads the crt and key files
 * Connects ldapConn to G Suite's LDAP server
 * Any error in here is fatal.
 */
func startLDAP() {
	var err error
	//Get the cert files
	crt, err := ioutil.ReadFile(crtFileName)
	if err != nil {
		log.Fatal("Error reading crt: " + err.Error())
	}
	key, err := ioutil.ReadFile(keyFileName)
	if err != nil {
		log.Fatal("Error reading key: " + err.Error())
	}
	//Combine the cert files
	cert, err = tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatal("Error creating x509 key pair: " + err.Error())
	}
	//Create the tls.config
	tlsConf = &tls.Config{ServerName: ldapServer, Certificates: []tls.Certificate{cert}}
	//Connect to G Suite
	ldapConn, err = ldap.DialTLS("tcp", ldapServer+":"+ldapPort, tlsConf)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error dialing: %v", err))
	}

	//Check to see if we should run in debug mode
	if debug == "true" {
		ldapConn.Debug.Enable(true)
	}

	//We actually don't want this closed.
	//defer ldapConn.Close()

	//Use our bound
	err = ldapConn.Bind(ldapBindUsername, ldapBindPassword)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error binding: %v", err))
	}

}

/**
 * Radius passes the username and password to this func.
 * returns an error.
 * If error is nil, then username and password are correct.
 * If error is not nil, then there was a problem.
 * This actually connects to G Suite LDAP as the user to test.
 * Meaning a whole new connection that immediately gets closed.
 */
func ldapTryUsernamePassword(username string, password string) error {
	searchRequest := ldap.NewSearchRequest(
		ldapDC,
		// The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=person)(uid=%s))", username), // The filter to apply
		[]string{"dn", "cn"}, // A list attributes to retrieve
		nil,
	)

	sr, err := ldapConn.Search(searchRequest)
	if err != nil {
		log.Println("Error in ldapRtUsernamePassword while ldapConn.Search: " + err.Error())
		e := errors.New("Error in ldapRtUsernamePassword while ldapConn.Search: " + err.Error())
		return e
	}
	sr.PrettyPrint(1)
	if len(sr.Entries) != 1 {
		err = errors.New("User does not exist or too many entries returned")
		return err
	}
	userdn := sr.Entries[0].DN
	l, err := ldap.DialTLS("tcp", ldapServer+":"+ldapPort, tlsConf)
	if err != nil {
		return err
	}
	if debug == "true" {
		l.Debug.Enable(true)
	}

	defer l.Close()
	err = l.Bind(userdn, password)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Calls startLDAP to connect to the server
 * Starts the radius server.
 */
func main() {
	logFile, err := os.OpenFile("gsuiteldapradius.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	startLDAP()

	handler := func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		password := rfc2865.UserPassword_GetString(r.Packet)

		var code radius.Code
		if ldapTryUsernamePassword(username, password) == nil {
			code = radius.CodeAccessAccept
		} else {
			code = radius.CodeAccessReject
		}
		log.Printf("Writing %v to %v", code, r.RemoteAddr)
		w.Write(r.Response(code))
	}

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(radiusSecret)),
	}

	log.Printf("Starting server on :1812")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error starting Radius Server: " + err.Error())
	}
}
