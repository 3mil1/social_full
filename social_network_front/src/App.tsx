import { Route, Routes } from "react-router-dom";
import { Private, Public } from "./hoc/routeWrappers";
import Pages from "./pages/pages";
import { useEffect } from "react";
import ProfileService from "./utilities/profile_service";
import { useDispatch } from "react-redux";
import WsApi from "./utilities/ws";
import * as helper from "./helpers/HelperFuncs";
import "./index.scss";

function App() {
  const profile_service = ProfileService();
  const dispatch = useDispatch();
  useEffect(() => {
    let id = helper.getTokenId();
    profile_service.checkAuth();
    if (id) {
      WsApi.start(id, dispatch);
    }
  }, []);

  return (
    <>
      <Routes>
        <Route element={<Private />}>
          <Route path={"/homepage"} element={<Pages.Homepage />} />
          <Route path={"/profile/:id"} element={<Pages.Profile />} />
          <Route path={"/group/:id"} element={<Pages.Group />} />
          <Route path={"/group/:groupId/post/:postId"} element={<Pages.GroupPost />} />
          <Route path={"post/:id"} element={<Pages.OnePost />} />
          <Route path={"/notifications"} element={<Pages.Notification />} />
          <Route path={"/chat"} element={<Pages.Chat />} />
          <Route path="/*" element={<Pages.OnePost />} />
        </Route>
        <Route element={<Public />}>
          <Route path={"register"} element={<Pages.Register />} />
          <Route path={"login"} element={<Pages.Login />} />
        </Route>
      </Routes>
    </>
  );
}

export default App;
