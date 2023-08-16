package template

var BuildTemplate = `rm -rf ./app
rm -rf %s.tar.gz
mkdir app
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./%s ./main.go
chmod +x ./%s
cp %s ./app/
cp -R ./conf ./app/
cp -R ./docs ./app/
cp run.sh ./app/
cp pm2.yml ./app/
tar -zcvf %s.tar.gz ./app
`
