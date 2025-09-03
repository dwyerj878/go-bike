FROM golang:1.24.5-bookworm

ARG USERNAME=dev
ARG USER_UID=1000
ARG USER_GID=$USER_UID

ENV LANG=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8
ENV LANGUAGE=en_US.UTF-8


RUN apt-get update \
    && apt-get install -y locales sudo curl ca-certificates git make openssh-client procps less unzip zip gnupg2 telnet iputils-ping git uuid awscli \
    && sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen \
    && locale-gen 

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m --shell /bin/bash $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME


COPY --chown=${USERNAME}:${USER_GID} .bash_profile /home/dev/.bash_profile
RUN  . /usr/share/bash-completion/completions/git

USER $USERNAME
WORKDIR /home/dev/

