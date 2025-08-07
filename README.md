# shared-browser-ide

This is an excuse to try out websockets and go to share text across several browsers in real time.
You can create sessions or join an existing one and share code across browsers.

This repo contains code for both the react frontend (built using vite) and the go backend. The FE is served through the go backend as static files.

## How to run

To get the app rolling locally simply run `docker-compose up --build` and both the react app and go server will build and start running on port 8080 of your machine.
