PROGNAME=SshEngine
OUTPUT=${PROGNAME}_${CI_COMMIT_TAG}_${GOOS}_${GOARCH}

cross:
	go build -o ${OUTPUT}
	sha256sum ${OUTPUT} >${OUTPUT}.sha256
