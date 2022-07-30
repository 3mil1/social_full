import axios, { AxiosResponse } from "axios";
import {
  removeAccessToken,
  removeRefreshToken,
  setAccessToken,
  setRefreshToken,
} from "./token";

const a = axios.create({
  baseURL: "http://localhost:8080/",
  headers: {
    "Content-type": "application/json",
  },
  withCredentials: true,
});

a.interceptors.request.use((config) => {
  const aConfig = config;
  const token = localStorage.getItem("accessToken");
  if (token) {
    if (aConfig.headers) {
      aConfig.headers.Authorization = `Bearer ${token}`;
    }
  }
  return aConfig;
});

export default {
  get: async (url: string): Promise<AxiosResponse<any, any>> => {
    try {
      return await a.get(url);
    } catch (e: any) {
      if (e.response) {
        if (e.response.data.status_code === 401) {
          await refresh();
          return await a.get(url);
        }
      }
      throw e;
    }
  },
  post: async (url: string, data?: any): Promise<AxiosResponse<any, any>> => {
    try {
      return await a.post(url, data);
    } catch (e: any) {
      if (e.response) {
        if (e.response.data.status_code === 401) {
          await refresh();
          return await a.post(url, data);
        }
      }
      throw e;
    }
  },
  put: async (url: string, data?: any): Promise<AxiosResponse<any, any>> => {
    try {
      return await a.put(url, data);
    } catch (e: any) {
      if (e.response) {
        if (e.response.data.status_code === 401) {
          await refresh();
          return await a.get(url);
        }
      }
      throw e;
    }
  },
  delete: async (url: string, data?: any) => {
    try {
      await a.delete(url, data);
    } catch (e) {
      console.log(e);
    }
  },
};

async function refresh() {
  const rT = localStorage.getItem("refreshToken");
  try {
    const response = await a.post("user/refresh", {
      refresh_token: rT,
    });
    setAccessToken(response.data.access_token);
    setRefreshToken(response.data.refresh_token);
    console.log("refreshed");
  } catch (e) {
    // @ts-ignore
    if (e.response.data.status_code === 401) {
      localStorage.removeItem("userInfo");
      window.location.replace("/login");
      removeAccessToken();
      removeRefreshToken();
    } else {
      throw e;
    }
  }
}
