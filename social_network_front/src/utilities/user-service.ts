import http from "./http-common";
import { RegisterForm } from "../pages/Register/register";
import {
  removeAccessToken,
  removeRefreshToken,
  setAccessToken,
  setRefreshToken,
} from "./token";



export default {
  async login(email: string, pwd: string){
    try {
      const response = await http.post("user/signin", {
        email: email,
        password: pwd,
      });

      setAccessToken(response.data.access_token);
      setRefreshToken(response.data.refresh_token);
      
    } catch (err) {
      throw err;
    }
  },

  async register(user: RegisterForm) {
    try {
      console.log(
        "%c Sending user registration data to server",
        "color:green",
        user
      );
      await http.post("user/signup", {
        email: user.email,
        password: user.password,
        nickname: user.nickname,
        first_name: user.first_name,
        last_name: user.last_name,
        birth_day: user.dob,
        about_me: user.desc,
        user_img: user.image_path,
      });
    } catch (err) {
      console.log("Error caught");
      throw err;
    }
  },

  async logout() {
    http
      .delete("user/signout")
      .then(() => {
        removeAccessToken();
        removeRefreshToken();
        console.log("User logged out");
      })
      .catch((err) => {
        console.log(err);
      });
  },
};


/// SAVED //////////////////////////////////////////////////

// import http from "./http-common";
// import { RegisterForm } from "../pages/Register/register";

// import {
//   removeAccessToken,
//   removeRefreshToken,
//   setAccessToken,
//   setRefreshToken,
// } from "./token";
// import { useState } from "react";

// export interface UserInfo {
//   firstName: string;
//   lastName: string;
//   id: string;
//   auth: boolean;
// }

// export default {
//   async login(email: string, pwd: string): Promise<UserInfo> {
//     try {
//       const response = await http.post("user/signin", {
//         email: email,
//         password: pwd,
//       });
//       console.log(response.data);

//       //get tokens and save to local storage
//       setAccessToken(response.data.access_token);
//       setRefreshToken(response.data.refresh_token);
//       // console.log(response);
//       return {
//         firstName: "",
//         lastName: "",
//         id: "",
//         auth: true,
//       };
//     } catch (err) {
//       throw err;
//     }
//   },

//   async register(user: RegisterForm) {
//     try {
//       console.log(
//         "%c Sending user registration data to server",
//         "color:green",
//         user
//       );
//       await http.post("user/signup", {
//         email: user.email,
//         password: user.password,
//         nickname: user.nickname,
//         first_name: user.first_name,
//         last_name: user.last_name,
//         birth_day: user.dob, //now it has type Date
//         about_me: user.desc,
//         user_img: user.image_path,
//       });
//     } catch (err) {
//       console.log("Error caught");
//       throw err;
//     }
//   },

//   // async auth(): Promise<UserInfo> {
//   //   try {
//   //     const user = await http.get("user/");
//   //     console.log("User auth", user);
//   //     if (user.status === 200) {
//   //       return {
//   //         firstName: user.data.firstName,
//   //         lastName: user.data.lastName,
//   //         id: user.data.ID,
//   //       };
//   //     }
//   //     return {
//   //       firstName: "",
//   //       lastName: "",
//   //       id: "",
//   //     };
//   //   } catch (err) {
//   //     console.log("auth wasn't completed", err);
//   //     throw err;
//   //   }
//   // },

//   async logout() {
//     http
//       .delete("user/signout")
//       .then(() => {
//         removeAccessToken();
//         removeRefreshToken();
//         console.log("User logged out");
//       })
//       .catch((err) => {
//         console.log(err);
//       });
//   },
//   // async followRequest() {
//   //   const resp = await http
//   //     .post("follower/user/", {
//   //       source_id: "c0df434a-3ea6-4796-818e-a3b7b1a6ec97",
//   //       target_id: "0e3e82bc-1808-456c-b37c-b6eefd88d60a",
//   //     })
//   //     .then((resp) => {
//   //       if (resp.status === 200) {
//   //         console.log("Followed successfully");
//   //       }
//   //     })
//   //     .catch((err) => {
//   //       console.log(err);
//   //     });
//   // },

//   // async profile() {
//   //   try {
//   //     const response = await http.get("user/me");
//   //     console.log("Response from profile", response.data);
//   //   } catch (e) {
//   //     console.log("Couldn't fetch", e);
//   //   }
//   // },
// };
