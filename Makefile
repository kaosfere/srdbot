S3_BUCKET = rcj-sam-packages
STACK_NAME = srdbot

build:
	for file in `ls handlers/*.go`; do \
		GOOS=linux go build -o `echo $$file | sed -e 's/.go$$//;'` $$file ;\
	done

package:
	aws cloudformation package --template-file template.yaml --s3-bucket $(S3_BUCKET) --output-template-file package.yaml

deploy:
	@aws cloudformation deploy --template-file package.yaml --capabilities CAPABILITY_IAM --stack-name $(STACK_NAME)
	@aws cloudformation describe-stacks --stack-name $(STACK_NAME) | jq '.Stacks[0].Outputs[] | select(.OutputKey=="Endpoint").OutputValue' 

destroy:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)