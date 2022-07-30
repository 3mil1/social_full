import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useSelector } from "react-redux";
import { RootState } from "../store/store";
import Navbar from "../components/Navbar";

export const Private = () => {
  const auth = useSelector((state: RootState) => state.profile.auth);
  const location = useLocation();
  if (!auth) {
    return <Navigate to={"/login"} state={{ from: location }} />;
  }
  return (
    <>
      <div className="left">
        <Navbar />
      </div>
      <div className="right">
        <Outlet />
      </div>
    </>
  );
};

export const Public = () => {
  const auth = useSelector((state: RootState) => state.profile.auth);
  const location = useLocation();
  if (auth) {
    let currentPath = location.pathname;
    if (currentPath === "/login") {
      return (
        <Navigate to={"/homepage"} state={{ from: location }} replace={true} />
      );
    }
    return <Navigate to={currentPath} state={{ from: location }} />;
  }
  return (
    <main className="screen-center">
      <Outlet />
    </main>
  );
};
