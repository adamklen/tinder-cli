# tinder-cli
Command line interface for Tinder, written in Go

## Token & User Id
Instructions to find your Facebook token and user id:

1. Navigate your browser to https://www.facebook.com/v2.6/dialog/oauth?redirect_uri=fb464891386855067%3A%2F%2Fauthorize%2F&display=touch&state=%7B%22challenge%22%3A%22IUUkEUqIGud332lfu%252BMJhxL4Wlc%253D%22%2C%220_auth_logger_id%22%3A%2230F06532-A1B9-4B10-BB28-B29956C71AB1%22%2C%22com.facebook.sdk_client_state%22%3Atrue%2C%223_method%22%3A%22sfvc_auth%22%7D&scope=user_birthday%2Cuser_photos%2Cuser_education_history%2Cemail%2Cuser_relationship_details%2Cuser_friends%2Cuser_work_history%2Cuser_likes&response_type=token%2Csigned_request&default_audience=friends&return_scopes=true&auth_type=rerequest&client_id=464891386855067&ret=login&sdk=ios&logger_id=30F06532-A1B9-4B10-BB28-B29956C71AB1&ext=1470840777&hash=AeZqkIcf-NEW6vBd"
2. Use your browser's dev tools to look at network requests. Find the request for the url https://www.facebook.com/v2.6/dialog/oauth/confirm?dpr=1
3. From the response, grab the access_token field. This is your Facebook token.
4. Navigate to the url https://graph.facebook.com/me?access_token={YOUR_ACCESS_TOKEN} and extract "id" from the json response. This is your Facebook user id.
5. Put these fields into main.go, do a make, and run ./tinder-cli

