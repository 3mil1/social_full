import http from "./http-common";

export default {
  async getMsgs(id, skip, limit, shouldDelete) {
    // console.log("R:", receiver);
    try {
      const msgs = await http.get(
        `/chat/?with=${id}&skip=${skip}&limit=${limit}`
      );
      console.log("Got history", msgs);
      if (shouldDelete) {
        await http.delete(`/user/notification/reply?id=${id}`);
      }
      return msgs.data;
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
      console.log(msgs);
      let m = [];
      if (msgs.data !== null) {
        msgs.data.forEach((msg) => {
          const ms = {
            content: msg.content,
            data: msg.created_at,
            from: msg.from,
            name: msg.first_name + " " + msg.last_name,
            read: msg.seen, //prob don't need this field
          };
          m.push(ms);
        });
      }
      console.log("Parsed msgs", m);
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
      console.log("New list of followers", list);
      return list.data;
    } catch (e) {
      console.error(e);
    }
  },
};
