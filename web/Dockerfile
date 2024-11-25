# stage 1 - build stage
FROM node:18-alpine as build

# stage 2 - serve stage
FROM nginx:alpine

# set working directory
WORKDIR /app

# install dependencies
COPY package.json ./
COPY package-lock.json ./
RUN npm install

# copy all source code to build it
COPY . ./

# build the react app
RUN npm run build

# copy custom nginx configuration
# COPY nginx.conf /etc/nginx/nginx.conf

# copy build from the previous stage
COPY --from=build /app/dist /usr/share/nginx/html

# expose port
EXPOSE 80

# start nginx
CMD ["nginx", "-g", "daemon off;"]
