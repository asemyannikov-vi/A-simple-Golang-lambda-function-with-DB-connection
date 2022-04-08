
SVC_DIR        := $(shell pwd)
GO             := $(shell which go)

LAMBDA_NAME    := main
LAMBDA_DIR     := ${SVC_DIR}/lambdas/${LAMBDA_NAME}
LAMBDAS_DOCKER := ${LAMBDA_DIR}/docker-compose.yaml
LAMBDA_NETWORK := service-network
SAM_SUBNET     := 172.21.0.0/16
SAM_GATEWAY    := 172.21.0.1

prepare-lambda-package:
	${GO} build ${LAMBDA_DIR}/${LAMBDA_NAME}.go && mv ${SVC_DIR}/${LAMBDA_NAME} ${LAMBDA_DIR}
	zip -r ${LAMBDA_DIR}/${LAMBDA_NAME}.zip -j ${LAMBDA_DIR}/${LAMBDA_NAME} ${LAMBDA_DIR}/scripts
	docker network create --subnet=${SAM_SUBNET} --gateway=${SAM_GATEWAY} ${LAMBDA_NETWORK}
	docker-compose --file ${LAMBDAS_DOCKER} up -d

lambda-network-down:
	docker-compose --file ${LAMBDAS_DOCKER} down
	docker network rm ${LAMBDA_NETWORK}
	rm ${LAMBDA_DIR}/${LAMBDA_NAME}.zip

launch-lambda-with-sam:
	sam build --template-file ${LAMBDA_DIR}/template.yaml
	sam local start-api