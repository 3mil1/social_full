export function setAccessToken(token: string) {
  localStorage.setItem("accessToken", token);
}

export function setRefreshToken(token: string) {
  localStorage.setItem("refreshToken", token);
}

export function removeAccessToken() {
  localStorage.removeItem("accessToken");
}

export function removeRefreshToken() {
  localStorage.removeItem("refreshToken");
}


