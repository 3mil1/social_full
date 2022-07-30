import { fireEvent, render, screen, waitFor } from "@testing-library/react";
// import Login from "./login";
// import { MemoryRouter, Router } from "react-router-dom";
// import userService from "../../utilities/user-service";
// import { createMemoryHistory } from "history";

// //checks whether login component is on the page
// it("should render Login component on the screen", () => {
//   render(
//     <MemoryRouter>
//       <Login />
//     </MemoryRouter>
//   );
//   expect(screen.getByTestId("login-form")).toBeInTheDocument();
// });

// const testEmail = "test@mail.com";
// const testPwd = "0000";

// it("sends request with correct data", async () => {
//   //history to check redirection after submit
//   const history = createMemoryHistory();
//   render(
//     <Router navigator={history} location={"/"}>
//       <Login />
//     </Router>
//   );
//   //getting fields from login-for
//   const email = screen.getByTestId("email-input");
//   const pwd = screen.getByTestId("pwd-input");
//   const btn = screen.getByTestId("submit-btn");
//   //insert values
//   fireEvent.change(email, { target: { value: testEmail } });
//   fireEvent.change(pwd, { target: { value: testPwd } });

//   //mock login method from userService
//   const spy = jest.spyOn(userService, "login").mockResolvedValue();

//   fireEvent.click(btn);

//   await waitFor(() => expect(spy).toBeCalled());
//   //check login arguments passed to fn are correct
//   expect(spy.mock.calls[0][0]).toBe(testEmail);
//   expect(spy.mock.calls[0][1]).toBe(testPwd);
//   //check redirection
//   await waitFor(() => {
//     expect(history.location.pathname).toBe("/profile");
//   });
// });

