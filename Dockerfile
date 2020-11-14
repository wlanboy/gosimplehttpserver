FROM busybox:1.31

COPY ./gosimplehttpserver /home/
COPY ./.env /home/
EXPOSE 7000

CMD ["/home/gosimplehttpserver"]
