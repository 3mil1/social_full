import { BrowserRouter } from "react-router-dom";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import App from "./App";
import store from "./store/store";
// UI Material
import { ThemeProvider } from "@mui/material";
import { theme } from "./theme";
import "./index.scss";
import { AlertSnackbar } from "./components/alert";

const container = document.getElementById("container");
const root = createRoot(container!);

root.render(
  <BrowserRouter>
    <Provider store={store}>
      <ThemeProvider theme={theme}>
        <AlertSnackbar />
        <App />
      </ThemeProvider>
    </Provider>
  </BrowserRouter>
);
