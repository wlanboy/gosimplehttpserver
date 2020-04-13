FROM busybox:1.31

COPY ./gosimplehttpserver /home/
EXPOSE 8000

CMD ["/home/gosimplehttpserver"]
