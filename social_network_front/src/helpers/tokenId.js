const getTokenId = () => {
  let t = localStorage.getItem("accessToken");

  if (!t) return "";
  let tokenEncode = localStorage.getItem("accessToken").split(".")[1];
  let token = tokenEncode.replace("-", "+").replace("_", "/");
  return JSON.parse(window.atob(token)).user_id;
};

export default getTokenId;
