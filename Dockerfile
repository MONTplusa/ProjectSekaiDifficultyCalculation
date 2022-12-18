FROM golang:latest

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go mod init github.com/MONTplusa/ProjectSekaiDifficultyCalculation
RUN go mod tidy
RUN go build

CMD ["./ProjectSekaiDifficultyCalculation"]
