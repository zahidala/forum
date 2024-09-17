FROM golang:latest
LABEL maintainer="Zahid Allaulddin (zbutt), Mohammed Ali (moali), malamri (Manal Alamri), and obhari (Omar Bahri)"
LABEL version="1.0"
LABEL description="Dockerized version of Forum"
LABEL source="https://learn.reboot01.com/git/zbutt/forum"
WORKDIR /home
ARG GITEA_URL=learn.reboot01.com
ARG PAT=90cf6c95c9bbe6e5c24849d6df06422f996a3d03
RUN git clone https://$PAT@${GITEA_URL}/git/zbutt/forum.git /home/forum
WORKDIR /home/forum
# create env file
RUN echo "FREEIMAGEHOST_API_KEY=6d207e02198a847aa98d0a2a901485a5" > .env
RUN go build -o /forum
EXPOSE 8080
CMD [ "/forum" ]