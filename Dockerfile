FROM busybox:1.31

ARG BIN_FILE
ADD ${BIN_FILE} /

EXPOSE 7000

ENTRYPOINT ["/gosimplehttpserver"]
