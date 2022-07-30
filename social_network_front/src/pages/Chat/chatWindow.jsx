import React, { useEffect, useRef, useState } from "react";
import "./chat.scss";
import chatService from "../../utilities/chat";
import { setAlert } from "../../store/alertSlice";
import {
  addToBegining,
  loadMsgs,
  setFollowerList,
} from "../../store/chatSlice";
import GroupService from "../../utilities/group_service";

// import Picker from "emoji-picker-react";
import InputEmoji from "react-input-emoji";

//mui material
import { Button, Divider, Grid, ListItem } from "@mui/material";
import Paper from "@mui/material/Paper";
import List from "@mui/material/List";
import ListItemText from "@mui/material/ListItemText";
import FaceIcon from "@mui/icons-material/Face";
import GroupsIcon from "@mui/icons-material/Groups";
import ChatBubbleOutlineIcon from "@mui/icons-material/ChatBubbleOutline";
import { useDispatch, useSelector } from "react-redux";
import * as helper from "../../helpers/HelperFuncs";
import WsApi from "../../utilities/ws";
import { removeNotification } from "../../store/notificationSlice";

export const Chat = () => {
  // let followerList = useSelector((state) => state.followers.followers);
  // const [followerList, setFollowerList] = useState({});
  const followerList = useSelector((state) => state.chat.followers);
  let notifications = useSelector((state) => state.notifications.messages);
  // console.log("Notifications from useSelector", notifications);
  let dispatch = useDispatch();
  const [receiver, setReceiver] = useState({ id: "", type: "" });
  // console.log("Receiver", receiver);

  const msgs = useSelector((state) => state.chat.msgHistory);
  let sender = helper.getTokenId();

  //load or update group store
  const group_service = GroupService();
  useEffect(() => {
    group_service.getCreatedGroups();
    group_service.getJoinedGroups();
  }, []);

  const getFollowers = async () => {
    await chatService.getUserList().then((followers) => {
      dispatch(setFollowerList(followers));
    });
  };

  useEffect(() => {
    getFollowers();
  }, []);

  //create list of groups
  const createdGroups = useSelector((state) => state.groups.createdGroups);
  const joinedGroups = useSelector((state) => state.groups.joinedGroups);
  let groups = createdGroups.concat(joinedGroups);

  // console.log("Groups:", groups);
  // console.log("Followers:", followerList);

  //lead to one obj type
  let members = [];
  followerList?.forEach((f) => {
    let chatMember = {
      name: f.first_name + " " + f.last_name,
      id: f.user_id,
      type: "person",
    };
    members.push(chatMember);
  });

  groups?.forEach((g) => {
    let group = {
      name: g.title,
      id: `${g.id}`,
      type: "group",
    };
    members.push(group);
  });
  // console.log("Members", members);

  // 👇️ scroll to bottom every time messages change
  const bottomRef = useRef(null);
  const topRef = useRef(null);
  const [lastMsg, setLastMsg] = useState("");

  useEffect(() => {
    if (msgs && lastMsg !== msgs[msgs.length - 1]) {
      setLastMsg(msgs[msgs.length - 1]);
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [msgs]);

  //load more msgs
  const [hasMore, setHasMore] = useState(true);

  function loadMore() {
    setTimeout(() => {
      loadHistory(msgs.length).then(() => topRef.current?.scrollIntoView());
    }, 2000);
  }

  //emoji
  const [chosenEmoji, setChosenEmoji] = useState(null);
  const onEmojiClick = (event, emojiObject) => {
    setChosenEmoji(emojiObject);
  };

  //load msg history
  const loadHistory = async (s) => {
    let msgHistory = [];
    let shouldDelete = notifications.includes(receiver.id);
    if (receiver.type === "person") {
      try {
        msgHistory = await chatService.getMsgs(
          receiver.id,
          s,
          10,
          shouldDelete
        );
      } catch (e) {
        console.log(e.message);
        const errorState = {
          text: "Can't load messages",
          severity: "warning",
        };
        dispatch(setAlert(errorState));
      }
    } else if (receiver.type === "group") {
      // console.log("receiver", receiver);
      try {
        msgHistory = await chatService.getGroupMsgs(
          receiver.id,
          s,
          10,
          shouldDelete
        );
        console.log("group messages", msgHistory);
      } catch (e) {
        console.error(e.message);
        const errorState = {
          text: "Can't load messages",
          severity: "warning",
        };
        dispatch(setAlert(errorState));
      }
    }
    console.log("response", msgHistory);
    setHasMore(msgHistory !== null && msgHistory.length === 10);
    if (s === 0) {
      dispatch(loadMsgs(msgHistory));
    } else {
      dispatch(addToBegining(msgHistory));
    }
  };
  useEffect(() => {
    if (receiver.id !== "") {
      loadHistory(0);
      dispatch(removeNotification(receiver.id));
    }
  }, [receiver]);

  //send chat message
  const [text, setText] = useState("");
  const sendMsg = (text) => {
    // console.log(text);
    if (text.trim().length > 0) {
      let jsonData = {};
      if (receiver.type === "person") {
        jsonData["action"] = "message";
      } else {
        jsonData["action"] = "group_chat";
      }
      jsonData["user"] = sender;
      jsonData["message_to"] = receiver.id;
      jsonData["message_content"] = text;
      // console.log("JSON DATA", JSON.stringify(jsonData));
      WsApi.sendChatMessage(JSON.stringify(jsonData));
      setText("");
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  };

  //alert if no chat chosen
  const fireAlert = () => {
    const errorState = {
      text: "Select chat to send message",
      severity: "warning",
    };
    dispatch(setAlert(errorState));
  };

  // let msgHeight = document.getElementsByClassName("messageArea");
  // console.log("height: ", msgHeight[0].scrollHeight);

  return (
    <div className={"fullWidth"}>
      <Grid container component={Paper} className="chatSection">
        <Grid item xs={3} className={"borderRight500"}>
          <List>
            {members.length === 0 ? (
              <ListItem key={"warningMsg"}>
                <ListItemText>
                  Follow somebody or enter a group to start chat
                </ListItemText>
              </ListItem>
            ) : (
              members.map((member, i) => {
                return (
                  <ListItem key={i.toString()}>
                    {/*<FaceIcon />*/}
                    {member.type === "person" ? <FaceIcon /> : <GroupsIcon />}
                    <ListItemText>
                      <Button
                        className={member.id === receiver.id ? "active" : ""}
                        onClick={() => {
                          localStorage.setItem("chat_with", member.id);
                          setReceiver({ id: member.id, type: member.type });
                        }}
                        fullWidth
                      >
                        {member.name}{" "}
                        {notifications.includes(member.id) && (
                          <ChatBubbleOutlineIcon
                            sx={{
                              color: "#D7B271",
                              fontSize: 20,
                              marginLeft: 1,
                            }}
                          />
                        )}
                      </Button>
                    </ListItemText>
                  </ListItem>
                );
              })
            )}
          </List>
        </Grid>
        <Grid item xs={9}>
          <List className={"messageArea"}>
            {receiver.id !== "" && hasMore && (
              <ListItem key="loadMoreBtn">
                <Button
                  sx={{ marginLeft: 25 }}
                  variant="text"
                  onClick={loadMore}
                  className={"load-more-btn"}
                >
                  {" "}
                  Load more...{" "}
                </Button>
              </ListItem>
            )}
            {receiver.id === "" ? (
              <ListItem key={"selectChat"}>
                <ListItemText>Select chat</ListItemText>
              </ListItem>
            ) : msgs.length === 0 ? (
              <ListItem key={"noMsg"}>
                <ListItemText>No messages yet</ListItemText>
              </ListItem>
            ) : (
              msgs.map((m, i, l) => {
                return (
                  <ListItem key={i.toString()}>
                    {i === 9 && <div ref={topRef} />}
                    <Grid container>
                      <Grid item xs={12}>
                        <ListItemText
                          className={m.from === sender ? "right" : "left"}
                        >
                          <div className={"additional-info user-name"}>
                            {m.name}
                          </div>
                          <div>{m.content}</div>
                          <div className={"additional-info"}>{m.data}</div>
                        </ListItemText>
                      </Grid>
                    </Grid>
                  </ListItem>
                );
              })
            )}
            <div ref={bottomRef} />
          </List>
          <Divider />
          {/*<form ref={form} onSubmit={handleSubmit(sendMsg)}>*/}
          <Grid container style={{ padding: "20px" }}>
            <Grid item xs={10}>
              <InputEmoji
                id="msg-input"
                label="Type a message"
                // multiline
                maxRows={2}
                height={600}
                // fullWidth
                // cleanOnEnter
                maxLength={400}
                value={text}
                onChange={setText}
                onEnter={receiver.id !== "" ? sendMsg : fireAlert}
              />
            </Grid>
            {/*<Grid item xs={1} marginLeft={1} align="right">*/}
            {/*  /!*<Picker onEmojiClick={onEmojiClick} />*!/*/}
            {/*  <Fab color="primary" aria-label="add" type={"submit"}>*/}
            {/*    <SendIcon />*/}
            {/*  </Fab>*/}
            {/*</Grid>*/}
          </Grid>
          {/*</form>*/}
        </Grid>
      </Grid>
    </div>
  );
};

export default Chat;
