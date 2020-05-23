#!/bin/bash
#rename this file to build.sh
#Change all the environmental vars below to what G Suite provides
#and what you want the radius secret to be.
#Change debug to true if you want to see verbose logs.
#Put the crt and key file from G Suite in this folder.
export LDAP_BIND_USERNAME="UsernameFromGSuiteLdap"
export LDAP_BIND_PASSWORD="PasswordFromGSuiteLdap"
export LDAP_DC="dc=foo,dc=com"
export CRT_FILENAME="From_GSuite_LDAP.crt"
export KEY_FILENAME="From_GSuite_LDAP.key"
export LDAP_SERVER="ldap.google.com"
export LDAP_PORT="636"
export RADIUS_SECRET="Long-Random-String-Probably-32-Characters"
export DEBUG="false"
rm gldapradius
go build -o gldapradius main.go
./gldapradius