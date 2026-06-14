FROM alpine

COPY validating-kontroller /validating-kontroller

ENTRYPOINT [ "./validating-kontroller" ]