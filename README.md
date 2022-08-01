# Social Network

Authors:<br>
 - Emil Varnomasing (3mil)<br>
 - Valeria Kharchenko (ValeriaKharchenko)<br>
 - Silver Luhtoja (SilverL)<br>
 - Anna Lazarenkova (anna_lazarenkova)

How to run: 
1. in terminal type: bash  start.sh
2. wait til docker has finished installing and setting up
3. open broswer on "http://localhost:3000/:

### Objectives
This project is a Facebook-like social network that will contain the following features:
* Followers
* Profile
* Posts
* Groups
* Notification
* Chats
- - - -
### Authentication
In order for the users to use the social network they will have to make an account. 
When the user logins, he/she should stay logged in until he/she chooses a logout option that should be available at all times. 
- - - -

### Followers
When navigating the social network the user should be able to follow and unfollow other users. 
In order to follow someone the user first needs to send a request to the user he/she wants to follow. Then the other user should be able to accept the request or decline it. If the second user has a public profile this step is skipped and the sending of the request is ignored.
- - - -

### Profile
There are two types of profiles: a public profile and a private profile. 
The user should be able to turn its profile public or private.
- - - -

### Posts
The user must be able to specify the privacy of the post(can include an image or GIF):
* public (all users in the social network will be able to see the post)
* private (only followers of the creator of the post will be able to see the post)
* almost private (only the followers chosen by the creator of the post will be able to see it)
- - - -

### Groups & Events
A user is able to create a group.<br>
Only the creator of the group would be allowed to accept or refuse the "joining" request.<br>
Posts and comments of the group will only be displayed to members of the group.<br>
A user belonging to the group can also create an event, making it available for the other group users.<br>
After creating the event every user of that group can choose one of the options (going/ not going/ interested) for the event.
- - - -

### Chat
The user is able to send private messages/emojis to users that he/she is following.
Groups have a common chat room, so if a user is a member of the group he/she is able to send and receive messages to this group chat.
Notification about any new message will be shown.
- - - -

### Notifications
A user will be notified if he/she:

* has a private profile and some other user sends him/her a following request (mandatory)
* receives a group invitation, so he can refuse or accept the request (mandatory)
* is the creator of a group and another user requests to join the group, so he can refuse or accept the request (mandatory)
* is member of a group and an event is created (mandatory)
* new comment to user post will be added (optional)
* group access will be opened for user (optional)
