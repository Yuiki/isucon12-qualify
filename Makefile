.PHONY: create-stack renew-pem

create-stack:
	aws cloudformation create-stack --template-body file://cloudformation.yaml --parameters ParameterKey=ImageId,ParameterValue=ami-<AMI_ID> --capabilities CAPABILITY_IAM --stack-name isucon12-qualify

renew-pem:
	bash -c "openssl x509 -in <(openssl req -subj '/CN=*.t.isucon.dev' -nodes -newkey rsa:2048 -keyout nginx/tls/key.pem) -req -signkey ./nginx/tls/key.pem -sha256 -days 3650 -out nginx/tls/fullchain.pem -extfile <(echo -e "basicConstraints=critical,CA:true,pathlen:0\nsubjectAltName=DNS.1:*.t.isucon.dev")"
