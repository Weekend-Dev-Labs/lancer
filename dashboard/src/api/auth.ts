import { AxiosResponse } from "axios";
import { axiosInstances } from "./axios";

interface ApiAuthPayload {
  email: string;
  password: string;
}

interface ApiAuthRes {
  token: string;
}

export const apiLogin = (data: ApiAuthPayload) => {
  return axiosInstances.post<any , AxiosResponse<ApiAuthRes>>("admin/login", data);
};

export const getAuthHeaders = () => {
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
};


