This reverse proxy helps with local testing of frontend and backend hosting both from same host, and avoid CORS and cookie security problems.

Backend URL path starts with `/api` and forwarded to port `:8080` without `/api`.

Other requests considered frontend, forwarded to port `:3000` as is.
