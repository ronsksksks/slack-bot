build-linux: 
	GOOS=linux GOARCH=amd64 go build main.go

deploy:
	mv main mura-bot
	scp -i ~/aws.pem mura-bot ec2-user@13.231.109.9: