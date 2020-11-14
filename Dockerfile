FROM busybox:1.31

ADD gosimplehttpserver /

RUN mkdir -p /usr/local/share/busybox && echo "/bin/busybox sh" > /usr/local/share/busybox/sh && chmod +x /usr/local/share/busybox/sh
RUN addgroup -S kanban && adduser -S kanban -G kanban -s /usr/local/share/busybox/sh

USER kanban

EXPOSE 7000

CMD ["/gosimplehttpserver"]
