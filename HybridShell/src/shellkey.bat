"C:\Program Files\OpenSSL-Win64\bin\openssl.exe" genrsa -out pvt.key 4096
"C:\Program Files\OpenSSL-Win64\bin\openssl.exe" rsa -in pvt.key -pubout > pub.pem