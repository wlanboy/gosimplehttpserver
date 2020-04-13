FROM busybox:1.31

ARG BIN_FILE
ADD ${BIN_FILE} /home/gosimplehttpserver

EXPOSE 7000

ENTRYPOINT ["/home/gosimplehttpserver"]
