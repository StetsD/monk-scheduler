FROM mustbestd/monk-docker-ci-precompiled:v0.1

ARG UID=1000
ARG GID=1000
ARG COMMAND=start

ENV USER_NAME=monk \
    APP_NAME=monkapp \
    W_DIR=/go/src \
    GID=${GID} \
    UID=${UID} \
    SRC_DIR=./gitbub.com/stetsd/monk-app \
    COMMAND=${COMMAND} \
    GOOS=linux \
    GOARCH=amd64

RUN groupadd --gid $GID $USER_NAME && \
    useradd -u $UID --gid $USER_NAME --shell /bin/bash --create-home $USER_NAME

WORKDIR $W_DIR

COPY . $SRC_DIR

RUN cd $SRC_DIR; make build && \
    mv $APP_NAME /usr/local/bin

USER $USER_NAME

CMD ["make", "setup"]
CMD ["make", "migrate"]

ENTRYPOINT monkapp $COMMAND