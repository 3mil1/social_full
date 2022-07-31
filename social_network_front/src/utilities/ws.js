import {
  addNotification,
  updateNotifications,
} from "../store/notificationSlice";
import { addMsg } from "../store/chatSlice";
import * as helper from "../helpers/HelperFuncs";

let ws;
const cleanUp = () => {
  ws?.removeEventListener("close", ws?.onclose);
  ws?.removeEventListener("message", ws?.onmessage);
  ws?.removeEventListener("open", ws?.onopen());
};

export default {
  start(id, dispatcher) {
    cleanUp();
    ws?.close();
    ws = new WebSocket("ws://localhost:8080/ws/");

    ws.onopen = () => {
      let jsonData = {};
      jsonData["action"] = "connect";
      jsonData["user"] = id;
      ws.send(JSON.stringify(jsonData));
    };

    ws.onmessage = (msg) => {
      const msgJSON = JSON.parse(msg.data);
      let notificationList = [];
      let receiver = localStorage.getItem("chat_with");
      let sender = helper.getTokenId();
      let location = window.location.href.includes("/chat");

      if (Array.isArray(msgJSON)) {
        msgJSON.forEach((m) => {
          switch (m.action_type) {
            case "private message":
              if (m.data.from === receiver || m.data.from === sender) {
                dispatcher(addMsg(m.data));
              }
              break;
            case "group message":
              const newMsg = {
                content: m.data.content,
                data: m.data.created_at,
                from: m.data.from,
                name: m.data.first_name + " " + m.data.last_name,
                group_id: m.data.group_id,
              };
              if (`${m.data.group_id}` === receiver || m.data.from === sender) {
                dispatcher(addMsg(newMsg));
              }
              break;
            case "new message in group chat":
              if (!location || `${m.data.group_id}` !== receiver) {
                dispatcher(addNotification(m.data.group_id));
              }
              break;
            case "new private message":
              if (!location || m.data.actor_id !== receiver) {
                dispatcher(addNotification(m.data.actor_id));
              }
              break;
            default:
              notificationList.push(m);
          }
        });
        dispatcher(updateNotifications(notificationList));
      }
    };
  },

  stop() {
    let jsonData = {};
    jsonData["action"] = "left";

    ws.send(JSON.stringify(jsonData));
    cleanUp();
    ws?.close();
  },

  sendChatMessage(message) {
    ws.send(message);
  },
};
