FROM golang:1.10-alpine

ARG subscriptionId
ARG clientId
ARG clientSecret
ARG tenantId
ARG metricAggregation
ARG resourceQueryTagName
ARG environment
ARG release
ARG elkHost
ARG elkIndex
ARG activeLogSegregation

ENV subscriptionId $subscriptionId
ENV clientId $clientId
ENV clientSecret $clientSecret
ENV tenantId $tenantId
ENV metricAggregation $metricAggregation
ENV resourceQueryTagName $resourceQueryTagName
ENV environment $environment
ENV release $release
ENV elkHost $elkHost
ENV elkIndex $elkIndex
ENV activeLogSegregation $activeLogSegregation

WORKDIR /go/src/app
COPY . .

RUN apk add util-linux && apk add --update git

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 9276
CMD ["app"]