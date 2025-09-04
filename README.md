### コマンド（Powershell）
$env:GOOS = "linux"  
go env GOOS CGO_ENABLED GOARCH  
set GOOS=linux  
set GOARCH=amd64  
set CGO_ENABLED=0  
go build -tags "lambda.norpc,dev" -o bootstrap .  
build-lambda-zip -o myFunction.zip bootstrap  