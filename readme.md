# GLDAPRadius

Radius Server for Google's G Suite LDAP Directory  
It's very simple and should have a pretty small footprint. It is tested on Ubuntu 20.04 on Raspberry Pi 4 and Ubuntu 22.04 amd64. The intent of this app is to be used on Google Cloud's Compute Engine or Cloud Run. 

Because Google's LDAP doesn't expose the password hash, you MUST use PEAP as the authentcation mechanism. CHAP and MSCHAPv2 do no send the password. They assume the server knows the password. This is a huge annoyance from Google. They should either expose the password hash or create their own Radius-as-a-service. Contact your rep and let them know that!

There is a Dockerfile here for Cloud Run. However I have only tested on Compute Engine.

## How To Use

To use this radius server, you first need to have Google [Workspace Enterprise](https://support.google.com/a/answer/7284269?hl=en) or [Cloud Identity](https://cloud.google.com/identity). You must also follow [Googles steps](https://support.google.com/a/topic/9048334?hl=en&ref_topic=7556782) setup the LDAP directory.

1. In linux, install git, golang and ufw.

2. Get the dependencies:
go mod tidy

3. Copy the example files to their correct 
cp build.example.sh build.sh

4. Copy the .crt and .key file you got from Google while adding the LDAP client.

5. Make build.sh executable
chmod +x build.sh

6. Edit build.sh and app.yaml
build.sh is used on linux to set environmental variables, compile the app and run the app.
app.yaml is used by app engine to do the same. Set environmental variables and run the app.

LDAP_BIND_USERNAME: Username from G Suite LDAP client creation  
LDAP_BIND_PASSWORD: Password from G Suite LDAP client creation  
LDAP_DC: Your domain. example: "dc=foo,dc=com" for foo.com  
CRT_FILENAME: crt file from G Suite LDAP client creation  
KEY_FILENAME: key file from G Suite LDAP client creation  
LDAP_SERVER: "ldap.google.com"  
LDAP_PORT: "636"  
RADIUS_SECRET: Long random string that you create. I use 30 or more characters.  
DEBUG: "false" is default. set to "true" if you want verbose logging. You don't want verbose logging  

7. Security, Firewall, Etc.

If you are using a linux server, use a firewall.

sudo ufw default deny incoming  
sudo ufw default allow outgoing  
sudo ufw allow ssh  
sudo ufw allow 1812  
sudo ufw enable  

The above commands closes all incoming but opens incoming ssh and radius ports to the whole internet.
If you want to use ufw to only allow acces from certain IPs or subnets, [look here](https://www.digitalocean.com/community/tutorials/ufw-essentials-common-firewall-rules-and-commands).

With Google Cloud Firewall or any firewall, you can restrict access to specific subnets.

Also probably a good idea to have DOS protection.

## Future Features
1. LDAP group restriction
2. IP and Subnet whitelists & blacklists
3. Accounting

## Built With

* [Basic LDAP v3 functionality for the GO programming language.](github.com/go-ldap/ldap)
* [Google Workspace](https://gsuite.google.com/)
* [layeh.com/radius ](https://github.com/layeh/radius)


## License

I really didn't build anything. I just put a few pieces together. Whatever license the go libraries use, just abide by those.
