S3_BUCKET = rcj-sam-packages
STACK_NAME = srdbot

build:
	@mkdir -p dist
	@for handler in `ls handlers`; do \
		GOOS=linux go build -o dist/$$handler handlers/$$handler/*.go ;\
	done

test:
	@for handler in `ls handlers`; do \
		BASE=`pwd`; \
		cd handlers/$$handler; \
		go test; \
		cd $$BASE; \
	done

clean:
	@rm -rf dist

package:
	@aws cloudformation package --template-file template.yaml --s3-bucket $(S3_BUCKET) --output-template-file package.yaml

deploy: package
	@aws cloudformation deploy --template-file package.yaml --capabilities CAPABILITY_IAM --stack-name $(STACK_NAME)
	@#aws cloudformation describe-stacks --stack-name $(STACK_NAME) | jq '.Stacks[0].Outputs[] | select(.OutputKey=="Endpoint").OutputValue' 

destroy:
	@aws cloudformation delete-stack --stack-name $(STACK_NAME)