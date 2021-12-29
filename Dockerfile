FROM alpine:edge
COPY dungeonQuest-linux-amd64 /usr/bin/dungeonQuest
RUN addgroup -g 1000 -S dungeonQuest && \
    adduser -h /home/dungeonQuest -g 'dungeonQuest,,,,' -s /bin/sh -S -D -u 1000 dungeonQuest
WORKDIR /home/dungeonQuest/
COPY maps /home/dungeonQuest/maps
CMD dungeonQuest -client /home/dungeonQuest/BrowserQuest -config /home/dungeonQuest/config.json