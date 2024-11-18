# # stage 1 - build stage
FROM node:18-alpine as build

# # set working directory
# WORKDIR /app

# # install dependencies
# COPY web/package*.json ./
# RUN npm install

# # copy all source code to build it
# COPY web ./

# # build the react app
# RUN npm run build

# # stage 2 - serve stage
# FROM nginx:alpine

# # copy build from the previous stage
# COPY --from=build /app/build /usr/share/nginx/html

# # expose port
# EXPOSE 80

# # start nginx
# CMD ["nginx", "-g", "daemon off;"]