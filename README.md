# go-searchme
Go API that utilizes Groupme APIs to parse recent group messages by keyword

# env vars
API needs a the following env variables to be set
* `SINCE_ID` - can be obtained from GroupMe account
* `GROUP_ID` - id of the group to query messages from
* `TOKEN` - API token from GroupMe
* `REDIS_PORT` - redis port
* `REDIS_PASS` - redis secret pass
* `REDIS_HOST` - redis host ip

# frontend assets
* cd into `frontend/` dir and run `npm run build` to compile assets and create dist folder for frontend changes to take affect