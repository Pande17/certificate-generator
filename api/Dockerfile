# Stage 1 - wkhtmltopdf dependencies
FROM surnet/alpine-wkhtmltopdf:3.20.0-0.12.6-full AS wkhtmltopdf

# Stage 2 - Go with wkhtmltopdf and air autoreload
FROM golang:1.22-alpine3.20

# Add wkhtmltopdf required package
RUN apk add --no-cache \
    libstdc++ libx11 libxrender \
    libxext fontconfig freetype \
    ttf-droid ttf-freefont ttf-liberation \
    bash ca-certificates wget

# Copy the wkhtmltopdf binary from "wkhtmltopdf" reference image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf    /usr/local/bin/wkhtmltopdf
COPY --from=wkhtmltopdf /bin/wkhtmltoimage  /usr/local/bin/wkhtmltoimage
COPY --from=wkhtmltopdf /bin/libwkhtmltox*  /usr/local/bin

# Ensure binaries are executable
RUN chmod +x /usr/local/bin/wkhtmltopdf && \
    chmod +x /usr/local/bin/wkhtmltoimage

# # perintah untuk mengecek keberadaan file wkhtmltopdf
# # RUN ls -l /usr/local/bin/wkhtmltopdf

# # Add /usr/local/bin to PATH explicitly
ENV PATH="/usr/local/bin:${PATH}"

# # Verify that wkhtmltopdf is correctly installed (optional)
# RUN wkhtmltopdf --version

# Install Go Air autoreload package
RUN go install github.com/air-verse/air@v1.52.3
# RUN go install github.com/cosmtrek/air@v1.52.3

# Set working directory to /app
WORKDIR /app

# # Copy Go module files and download dependencies
# COPY /go.mod /go.sum ./
# RUN go mod download

# # Copy all files from api to /app in container
COPY . .

# # Build the Go application
# RUN go build -o certifgen-backend

# Run autoreload or direct binary depending on environment
CMD ["air", "-c", ".air.toml"]

# Expose port 3000
EXPOSE 3000
