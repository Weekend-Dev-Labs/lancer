import axios from "axios";

export const axiosInstances = axios.create({
  baseURL: "/api",
});
