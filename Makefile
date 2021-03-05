deploy:
	env GOOS=linux GOARCH=arm go build
	scp backend-projekt mama.sh:.
	ssh mama.sh 'mv backend-projekt tmp.mama.sh/api; echo newpdk1234 | sudo -S systemctl restart tmp.service'
