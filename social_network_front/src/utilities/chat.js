import http from "./http-common";
import { parseDate } from "../helpers/parseDate";

export default {
  async getMsgs(id, skip, limit, shouldDelete) {
    try {
      const msgs = await http.get(
        `/chat/?with=${id}&skip=${skip}&limit=${limit}`
      );
      if (shouldDelete) {
        await http.delete(`/user/notification/reply?id=${id}`);
      }
      let messages = [];
      msgs?.data.forEach((m) => {
        const msg = {
          content: m.content,
          data: m.created_at,
          from: m.from,
        };
        messages.push(msg);
      });
      return messages;
    } catch (err) {
      console.error(err);
      throw err;
    }
  },

  async getGroupMsgs(id, skip, limit, shouldDelete) {
    try {
      const msgs = await http.get(
        `/group/chat?groupId=${id}&skip=${skip}&limit=${limit}`
      );
      let m = [];
      if (msgs.data !== null) {
        msgs.data.forEach((msg) => {
          const ms = {
            content: msg.content,
            data: parseDate(msg.created_at, false),
            from: msg.from,
            name: msg.first_name + " " + msg.last_name,
          };
          m.push(ms);
        });
      }
      if (shouldDelete) {
        await http.delete(`/user/notification/reply?id=${id}`);
      }
      return m;
    } catch (err) {
      console.error(err);
    }
  },

  async getUserList() {
    try {
      const list = await http.get("/follower/chat");
      return list.data;
    } catch (e) {
      console.error(e);
    }
  },
};
