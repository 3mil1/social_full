import userService from "../../utilities/user-service";
import { useNavigate } from "react-router-dom";
import { Button } from "@mui/material";
import { useDispatch } from "react-redux";
import LogoutIcon from "@mui/icons-material/Logout";
import "./logout.scss";
import WsApi from "../../utilities/ws";
import * as helper from "../../helpers/HelperFuncs";

export default function Logout() {
  let redirect = useNavigate();
  const dispatch = useDispatch();

  const handleLogout = async () => {
    console.log("logout fired");
    await userService.logout().then(() => {
      dispatch({ type: "LOGOUT" });
      WsApi.stop();
      localStorage.removeItem("chat_with");
      redirect("/login", { replace: true });
    });
  };
  return (
    <>
      <Button className="logout_link" onClick={handleLogout}>
        {" "}
        <LogoutIcon fontSize="large" />
      </Button>
    </>
  );
}
