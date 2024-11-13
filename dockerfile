# add wkhtmltopdf image source based on alpine 3.20 with "wkhtmltopdf" as reference name
FROM surnet/alpine-wkhtmltopdf:3.20.0-0.12.6-full as wkhtmltopdf

# use golang 1.22 on alpine 3.20 as image base
FROM golang:1.22-alpine3.20

# add wkhtmltopdf required package
RUN apk add --no-cache \
    libstdc++ libx11 libxrender \
    libxext fontconfig freetype \
    ttf-droid ttf-freefont ttf-liberation

# copy the wkhtmltopdf binary from "wkhtmltopdf" reference image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf    /usr/local/bin/wkhtmltopdf
COPY --from=wkhtmltopdf /bin/wkhtmltoimage  /usr/local/bin/wkhtmltoimage
COPY --from=wkhtmltopdf /bin/libwkhtmltox*  /usr/local/bin

#install golang air autoreload package
RUN go install github.com/air-verse/air@v1.52.3

# mount project to this directory in container
WORKDIR /app

# run autoreload
CMD ["air", "-c", ".air.toml"]

EXPOSE 3000